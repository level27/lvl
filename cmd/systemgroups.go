package cmd

import (
	"fmt"
	"log"
	"strconv"

	"bitbucket.org/level27/lvl/types"
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
	//flag to skip confirmation when deleting a systemgroup
	systemgroupsDeleteCmd.Flags().BoolVarP(&systemgroupDeleteConfirmed, "yes", "y", false, "Set this flag to skip confirmation when deleting a systemgroup")

}

func resolveSystemgroup(arg string) int {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id
	}

	return resolveShared(
		Level27Client.SystemgroupLookup(arg),
		arg,
		"systemgroup",
		func(group types.Systemgroup) string { return fmt.Sprintf("%s (%d)", group.Name, group.ID) }).ID
}


//------------------------------------------------- SYSTEMSGROUPS (GET / CREATE  / UPDATE / DELETE)-------------------------------------------------
// ---------------- DESCRIBE
var systemgroupDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Get detailed info about a systemgroup",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//check for valid systemgroupId type
		systemgroupID := checkSingleIntID(args[0], "systemgroup")

		systemgroup := Level27Client.SystemgroupsgetSingle(systemgroupID)

		// create output on template
		outputFormatTemplate(systemgroup, "templates/systemgroup.tmpl")
	},
}

// ---------------- GET
var systemgroupsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Show list of all available systemgroups.",
	Run: func(cmd *cobra.Command, args []string) {

		systemgroups := Level27Client.SystemgroupsGet(optGetParameters)

		outputFormatTable(systemgroups, []string{"ID", "NAME", "ORGANISATION"}, []string{"ID", "Name", "Organisation.Name"})

	},
}

// ---------------- CREATE
var systemgroupCreateName, systemgroupCreateOrg string
var systemgroupsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new systemgroup.",
	Run: func(cmd *cobra.Command, args []string) {

		// this function accepts the organisation name (string)
		// and will look up the ID if the name is found
		organisationId := resolveOrganisation(systemgroupCreateOrg)

		// when given an empty string as name
		if len(systemgroupCreateName) == 0 {
			cobra.CheckErr("Name cannot be empty!")
		}

		// fill in given data in request type
		request := types.SystemgroupRequest{
			Name:         systemgroupCreateName,
			Organisation: organisationId,
		}

		systemgroup := Level27Client.SystemgroupsCreate(request)

		// will only print if systemgroup is created successfully
		log.Printf("systemgroup succesfully created. [ID: '%v' - NAME: '%v'].", systemgroup.ID, systemgroup.Name)

	},
}

// ---------------- UPDATE
var systemgroupUpdateName, systemgroupUpdateOrg string
var systemgroupsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update a systemgroups name or organisation.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		//check for valid systemgroupId type
		systemgroupID := checkSingleIntID(args[0], "systemgroup")

		// when no flag has been set. -> dont need to do an update
		if !cmd.Flag("name").Changed && !cmd.Flag("organisation").Changed {
			cobra.CheckErr("Use at least one flag to change a value of the systemgroup.")
		} else {
			// get current data from the systemgroup
			currentData := Level27Client.SystemgroupsgetSingle(systemgroupID)

			// fill in current data in request type (case only one thing has to change. the other one still needs to be send aswell (put))
			request := types.SystemgroupRequest{
				Name:         currentData.Name,
				Organisation: currentData.Organisation.ID,
			}

			// when organisation flag is used
			if cmd.Flag("organisation").Changed {
				// this function accepts the organisation name (string)
				// and will look up the ID if the name is found
				organisationId := resolveOrganisation(systemgroupUpdateOrg)
				request.Organisation = organisationId
			}

			// when name flag is used
			if cmd.Flag("name").Changed {
				// when given an empty string as name for systemgroup
				if len(systemgroupUpdateName) == 0 {
					cobra.CheckErr("Name cannot be empty!")
				} else {
					request.Name = systemgroupUpdateName
				}
			}

			Level27Client.SystemgroupsUpdate(systemgroupID, request)

		}

	},
}

// ---------------- DELETE
var systemgroupDeleteConfirmed bool
var systemgroupsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a systemgroup",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		//check for valid systemgroupId type
		systemgroupId := checkSingleIntID(args[0], "systemgroup")

		Level27Client.SystemgroupDelete(systemgroupId, systemgroupDeleteConfirmed)

	},
}
