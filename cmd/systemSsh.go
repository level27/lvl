package cmd

import (
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/kevinburke/ssh_config"
	"github.com/level27/l27-go"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

// Contains commands like lvl system ssh.

func init() {
	// SYSTEM SSH
	systemCmd.AddCommand(systemSshCmd)

	// SYSTEM SCP
	systemCmd.AddCommand(systemScpCommand)

	// SYSTEM SSHCONFIG
	systemCmd.AddCommand(systemSshConfigCmd)
}

// SYSTEM SSH
var systemSshCmd = &cobra.Command{
	Use:   "ssh <system> [flags] [--] [ssh args]",
	Short: "Connect to a system via SSH, automatically adding SSH keys to the system if necessary",
	Long: `Connect to a system via SSH, automatically adding SSH keys to the system if necessary.
The command will automatically add your favorite SSH key to the system if necessary. You'll need to use 'lvl sshkey favorite' to configure it the first time.
The command figures out a valid IP address to connect to and passes it through to the ssh command.
Arguments passed after the system ID/name are passed to ssh literally. Note that for any flags starting with "-", you'll want to put "--" before them so lvl does not try to interpret them as flags itself.`,
	Example: `lvl system ssh my-awesome-server
lvl system ssh my-awesome-server ls "~"
lvl system ssh my-awesome-server -- ls -l "~"`,

	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		favoriteKeyID := viper.GetInt32("ssh_favoritekey")
		if favoriteKeyID == 0 {
			return fmt.Errorf("no favorite SSH key configured. Use 'lvl sshkey favorite' to configure one")
		}

		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		// We need to do two things:
		// 1. Make sure we have an SSH key on the system.
		// 2. Fetch the host to pass in the ssh command.
		// We send these as concurrent tasks to reduce latency on the command.

		taskSshKey := taskRunVoid(func() error {
			return waitEnsureSshKey(systemID, favoriteKeyID)
		})

		taskSshHost := taskRun(func() (string, error) {
			return sshResolveHost(systemID, false)
		})

		sshHost := <-taskSshHost
		if sshHost.Error != nil {
			return sshHost.Error
		}

		err = <-taskSshKey
		if err != nil {
			return err
		}

		sshArgs := []string{fmt.Sprintf("root@%s", sshHost.Result)}
		sshArgs = append(sshArgs, args[1:]...)

		return tailExecProcess("ssh", sshArgs)
	},
}

// Ensure the given SSH key is available and 'ok' on a system.
func waitEnsureSshKey(systemID l27.IntID, sshKeyID l27.IntID) error {
	_, err := Level27Client.SystemSshKeysGetSingle(systemID, sshKeyID)
	if err == nil {
		// No error, so key exists.
		return nil
	}

	// Error, might indicate SSH key doesn't exist yet.
	_, ok := err.(l27.ErrorResponse)
	if !ok {
		// Not an API error, could be network failure or something instead, abort.
		return err
	}

	// TODO: check error code above, isn't currently correct thanks to PL-7611
	// For now we assume it's just a 404, so try to add the SSH key.

	err = waitAddSshKey(systemID, sshKeyID)
	return err
}

// Add an SSH key to a system, waiting for the status to change to 'ok'.
func waitAddSshKey(systemID l27.IntID, sshKeyID l27.IntID) error {
	fmt.Fprint(os.Stderr, "Adding SSH key to system")

	key, err := Level27Client.SystemAddSshKey(systemID, sshKeyID)
	if err != nil {
		return err
	}

	waitIndicator(func() {
		_, err = waitForStatus(
			func() (l27.SystemSshkey, error) { return Level27Client.SystemSshKeysGetSingle(systemID, key.ID) },
			func(ss l27.SystemSshkey) string { return ss.ShsStatus },
			"ok",
			[]string{"updating"},
		)
	})

	if err != nil {
		return fmt.Errorf("waiting for SSH key to change status failed: %s", err.Error())
	}

	return nil
}

// Resolve the hostname to SSH into a system.
// If the FQDN properly resolves, we use that.
// Otherwise we try the IP addresses in the system's networks.
func sshResolveHost(systemID l27.IntID, preferIP bool) (string, error) {
	system, err := Level27Client.SystemGetSingle(systemID)
	if err != nil {
		return "", err
	}

	ips, err := net.LookupIP(system.Fqdn)
	if err == nil && len(ips) > 0 {
		// FQDN resolves, pass it to the ssh command.
		if preferIP {
			return ips[0].String(), nil
		}

		return system.Fqdn, nil
	}

	sort.Slice(system.Networks, func(i int, j int) bool {
		netA := system.Networks[i]
		netB := system.Networks[j]

		return netA.NetPublic && netB.NetInternal
	})

	for _, net := range system.Networks {
		for _, ip := range net.Ips {
			if ip.PublicIpv4 != "" {
				return ip.PublicIpv4, nil
			}

			if ip.Ipv4 != "" {
				return ip.Ipv4, nil
			}
		}
	}

	// Couldn't find anything.
	return "", fmt.Errorf("unable to find a suitable address to connect to on system '%s' (%d)", system.Name, system.ID)
}

