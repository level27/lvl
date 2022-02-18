package cmd

import (
	"bitbucket.org/level27/lvl/types"
	"github.com/spf13/cobra"
)

// MAIN COMMAND
var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Commands for managing systems",
}

func init() {
	//adding main command to root
	RootCmd.AddCommand(systemCmd)

	//Toplevel subcommands (get/post)
	systemCmd.AddCommand(systemGetCmd)
}

var systemGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get a list of all curent systems",
	Run: func(cmd *cobra.Command, args []string) {
		outputFormatTable(getSystems(args), []string{"ID", "NAME", "STATUS"}, []string{"Id", "Name", "Status"})

	},
}

func getSystems(amount []string) []types.SystemGet {

	if len(amount) == 0 {
		return Level27Client.SystemGetList()
	}
	return Level27Client.SystemGetList()

}

// // func getDomains(ids []string) []types.Domain {
// // 	c := Level27Client
// // 	if len(ids) == 0 {
// // 		return c.Domains(optGetParameters)
// // 	} else {
// // 		domains := make([]types.Domain, len(ids))
// // 		for idx, id := range ids {
// // 			domains[idx] = c.Domain("GET", id, nil)
// // 		}
// // 		return domains
// // 	}
// // }
