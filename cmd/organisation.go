package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"bitbucket.org/level27/lvl/types"
	"github.com/spf13/cobra"
)

var organisationCmd = &cobra.Command{
	Use: "organisation",
	Short: "Commands for managing organisations",
}

var organisationGetCmd = &cobra.Command{
	Use:   "get",
	
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
	RootCmd.AddCommand(organisationCmd)

	organisationCmd.AddCommand(organisationGetCmd)
	addCommonGetFlags(organisationGetCmd)
}

func getOrganisations(ids []string) []types.StructOrganisation {
	c := Level27Client
	if len(ids) == 0 {
		return c.Organisations(optFilter, optNumber).Data
	} else {
		organisations := make([]types.StructOrganisation, len(ids))
		for idx, id := range ids {
			organisations[idx] = c.Organisation("GET", id, nil).Data
		}
		return organisations
	}
}
