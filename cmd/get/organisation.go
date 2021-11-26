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

var organisationGetCmd = &cobra.Command{
	Use:   "organisation [IDs to retrieve]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ArbitraryArgs,
	Run: func(ccmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\t")

		organisations := getOrganisations(args)
		for _, organisation := range organisations {
			fmt.Fprintln(w, strconv.Itoa(organisation.ID)+"\t"+organisation.Name+"\t")
		}
	
		w.Flush()
	},
}

func init() {
	GetCmd.AddCommand(organisationGetCmd)
}

func getOrganisations(ids []string) []types.StructOrganisation {
	c := cmd.Level27Client
	if len(ids) == 0 {
		numberToGet := viper.GetString("number")
		organisationFilter := viper.GetString("filter")
	
		return c.Organisations(organisationFilter, numberToGet).Data
	} else {
		organisations := make([]types.StructOrganisation, len(ids))
		for idx, id := range ids {
			organisations[idx] = c.Organisation("GET", id, nil).Data
		}
		return organisations
	}
}
