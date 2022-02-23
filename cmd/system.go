package cmd

import (
	"log"
	"strings"

	"bitbucket.org/level27/lvl/types"
	"github.com/spf13/cobra"
)

// MAIN COMMAND
var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Commands for managing systems",
}

func init() {
	//adding main command to root
	RootCmd.AddCommand(systemCmd)

	//-------------------------------------  Toplevel subcommands (get/post) --------------------------------------
	// --- GET
	systemCmd.AddCommand(systemGetCmd)
	addCommonGetFlags(systemGetCmd)

	// --- CREATE
	systemCmd.AddCommand(systemCreateCmd)
	flags := systemCreateCmd.Flags()
	flags.StringVarP(&systemCreateName, "name", "n", "", "The name you want to give the system")
	flags.StringVarP(&systemCreateFqdn, "Fqdn", "", "", "Valid hostname for the system")
	flags.StringVarP(&systemCreateRemarks, "remarks", "", "", "Remarks (Admin only)")
	flags.IntVarP(&systemCreateDisk, "disk", "", 0, "Disk (non-editable)")
	flags.IntVarP(&systemCreateCpu, "cpu", "", 0, "Cpu (Required for Level27 systems)")
	flags.IntVarP(&systemCreateMemory, "memory", "", 0, "Memory (Required for Level27 systems)")
	flags.StringVarP(&systemCreateManageType, "management", "", "basic", "Managament type (default: basic)")
	flags.BoolVarP(&systemCreatePublicNetworking, "publicNetworking", "", true, "For digitalOcean servers always true. (non-editable)")
	flags.IntVarP(&systemCreateImage, "image", "", 0, "The ID of a systemimage. (must match selected configuration and zone. non-editable)")
	flags.IntVarP(&systemCreateOrganisation, "organisation", "", 0, "The unique ID of an organisation")
	flags.IntVarP(&systemCreateProviderConfig, "provider", "", 0, "The unique ID of a SystemproviderConfiguration")
	flags.IntVarP(&systemCreateZone, "zone", "", 0, "The unique ID of a zone")
	flags.StringVarP(&systemCreateSecurityUpdates, "security", "", "", "installSecurityUpdates (default: random POST:1-8, PUT:0-12)")
	flags.StringVarP(&systemCreateAutoTeams, "autoTeams", "", "", "A csv list of team ID's")
	flags.StringVarP(&systemCreateExternalInfo, "externalInfo", "", "", "ExternalInfo (required when billableItemInfo entities for an organisation exist in db)")
	flags.IntVarP(&systemCreateOperatingSystemVersion, "version", "", 0, "The unique ID of an OperatingsystemVersion (non-editable)")
	flags.IntVarP(&systemCreateParentSystem, "parent", "", 0, "The unique ID of a system (parent system)")
	flags.StringVarP(&systemCreateType, "type", "", "", "System type")
	flags.StringArrayP("networks", "", []string{""}, "Array of network IP's. (default: null)")

	// Required flags for create system.
	requiredFlags := []string{"name", "image", "organisation", "provider", "zone"}
	for _, flag := range requiredFlags {
		systemCreateCmd.MarkFlagRequired(flag)
	}

	//-------------------------------------  SYSTEMS/CHECKS TOPLEVEL (get/post) --------------------------------------
	systemCmd.AddCommand(systemCheckCmd)
	// GET LIST OF ALL CHECKS
	systemCheckCmd.AddCommand(systemCheckGetCmd)
	addCommonGetFlags(systemGetCmd)
}

//------------------------------------------------- SYSTEM TOPLEVEL (GET / CREATE) ----------------------------------
//----------------------------------------- GET ---------------------------------------
var systemGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get a list of all curent systems",
	Run: func(cmd *cobra.Command, args []string) {
		ids, err := convertStringsToIds(args)
		if err != nil {
			log.Fatalln("Invalid domain ID")
		}
		outputFormatTable(getSystems(ids), []string{"ID", "NAME", "STATUS"}, []string{"Id", "Name", "Status"})

	},
}

func getSystems(ids []int) []types.System {

	if len(ids) == 0 {
		return Level27Client.SystemGetList(optGetParameters)
	} else {
		systems := make([]types.System, len(ids))
		for idx, id := range ids {
			systems[idx] = Level27Client.SystemGetSingle(id)
		}
		return systems
	}

}

