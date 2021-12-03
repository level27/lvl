package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"bitbucket.org/level27/lvl/types"
	"github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
	Use: "app",
	Short: "Commands to manage apps",
}

var appGetCmd = &cobra.Command{
	Use:   "get",
	Args: cobra.ArbitraryArgs,
	Run: func(acmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\t")

		apps := getApps(args)
		for _, app := range apps {
			fmt.Fprintln(w, strconv.Itoa(app.ID)+"\t"+app.Name+"\t"+app.Status+"\t")
		}
	
		w.Flush()
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
		return c.Apps(optFilter, optNumber).Data
	} else {
		apps := make([]types.StructApp, len(ids))
		for idx, id := range ids {
			apps[idx] = c.App("GET", id, nil).Data
		}
		return apps
	}
}