// SYSTEM SCP
var systemScpCommand = &cobra.Command{
	Use:     "scp [system1:]file1 ... [system2:]file2",
	Short:   "Copy files to/from the system using scp",
	Long:    "Uses the same syntax as regular scp. Arguments are passed through, but host names (before the :) are interpreted as system names/IDs and resolved. To pass flags through to scp, put them after a --",
	Example: "lvl system scp foo.txt mySystem:~/foo.txt",

	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		favoriteKeyID := viper.GetInt32("ssh_favoritekey")
		if favoriteKeyID == 0 {
			return fmt.Errorf("no favorite SSH key configured. Use 'lvl sshkey favorite' to configure one")
		}

		// This code is quite complex to be able to be as optimally concurrent (and fast) as possible.
		// Basically how it works:
		// For each argument, we need to resolve the system if it's a system:file argument
		// These resolves are done concurrently. They send back the new value of the arg when they're done.
		// We also need to make sure SSH keys are added, this goes via another set of channels to also be concurrent.

		// Goroutine to asynchronously add SSH keys to systems while we go through resolving systems down below.
		// We need this to avoid trying to add an SSH key to the same system twice, causing race conditions.
		keyAddChannel := make(chan l27.IntID)
		keyDone := taskRunVoid(func() error {
			var group errgroup.Group
			// Map of systems we're already handling SSH keys on, to avoid running them twice.
			systemsEnsured := map[l27.IntID]bool{}
			for systemID := range keyAddChannel {
				sysID := systemID
				if _, ok := systemsEnsured[systemID]; ok {
					continue
				}

				systemsEnsured[systemID] = true
				group.Go(func() error { return waitEnsureSshKey(sysID, favoriteKeyID) })
			}

			return group.Wait()
		})

		var systemArgsGroup errgroup.Group
		systemArgsChannel := make(chan tuple2[int, string])

		// Copy input arguments to pass them through to scp.
		// We will modify the ones that are remote files to replace system hosts with the real IP/domain
		scpArgs := append([]string{}, args...)

		argTask := taskRunVoid(func() error {
			// Go over remote arguments, and resolve the system.
			for i, arg := range args {
				split := strings.SplitN(arg, ":", 2)
				if len(split) == 1 {
					// No host specified, so local file or flag or something.
					// TODO: this means of parsing mostly works, but it means that any flag parameters with a colon in them
					// will be interpreted as a remote file.
					// It might be a good idea to manually pass-through flag args for well-known flags.
					continue
				}

				// Do this all in parallel with goroutines to avoid chaining latency, nice and fast.
				ii := i
				systemArgsGroup.Go(func() error {
					system := split[0]
					file := split[1]

					systemID, err := resolveSystem(system)
					if err != nil {
						return err
					}

					// Send ID to SSH key channel so the SSH key gets added.
					keyAddChannel <- systemID

					host, err := sshResolveHost(systemID, false)
					if err != nil {
						return err
					}

					// Send arg index and new value so the value gets replaced.
					systemArgsChannel <- makeTuple2(ii, fmt.Sprintf("root@%s:%s", host, file))
					return nil
				})
			}

			// Wait for all systems to finish resolving.
			err := systemArgsGroup.Wait()
			close(systemArgsChannel)
			return err
		})

		// Update all the args that need updating from the above loop.
		for tuple := range systemArgsChannel {
			scpArgs[tuple.Item1] = tuple.Item2
		}

		// Handle errors from argument processing.
		if err := <-argTask; err != nil {
			return err
		}

		// All args have been processed, so also all SSH keys have been dispatched at least.
		close(keyAddChannel)

		// Wait for all SSH keys to be available on systems.
		err := <-keyDone
		if err != nil {
			return err
		}

		return tailExecProcess("scp", scpArgs)
	},
}

var systemSshConfigCmd = &cobra.Command{
	Use:   "sshconfig",
	Short: "Add system's name to your user SSH config for easy access",
	Long: `Add system's name to your user SSH config for easy access
This will add a Host entry to your SSH config, so afterwards you can use commands outside lvl to access the system by name.
For example: rsync foo.txt my-awesome-system:~/

The new host names are written into a separate ~/.ssh/lvl config file, which gets added to your ~/.ssh/config via an Include directive.`,

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		address, err := sshResolveHost(systemID, true)
		if err != nil {
			return err
		}

		system, err := Level27Client.SystemGetSingle(systemID)
		if err != nil {
			return err
		}

		f, configPath, err := ensureSshConfig()
		if err != nil {
			return err
		}

		var cfg *ssh_config.Config
		(func() {
			defer f.Close()

			cfg, err = ssh_config.Decode(f)
		})()

		if err != nil {
			return fmt.Errorf("unable to parse SSH config: %s", err.Error())
		}

		host := findSshConfigHost(cfg, system.Name)
		if host == nil {
			pattern, err := ssh_config.NewPattern(system.Name)
			if err != nil {
				return fmt.Errorf("invalid SSH config host value: '%s'", system.Name)
			}

			// Host not in config file yet, add a new one!
			host = &ssh_config.Host{
				Patterns: []*ssh_config.Pattern{pattern},
			}

			cfg.Hosts = append(cfg.Hosts, host)
		}

		setSshConfigHostNode(host, "HostName", address)
		setSshConfigHostNode(host, "User", "root")

		// Write new config back to file.
		err = os.WriteFile(configPath, []byte(cfg.String()), 0o644)
		if err != nil {
			return fmt.Errorf("failed to write %s: %s", configPath, err.Error())
		}

		err = ensureSshIncludeWritten()
		if err != nil {
			return err
		}

		outputFormatTemplate(system, "templates/entities/system/sshConfigConfirm.tmpl")

		return nil
	},
}

