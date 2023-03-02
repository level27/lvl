//go:build windows

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

// Windows-specific code for self-updating lvl.
//
// Updating the exe is annoying, Windows won't let us just delete the file while it's running.
// We have to:
// 1. rename the old executable while it's running (this is allowed)
// 2. rename the new executable into place (no problem)
// 3. spawn a new process from the new exe to delete the old exe
// 4. exit the old process
// 5. now the new process can delete the old one (lock has to be released)

func init() {
	RootCmd.AddCommand(winFinishUpdateCmd)
}

var winFinishUpdateCmd = &cobra.Command{
	Use:    "__winfinishupdate",
	Hidden: true,

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]

		var err error
		// Try a few times in case the old process hasn't exited yet.
		for i := 0; i < 5; i++ {
			time.Sleep(1 * time.Second)

			err = os.Remove(file)
		}

		return err
	},
}

func updateSwapFile(new string, old string) error {
	execSwapPath := old + ".old"
	err := os.Rename(old, execSwapPath)
	if err != nil {
		return err
	}

	err = os.Rename(new, old)
	if err != nil {
		return err
	}

	command := exec.Command(old, "__winfinishupdate", execSwapPath)
	err = command.Start()
	if err != nil {
		return err
	}

	return nil
}

func getAssetFileName() string {
	return fmt.Sprintf("lvl-windows-%s.exe", runtime.GOARCH)
}
