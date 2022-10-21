package cmd

import (
	"fmt"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

var organisationCmd = &cobra.Command{
	Use:   "organisation",
	Short: "Commands for managing organisations",
}

var organisationGetCmd = &cobra.Command{
	Use: "get",

	Args: cobra.ArbitraryArgs,
	RunE: func(ccmd *cobra.Command, args []string) error {
		ids, err := convertStringsToIds(args)
		if err != nil {
			return err
		}

		options, err := getOrganisations(ids)
		if err != nil {
			return err
		}

		outputFormatTable(options, []string{"ID", "NAME"}, []string{"ID", "Name"})
		return nil
	},
}

func init() {
	RootCmd.AddCommand(organisationCmd)

	organisationCmd.AddCommand(organisationGetCmd)
	addCommonGetFlags(organisationGetCmd)
}

func resolveOrganisation(arg string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.LookupOrganisation(arg)
	if err != nil {
		return id, nil
	}

	res, err := resolveShared(
		options,
		arg,
		"organisation",
		func(app l27.Organisation) string { return fmt.Sprintf("%s (%d)", app.Name, app.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, nil
}

func getOrganisations(ids []l27.IntID) ([]l27.Organisation, error) {
	c := Level27Client
	if len(ids) == 0 {
		return c.Organisations(optGetParameters)
	} else {
		organisations := make([]l27.Organisation, len(ids))
		for idx, id := range ids {
			var err error
			organisations[idx], err = c.Organisation(id)
			if err != nil {
				return nil, err
			}
		}

		return organisations, nil
	}
}
