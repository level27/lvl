package cmd

import (
	"github.com/spf13/cobra"
)

// MAIN COMMAND
var systemgroupCmd = &cobra.Command{
	Use:   "systemgroups",
	Short: "Commands for managing systemgroups",
}

func init() {
	//adding main command to root
	RootCmd.AddCommand(systemgroupCmd)

	// --- GET (LIST)
	systemgroupCmd.AddCommand(systemgroupsGetCmd)
	// add optional get parameters (filters)
	addCommonGetFlags(systemgroupsGetCmd)

}

//------------------------------------------------- SYSTEMSGROUPS (GET / ADD  / UPDATE / DELETE)-------------------------------------------------
// ---------------- MAIN COMMAND (groups)
var systemgroupsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Show list of all available systemgroups.",
	Run: func(cmd *cobra.Command, args []string) {
		
		systemgroups := Level27Client.SystemgroupsGet(optGetParameters)
		
		outputFormatTable(systemgroups, []string{"ID", "NAME", "ORGANISATION"}, []string{"ID", "Name", "Organisation.Name"})
		
	},
}
