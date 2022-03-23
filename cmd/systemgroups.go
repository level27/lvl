package cmd

import (
	"github.com/spf13/cobra"
)

// MAIN COMMAND
var systemgroupCmd = &cobra.Command{
	Use:   "systemgroup",
	Short: "Commands for managing systemgroups",
}

func init() {
	//adding main command to root
	RootCmd.AddCommand(systemgroupCmd)

	// --- GET (LIST)
	systemgroupCmd.AddCommand(systemgroupsGetCmd)

}

//------------------------------------------------- SYSTEMSGROUPS (GET / ADD  / UPDATE / DELETE)-------------------------------------------------
// ---------------- MAIN COMMAND (groups)
var systemgroupsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Show list of all available systemgroups.",
	Run: func(cmd *cobra.Command, args []string) {
		loginCmd.Print("hi")
	},
}
