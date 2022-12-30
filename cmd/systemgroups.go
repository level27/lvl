package cmd

import (
	"errors"
	"fmt"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

// MAIN COMMAND
var systemgroupCmd = &cobra.Command{
	Use:   "systemgroup",
	Short: "Commands for managing systemgroups",
}

func init() {
	//adding main command to root
	RootCmd.AddCommand(systemgroupCmd)

	//------------------------------------------------- SYSTEMSGROUPS (GET / CREATE  / UPDATE / DELETE)-------------------------------------------------

	// --- GET (LIST)
	systemgroupCmd.AddCommand(systemgroupsGetCmd)
	// add optional get parameters (filters)
	addCommonGetFlags(systemgroupsGetCmd)

	// --- DESCRIBE
	systemgroupCmd.AddCommand(systemgroupDescribeCmd)

	// --- CREATE
	systemgroupCmd.AddCommand(systemgroupsCreateCmd)
	// flags for creating systemgroup
	flags := systemgroupsCreateCmd.Flags()
	flags.StringVarP(&systemgroupCreateName, "name", "n", "", "The name you want to give the systemgroup.")
	flags.StringVarP(&systemgroupCreateOrg, "organisation", "", "", "The name of the organisation this systemgroup belongs to.")
	systemgroupsCreateCmd.MarkFlagRequired("name")
	systemgroupsCreateCmd.MarkFlagRequired("organisation")

	// --- UPDATE
	systemgroupCmd.AddCommand(systemgroupsUpdateCmd)
	// flags for creating systemgroup
	flags = systemgroupsUpdateCmd.Flags()
	flags.StringVarP(&systemgroupUpdateName, "name", "n", "", "The name you want to give the systemgroup.")
	flags.StringVarP(&systemgroupUpdateOrg, "organisation", "", "", "The name of the organisation this systemgroup belongs to.")

	// --- DELETE
	systemgroupCmd.AddCommand(systemgroupsDeleteCmd)
	addDeleteConfirmFlag(systemgroupsDeleteCmd)
}

func resolveSystemgroup(arg string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.SystemgroupLookup(arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"systemgroup",
		func(group l27.Systemgroup) string { return fmt.Sprintf("%s (%d)", group.Name, group.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, err
}

// ------------------------------------------------- SYSTEMSGROUPS (GET / CREATE  / UPDATE / DELETE)-------------------------------------------------
// ---------------- DESCRIBE
var systemgroupDescribeCmd = &cobra.Command{
	Use: "describe",
	RunE: func(cmd *cobra.Command, args []string) error {
		//check for valid systemgroupID type
		systemgroupID, err := checkSingleIntID(args[0], "systemgroup")
		if err != nil {
			return err
		}

		systemgroup, err := Level27Client.SystemgroupsgetSingle(systemgroupID)
		if err != nil {
			return err
		}

		// create output on template
		outputFormatTemplate(systemgroup, "templates/systemgroup.tmpl")
		return nil
	},
}

// ---------------- GET
var systemgroupsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Show list of all available systemgroups.",
	RunE: func(cmd *cobra.Command, args []string) error {
		systemgroups, err := Level27Client.SystemgroupsGet(optGetParameters)
		if err != nil {
			return err
		}

		outputFormatTable(systemgroups, []string{"ID", "NAME", "ORGANISATION"}, []string{"ID", "Name", "Organisation.Name"})
		return nil
	},
}

// ---------------- CREATE
var systemgroupCreateName, systemgroupCreateOrg string
var systemgroupsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new systemgroup.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if systemgroupCreateName == "" {
			return errors.New("name cannot be empty")
		}

		organisationID, err := resolveOrganisation(systemgroupCreateOrg)
		if err != nil {
			return err
		}

		request := l27.SystemgroupRequest{
			Name:         systemgroupCreateName,
			Organisation: organisationID,
		}

		systemgroup, err := Level27Client.SystemgroupsCreate(request)
		if err != nil {
			return err
		}

		outputFormatTemplate(systemgroup, "templates/entities/systemgroup/create.tmpl")
		return nil
	},
}

// ---------------- UPDATE
var systemgroupUpdateName, systemgroupUpdateOrg string
var systemgroupsUpdateCmd = &cobra.Command{
	Use: "update",
	RunE: func(cmd *cobra.Command, args []string) error {
		systemgroupID, err := resolveSystemgroup(args[0])
		if err != nil {
			return err
		}

		// when no flag has been set. -> dont need to do an update
		if !cmd.Flag("name").Changed && !cmd.Flag("organisation").Changed {
			return nil
		}

		// get current data from the systemgroup
		currentData, err := Level27Client.SystemgroupsgetSingle(systemgroupID)
		if err != nil {
			return err
		}

		// fill in current data in request type (case only one thing has to change. the other one still needs to be send aswell (put))
		request := l27.SystemgroupRequest{
			Name:         currentData.Name,
			Organisation: currentData.Organisation.ID,
		}

		// when organisation flag is used
		if cmd.Flag("organisation").Changed {
			// this function accepts the organisation name (string)
			// and will look up the ID if the name is found
			organisationID, err := resolveOrganisation(systemgroupUpdateOrg)
			if err != nil {
				return err
			}

			request.Organisation = organisationID
		}

		// when name flag is used
		if cmd.Flag("name").Changed {
			// when given an empty string as name for systemgroup
			if len(systemgroupUpdateName) == 0 {
				return errors.New("name cannot be empty")
			}

			request.Name = systemgroupUpdateName
		}

		err = Level27Client.SystemgroupsUpdate(systemgroupID, request)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/systemgroup/update.tmpl")
		return nil
	},
}

// ---------------- DELETE
var systemgroupsDeleteCmd = &cobra.Command{
	Use: "delete",
	RunE: func(cmd *cobra.Command, args []string) error {
		systemgroupID, err := resolveSystemgroup(args[0])
		if err != nil {
			return err
		}

		if !optDeleteConfirmed {
			group, err := Level27Client.SystemgroupsgetSingle(systemgroupID)
			if err != nil {
				return err
			}

			if !confirmPrompt(fmt.Sprintf("Delete system group %s (%d)?", group.Name, group.ID)) {
				return nil
			}
		}

		err = Level27Client.SystemgroupDelete(systemgroupID)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/systemgroup/delete.tmpl")
		return nil
	},
}