//----------------------------------------- CREATE ---------------------------------------
// vars needed to save flag data.
var systemCreateName, systemCreateFqdn, systemCreateRemarks string
var systemCreateDisk, systemCreateCpu, systemCreateMemory int
var systemCreateManageType string
var systemCreatePublicNetworking bool
var systemCreateImage, systemCreateOrganisation, systemCreateProviderConfig, systemCreateZone int
var systemCreateSecurityUpdates string
var systemCreateAutoTeams, systemCreateExternalInfo string
var systemCreateOperatingSystemVersion, systemCreateParentSystem int
var systemCreateType string
var systemCreateAutoNetworks []interface{}
var managementTypeArray = []string{"basic", "professional", "enterprise", "professional_level27"}
// var securityUpdatesArray = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

var systemCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new system",
	Run: func(cmd *cobra.Command, args []string) {

		managementTypeValue := cmd.Flag("management").Value.String()
		// securityUpdateValue := cmd.Flag("security").Value.String()
		// var checkedSecurityUpdateValue int

		//  checking if the management flag has been changed/set
		if cmd.Flag("management").Changed {

			// checking if given managamentType is one of the possible options.
			var isValidManagementType bool
			for _, arrayItem := range managementTypeArray {
				if strings.ToLower(managementTypeValue) == arrayItem {
					managementTypeValue = arrayItem
					isValidManagementType = true
				}
			}
			// if no valid management type was given -> error for user
			if !isValidManagementType {
				log.Printf("ERROR: given managementType is not valid: '%v'", managementTypeValue)
			}
		}

		// Using data from the flags to make the right type used for posting a new system. (types systemPost)
		RequestData := types.SystemPost{
			Name:                        systemCreateName,
			CustomerFqdn:                systemCreateFqdn,
			Remarks:                     systemCreateRemarks,
			Disk:                        &systemCreateDisk,
			Cpu:                         &systemCreateCpu,
			Memory:                      &systemCreateMemory,
			MamanagementType:            managementTypeValue,
			PublicNetworking:            systemCreatePublicNetworking,
			SystemImage:                 systemCreateImage,
			Organisation:                systemCreateOrganisation,
			SystemProviderConfiguration: systemCreateProviderConfig,
			Zone:                        systemCreateZone,
			// InstallSecurityUpdates:      &checkedSecurityUpdateValue,
			AutoTeams:              systemCreateAutoTeams,
			ExternalInfo:           systemCreateExternalInfo,
			OperatingSystemVersion: &systemCreateOperatingSystemVersion,
			ParentSystem:           &systemCreateParentSystem,
			Type:                   systemCreateType,
			AutoNetworks:           systemCreateAutoNetworks,
		}

		// checking if the install securityUpdates flag is changed
		if cmd.Flag("security").Changed {

			// check if given value is a valid int
			// checkedSecurityUpdateValue, err := strconv.Atoi(securityUpdateValue)
			// if err != nil {
			// 	log.Printf("ERROR: given installSecurityUpdates not valid: '%v'", securityUpdateValue)
			// } else {
			// 	//if given value is okay, check if value is one of the allowed values
			// 	var isValidSecurityUpdate bool
			// 	for _, arrayItem := range securityUpdatesArray {
			// 		if checkedSecurityUpdateValue == arrayItem {

			// 			RequestData.InstallSecurityUpdates = new(int)
			// 			RequestData.InstallSecurityUpdates = &arrayItem
			// 			log.Printf("value gevonden: '%v'", arrayItem)

			// 			isValidSecurityUpdate = true
			// 		}
			// 	}

			// 	if !isValidSecurityUpdate {
			// 		log.Printf("ERROR: given installSecurityUpdates is not valid: '%v'", securityUpdateValue)
			// 		RequestData.InstallSecurityUpdates = nil
			// 	}
			// }

		}

		if *RequestData.Disk == 0 {
			RequestData.Disk = nil
		}

		if *RequestData.Cpu == 0 {
			RequestData.Cpu = nil
		}

		if *RequestData.Memory == 0 {
			RequestData.Memory = nil
		}

		if *RequestData.OperatingSystemVersion == 0 {
			RequestData.OperatingSystemVersion = nil
		}

		if *RequestData.ParentSystem == 0 {
			RequestData.ParentSystem = nil
		}
		Level27Client.SystemCreate(args, RequestData)

	},
}

//------------------------------------------------- SYSTEM/CHECKS TOPLEVEL (GET / CREATE) ----------------------------------
// ---------------- MAIN COMMAND (checks)
var systemCheckCmd = &cobra.Command{
	Use:   "checks",
	Short: "Command for managing systems checks",
}

// ---------------- GET

var systemCheckGetCmd = &cobra.Command{
	Use:   "checks",
	Short: "Command for managing systems checks",
}