func findSshConfigHost(cfg *ssh_config.Config, name string) *ssh_config.Host {
	for _, host := range cfg.Hosts {
		for _, pattern := range host.Patterns {
			if pattern.String() == name {
				return host
			}
		}
	}

	return nil
}

func setSshConfigHostNode(host *ssh_config.Host, key string, value string) {
	for _, node := range host.Nodes {
		kv, ok := node.(*ssh_config.KV)
		if !ok {
			continue
		}

		if kv.Key == key {
			kv.Value = value
			return
		}
	}

	// Didn't find the key in the config yet, add a new one.
	host.Nodes = append(host.Nodes, &ssh_config.KV{
		Key:   key,
		Value: value,
	})
}

func ensureSshConfig() (*os.File, string, error) {
	return ensureSshDirFileExists(getSshConfigFileName())
}

// Ensure there is an include directive to lvl's SSH config file in the user's ~/.ssh/config.
func ensureSshIncludeWritten() error {
	f, fullPath, err := ensureSshDirFileExists("config")
	if err != nil {
		return fmt.Errorf("error opening SSH config file: %s", err.Error())
	}

	var configData []byte
	(func() {
		defer f.Close()

		configData, err = io.ReadAll(f)
	})()
	if err != nil {
		return fmt.Errorf("error reading SSH config file: %s", err.Error())
	}

	configName := getSshConfigFileName()
	regex := fmt.Sprintf(`Include\s+%s\b`, regexp.QuoteMeta(configName))
	matched, err := regexp.Match(regex, configData)
	if err != nil {
		return fmt.Errorf(
			"error checking for existing include directive in SSH config file: %s",
			err.Error())
	}

	if matched {
		// Already has the include directive, don't need to do anything!
		return nil
	}

	fmt.Fprintf(os.Stderr, "Adding Include for lvl config file to SSH config\n")

	if len(configData) != 0 {
		// Copy SSH config file as a backup in case we screw up.
		err = (func() error {
			backupPath := fullPath + ".bak"
			fmt.Fprintf(os.Stderr, "Copying %s to %s as backup", fullPath, backupPath)

			backup, err := os.Create(backupPath)
			if err != nil {
				return fmt.Errorf("failed to open %s: %s", backupPath, err.Error())
			}

			backup.Write(configData)
			return nil
		})()

		if err != nil {
			return err
		}
	}

	f, err = os.OpenFile(fullPath, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %s", fullPath, err.Error())
	}

	defer f.Close()

	includeBytes := []byte(fmt.Sprintf("Include %s\n\n", configName))
	_, err = f.Write(includeBytes)

	if err != nil {
		return fmt.Errorf("error while writing to %s: %s", fullPath, err.Error())
	}

	return nil
}

// Ensure ~/.ssh/<name> exists, or create it (and the parent directory) if necessary.
func ensureSshDirFileExists(name string) (*os.File, string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, "", fmt.Errorf("unable to determine home directory: %s", err.Error())
	}

	path_ssh := path.Join(home, ".ssh")
	ssh_config := path.Join(path_ssh, name)
	f, err := os.Open(ssh_config)
	if os.IsNotExist(err) {
		// Automatically make an empty ~/.ssh/<config file>.
		_, err = os.Stat(path_ssh)
		if os.IsNotExist(err) {
			// Have to make ~/.ssh/ too...
			err = os.Mkdir(path_ssh, 0o700)
			if err != nil {
				return nil, "", fmt.Errorf("failed to create %s: %s", path_ssh, err.Error())
			}
		}

		// Make ~/.ssh/config.
		f, err = os.Create(ssh_config)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create %s: %s", ssh_config, err.Error())
		}
	} else if err != nil {
		return nil, "", fmt.Errorf("failed to open %s: %s", ssh_config, err.Error())
	}

	return f, ssh_config, nil
}

// Get the file name to write system SSH config entries to.
// This is not ~/.ssh/config itself, but we write an include directive into ~/.ssh/config for it.
func getSshConfigFileName() string {
	name := viper.GetString("ssh_config_name")
	if name != "" {
		return name
	}

	// No value given, return default.
	return "lvl"
}
