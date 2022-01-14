package cmd

import (
	"bitbucket.org/level27/lvl/types"
	"github.com/spf13/cobra"
)

var organisationCmd = &cobra.Command{
	Use:   "organisation",
	Short: "Commands for managing organisations",
}

var organisationGetCmd = &cobra.Command{
	Use: "get",

	Args: cobra.ArbitraryArgs,
	Run: func(ccmd *cobra.Command, args []string) {
		outputFormatTable(getOrganisations(args), []string{"ID", "NAME"}, []string{"ID", "Name"})
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
		return c.Organisations(optGetParameters).Data
	} else {
		organisations := make([]types.StructOrganisation, len(ids))
		for idx, id := range ids {
			organisations[idx] = c.Organisation("GET", id, nil).Data
		}
		return organisations
	}
}
