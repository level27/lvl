package get

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"bitbucket.org/level27/lvl/cmd"
	"bitbucket.org/level27/lvl/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var appGetCmd = &cobra.Command{
	Use:   "app [IDs to retrieve]",
	Short: "Get available apps",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	GetCmd.AddCommand(appGetCmd)
}

func getApps(ids []string) []types.StructApp {
	c := cmd.Level27Client
	if len(ids) == 0 {
		numberToGet := viper.GetString("number")
		appFilter := viper.GetString("filter")
	
		return c.Apps(appFilter, numberToGet).Data
	} else {
		apps := make([]types.StructApp, len(ids))
		for idx, id := range ids {
			apps[idx] = c.App("GET", id, nil).Data
		}
		return apps
	}
}
