package cmd

import (
	"bitbucket.org/level27/lvl/types"
	"github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Commands to manage apps",
}

var appGetCmd = &cobra.Command{
	Use:  "get",
	Args: cobra.ArbitraryArgs,
	Run: func(acmd *cobra.Command, args []string) {
		outputFormatTable(
			getApps(args),
			[]string{"ID", "NAME", "STATUS"},
			[]string{"ID", "Name", "Status"})
	},
}

func init() {
	RootCmd.AddCommand(appCmd)

	appCmd.AddCommand(appGetCmd)
	addCommonGetFlags(appGetCmd)
}

func getApps(ids []string) []types.StructApp {
	c := Level27Client
	if len(ids) == 0 {
		return c.Apps(optGetParameters).Data
	} else {
		apps := make([]types.StructApp, len(ids))
		for idx, id := range ids {
			apps[idx] = c.App("GET", id, nil).Data
		}
		return apps
	}
}
