package cmd

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

func init() {
	systemNetworkCmd.AddCommand(systemNetworkIpCmd)

	systemNetworkIpCmd.AddCommand(systemNetworkIpGetCmd)

	systemNetworkIpCmd.AddCommand(systemNetworkIpAddCmd)
	systemNetworkIpAddCmd.Flags().StringVar(&systemNetworkIpAddHostname, "hostname", "", "Hostname for the IP address. If not specified the system hostname is used.")

	systemNetworkIpCmd.AddCommand(systemNetworkIpRemoveCmd)

	systemNetworkIpCmd.AddCommand(systemNetworkIpUpdateCmd)
	settingsFileFlag(systemNetworkIpUpdateCmd)
	settingString(systemNetworkIpUpdateCmd, updateSettings, "hostname", "New hostname for this IP")
}

func resolveSystemHasNetworkIP(systemID l27.IntID, hasNetworkID l27.IntID, arg string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.LookupSystemHasNetworkIp(systemID, hasNetworkID, arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"system network IP",
		func(app l27.SystemHasNetworkIp) string { return fmt.Sprintf("%d", app.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, err
}

var systemNetworkIpCmd = &cobra.Command{
	Use:   "ip",
	Short: "Manage IP addresses on network connections",
}

var systemNetworkIpGetCmd = &cobra.Command{
	Use:   "get [system] [network]",
	Short: "Get all IP addresses for a system network",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		networkID, err := resolveSystemHasNetwork(systemID, args[1])
		if err != nil {
			return err
		}

		ips, err := Level27Client.SystemGetHasNetworkIps(systemID, networkID)
		if err != nil {
			return err
		}

		outputFormatTableFuncs(ips, []string{"ID", "Public IP", "IP", "Hostname", "Status"}, []interface{}{"ID", func(i l27.SystemHasNetworkIp) string {
			if i.PublicIpv4 != "" {
				i, _ := strconv.Atoi(i.PublicIpv4)
				if i == 0 {
					return ""
				} else {
					return ipv4IntToString(i)
				}
			} else if i.PublicIpv6 != "" {
				ip := net.ParseIP(i.PublicIpv6)
				return fmt.Sprint(ip)
			} else {
				return ""
			}
		},
			func(i l27.SystemHasNetworkIp) string {
				if i.Ipv4 != "" {
					i, _ := strconv.Atoi(i.Ipv4)
					if i == 0 {
						return ""
					} else {
						return ipv4IntToString(i)
					}
				} else if i.Ipv6 != "" {
					ip := net.ParseIP(i.Ipv6)
					return fmt.Sprint(ip)
				} else {
					return ""
				}
			}, "Hostname", "Status"})

		return nil
	},
}

var systemNetworkIpAddHostname string

var systemNetworkIpAddCmd = &cobra.Command{
	Use:   "add [system] [network] [address]",
	Short: "Add IP address to a system network",
	Long:  "Adds an IP address to a system network. Address can be either IPv4 or IPv6. The special values 'auto' and 'auto-v6' automatically fetch an unused address to use.",

	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		system, err := Level27Client.SystemGetSingle(systemID)
		if err != nil {
			return err
		}

		hasNetworkID, err := resolveSystemHasNetwork(systemID, args[1])
		if err != nil {
			return err
		}

		network, err := Level27Client.GetSystemHasNetwork(systemID, hasNetworkID)
		if err != nil {
			return err
		}

		networkID := network.Network.ID
		address := args[2]

		if address == "auto" || address == "auto-v6" {
			located, err := Level27Client.NetworkLocate(networkID)
			if err != nil {
				return err
			}

			var choices []string
			if address == "auto" {
				choices = located.Ipv4
			} else {
				choices = located.Ipv6
			}

			if len(choices) == 0 {
				return errors.New("unable to find a free IP address")
			}

			address = choices[0]
		}

		var data l27.SystemHasNetworkIpAdd
		public := network.Network.Public

		if strings.Contains(address, ":") {
			// IPv6
			if public {
				data.PublicIpv6 = address
			} else {
				data.Ipv6 = address
			}
		} else {
			// IPv4
			if public {
				data.PublicIpv4 = address
			} else {
				data.Ipv4 = address
			}
		}

		data.Hostname = system.Hostname
		if systemNetworkIpAddHostname != "" {
			data.Hostname = systemNetworkIpAddHostname
		}

		ip, err := Level27Client.SystemAddHasNetworkIps(systemID, hasNetworkID, data)
		if err != nil {
			return err
		}

		outputFormatTemplate(ip, "templates/entities/systemNetworkIp/add.tmpl")
		return nil
	},
}

var systemNetworkIpRemoveCmd = &cobra.Command{
	Use:   "remove [system] [network] [address | id]",
	Short: "Remove IP address from a system network",

	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		hasNetworkID, err := resolveSystemHasNetwork(systemID, args[1])
		if err != nil {
			return err
		}

		ipID, err := resolveSystemHasNetworkIP(systemID, hasNetworkID, args[2])
		if err != nil {
			return err
		}

		err = Level27Client.SystemRemoveHasNetworkIps(systemID, hasNetworkID, ipID)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/systemNetworkIp/remove.tmpl")
		return nil
	},
}

var systemNetworkIpUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update settings on a system network IP",

	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		settings, err := loadMergeSettings(updateSettingsFile, updateSettings)
		if err != nil {
			return err
		}

		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		hasNetworkID, err := resolveSystemHasNetwork(systemID, args[1])
		if err != nil {
			return err
		}

		ipID, err := resolveSystemHasNetworkIP(systemID, hasNetworkID, args[2])
		if err != nil {
			return err
		}

		ip, err := Level27Client.SystemGetHasNetworkIp(systemID, hasNetworkID, ipID)
		if err != nil {
			return err
		}

		ipPut := l27.SystemHasNetworkIpPut{
			Hostname: ip.Hostname,
		}

		data := mergeSettingsWithEntity(ipPut, settings)

		err = Level27Client.SystemHasNetworkIpUpdate(systemID, hasNetworkID, ipID, data)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/systemNetworkIp/update.tmpl")
		return err
	},
}
