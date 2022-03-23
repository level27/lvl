package cmd

import (
	"log"

	"bitbucket.org/level27/lvl/types"
	"github.com/spf13/cobra"
)

// MAIN COMMAND
var systemgroupCmd = &cobra.Command{
	Use:   "systemgroups",
	Short: "Commands for managing systemgroups",
}

func init() {
	//adding main command to root
	RootCmd.AddCommand(systemgroupCmd)

	// --- GET (LIST)
	systemgroupCmd.AddCommand(systemgroupsGetCmd)
	// add optional get parameters (filters)
	addCommonGetFlags(systemgroupsGetCmd)

	// --- CREATE 
	systemgroupCmd.AddCommand(systemgroupsCreateCmd)
	// flags for creating systemgroup
	flags := systemgroupsCreateCmd.Flags()
	flags.StringVarP(&systemgroupCreateName, "name" , "n", "", "The name you want to give the systemgroup.")
	flags.StringVarP(&systemgroupCreateOrg, "organisation", "", "", "The name of the organisation this systemgroup belongs to.")
	systemgroupsCreateCmd.MarkFlagRequired("name")
	systemgroupsCreateCmd.MarkFlagRequired("organisation")

}

//------------------------------------------------- SYSTEMSGROUPS (GET / CREATE  / UPDATE / DELETE)-------------------------------------------------
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
	Use: "create",
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
			Name: systemgroupCreateName,
			Organisation: organisationId,
		}

		systemgroup := Level27Client.SystemgroupsCreate(request)

		// will only print if systemgroup is created successfully
		log.Printf("systemgroup succesfully created. [ID: '%v' - NAME: '%v'].", systemgroup.ID, systemgroup.Name)


	},
}