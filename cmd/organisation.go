package cmd

import (
	"log"

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
		ids, err := convertStringsToIds(args)
		if err != nil {
			log.Fatalln("Invalid organisation ID")
		}
		outputFormatTable(getOrganisations(ids), []string{"ID", "NAME"}, []string{"ID", "Name"})
	},
}

func init() {
	RootCmd.AddCommand(organisationCmd)

	organisationCmd.AddCommand(organisationGetCmd)
	addCommonGetFlags(organisationGetCmd)
}

func getOrganisations(ids []int) []types.Organisation {
	c := Level27Client
	if len(ids) == 0 {
		return c.Organisations(optGetParameters)
	} else {
		organisations := make([]types.Organisation, len(ids))
		for idx, id := range ids {
			organisations[idx] = c.Organisation(id)
		}
		return organisations
	}
}
