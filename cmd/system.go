package cmd

import (
	"log"
	"strconv"
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

	// --- DESCRIBE
	systemCmd.AddCommand(systemDescribeCmd)
	systemDescribeCmd.Flags().BoolVar(&systemDescribeHideJobs, "hide-jobs", false, "Hide jobs in the describe output.")

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
	//	flags.StringVarP(&systemCreateSecurityUpdates, "security", "", "", "installSecurityUpdates (default: random POST:1-8, PUT:0-12)") NOT NEEDED FOR CREATE REQUEST
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

	// --- ACTIONS
	systemCmd.AddCommand(systemActionsCmd)

	systemActionsCmd.AddCommand(systemActionsStartCmd)
	systemActionsCmd.AddCommand(systemActionsStopCmd)
	systemActionsCmd.AddCommand(systemActionsShutdownCmd)
	systemActionsCmd.AddCommand(systemActionsRebootCmd)
	systemActionsCmd.AddCommand(systemActionsResetCmd)
	systemActionsCmd.AddCommand(systemActionsEmergencyPowerOffCmd)
	systemActionsCmd.AddCommand(systemActionsDeactivateCmd)
	systemActionsCmd.AddCommand(systemActionsActivateCmd)
	systemActionsCmd.AddCommand(systemActionsAutoInstallCmd)

	//-------------------------------------  SYSTEMS/CHECKS TOPLEVEL (get/post) --------------------------------------
	systemCmd.AddCommand(systemCheckCmd)
	// ---- GET LIST OF ALL CHECKS
	systemCheckCmd.AddCommand(systemCheckGetCmd)
	addCommonGetFlags(systemCheckGetCmd)

	// ---- CREATE NEW CHECK
	systemCheckCmd.AddCommand(systemCheckCreateCmd)

	// -- flags needed to create a check
	flags = systemCheckCreateCmd.Flags()
	flags.StringVarP(&systemCheckCreate, "type", "t", "", "Check type (non-editable)")
	systemCheckCreateCmd.MarkFlagRequired("type")

	// -- optional flags, only for creating a http check
	flags.IntVarP(&systemCreateCheckPort, "port", "p", 80, "Port for http checktype.")
	flags.StringVarP(&systemCreateCheckHost, "host", "", "", "Hostname for http checktype.")
	flags.StringVarP(&systemCreateCheckUrl, "url", "", "", "Url for http checktype.")
	flags.StringVarP(&systemCreateCheckContent, "content", "c", "", "Content for http checktype.")


	//-------------------------------------  SYSTEMS/CHECKS ACTIONS (get/ delete/ update) --------------------------------------
	// --- DESCRIBE CHECK
	systemCheckCmd.AddCommand(systemCheckGetSingleCmd)
	// --- DELETE CHECK
	systemCheckCmd.AddCommand(systemCheckDeleteCmd)

	//flag to skip confirmation when deleting a check
	systemCheckDeleteCmd.Flags().BoolVarP(&systemCheckDeleteConfirmed, "yes", "y", false, "Set this flag to skip confirmation when deleting a check")



	//-------------------------------------  SYSTEMS/COOKBOOKS TOPLEVEL (get/post) --------------------------------------
	// adding cookbook subcommand to system command
	systemCmd.AddCommand(systemCookbookCmd)

	// ---- GET cookbooks
	systemCookbookCmd.AddCommand(systemCookbookGetCmd)

}

//------------------------------------------------- SYSTEM TOPLEVEL (GET / CREATE) ----------------------------------
//----------------------------------------- GET ---------------------------------------
var systemGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get a list of all curent systems",
	Run: func(cmd *cobra.Command, args []string) {
		ids, err := convertStringsToIds(args)
		if err != nil {
			log.Fatalln("Invalid system ID")
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

//----------------------------------------- DESCRIBE ---------------------------------------
var systemDescribeHideJobs = false

var systemDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Get detailed information about a system.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		systemID, err := convertStringToId(args[0])
		if err != nil {
			log.Fatalln("Invalid system ID")
		}

		var system types.DescribeSystem
		system.System = Level27Client.SystemGetSingle(systemID)
		if !systemDescribeHideJobs {
			system.Jobs = Level27Client.EntityJobHistoryGet("system", systemID)
			for idx, j := range system.Jobs {
				system.Jobs[idx] = Level27Client.JobHistoryRootGet(j.Id)
			}
		}

		system.SshKeys = Level27Client.SystemGetSshKeys(systemID, types.CommonGetParams{})
		securityUpdates := Level27Client.SecurityUpdateDates()
		system.InstallSecurityUpdatesString = securityUpdates[system.InstallSecurityUpdates]
		system.HasNetworks = Level27Client.SystemGetHasNetworks(systemID)
		system.Volumes = Level27Client.SystemGetVolumes(systemID, types.CommonGetParams{})

		outputFormatTemplate(system, "templates/system.tmpl")
	},
}

//----------------------------------------- CREATE ---------------------------------------
// vars needed to save flag data.
var systemCreateName, systemCreateFqdn, systemCreateRemarks string
var systemCreateDisk, systemCreateCpu, systemCreateMemory int
var systemCreateManageType string
var systemCreatePublicNetworking bool
var systemCreateImage, systemCreateOrganisation, systemCreateProviderConfig, systemCreateZone int

var systemCreateAutoTeams, systemCreateExternalInfo string
var systemCreateOperatingSystemVersion, systemCreateParentSystem int
var systemCreateType string
var systemCreateAutoNetworks []interface{}

// ARRAY NOG DYNAMIC MAKEN!!!!!
var managementTypeArray = []string{"basic", "professional", "enterprise", "professional_level27"}

// var securityUpdatesArray = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}        - not needed for create request
// var systemCreateSecurityUpdates string 											/

var systemCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new system",
	Run: func(cmd *cobra.Command, args []string) {

		managementTypeValue := cmd.Flag("management").Value.String()

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
			// InstallSecurityUpdates:      &checkedSecurityUpdateValue, NOT NEEDED IN CREATE REQUEST//
			AutoTeams:              systemCreateAutoTeams,
			ExternalInfo:           systemCreateExternalInfo,
			OperatingSystemVersion: &systemCreateOperatingSystemVersion,
			ParentSystem:           &systemCreateParentSystem,
			Type:                   systemCreateType,
			AutoNetworks:           systemCreateAutoNetworks,
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
		Level27Client.SystemCreate(RequestData)

	},
}

//------------------------------------------------- SYSTEM/CHECKS TOPLEVEL (GET / CREATE) ----------------------------------
// ---------------- MAIN COMMAND (checks)
var systemCheckCmd = &cobra.Command{
	Use:   "checks",
	Short: "Manage systems checks",
}

// ---------------- GET

var systemCheckGetCmd = &cobra.Command{
	Use:   "get [system ID]",
	Short: "Get a list of all checks for a system",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// check for valid system ID
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalln("Not a valid system ID!")
		}

		// Creating readable output
		outputFormatTable(getSystemChecks(id), []string{"ID", "CHECKTYPE", "STATUS", "INFORMATION"}, []string{"Id", "CheckType", "Status", "StatusInformation"})

	},
}

func getSystemChecks(id int) []types.SystemCheck {

	return Level27Client.SystemCheckGetList(id, optGetParameters)

}

// ---------------- CREATE CHECK
var systemCheckCreate, systemCreateCheckUrl, systemCreateCheckContent, systemCreateCheckHost string
var systemCreateCheckPort int

var systemCheckCreateCmd = &cobra.Command{
	Use:   "create [system ID] [parameters]",
	Short: "create a new check for a specific system",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//check for valid system ID
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalln("Not a valid system ID!")
		}

		// get the value of the flag type set by user
		checkTypeInput := cmd.Flag("type").Value.String()

		// bool value to see if user input is valid, and bool to check if chosen type is http
		var isChecktypeValid, isCheckTypeHttp bool
		if cmd.Flag("type").Changed {

			// GET REQUEST to see what all curent valid checktypes are (function gives back an array of valid types)
			systemCheckCreateArray := Level27Client.SystemCheckTypeGet()

			//when user input is one of the valid options -> validation bool is true
			for _, validOption := range systemCheckCreateArray {
				if strings.ToLower(checkTypeInput) == validOption {

					checkTypeInput = validOption

					// check if chosen type is http
					if checkTypeInput == "http" {
						isCheckTypeHttp = true
					}
					isChecktypeValid = true

				}
			}
			// if user input not in valid options array -> error
			if !isChecktypeValid {
				log.Fatalln("Given checktype is not valid")
			} else {
				//when user chose http type, aditional flags can be set
				if isCheckTypeHttp {
					request := types.SystemCheckRequestHttp{
						Checktype: checkTypeInput,
						Port:      systemCreateCheckPort,
						Url:       systemCreateCheckUrl,
						Hostname:  systemCreateCheckHost,
						Content:   systemCreateCheckContent,
					}
					Level27Client.SystemCheckCreate(id, request)
					//when chosen type NOT http -> only checktype will be needed for request
				} else {
					request := types.SystemCheckRequest{
						Checktype: checkTypeInput,
					}
					Level27Client.SystemCheckCreate(id, request)
				}

			}

		}

	},
}
//------------------------------------------------- SYSTEM/CHECKS ACTIONS (GET / DELETE / UPDATE) ----------------------------------
// -------------- GET DETAILS FROM A CHECK
var systemCheckGetSingleCmd = &cobra.Command{
	Use: "describe [systemID] [checkID]",
	Short: "Get detailed info about a specific check.",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		//check for valid system ID
		systemID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalln("Not a valid system ID!")
		}

		//check for valid system checkID
		checkID, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalln("Not a valid check ID!")
		}

		check := Level27Client.SystemCheckDescribe(systemID, checkID)

		outputFormatTemplate(check, "templates/systemCheck.tmpl")
	},
}
// -------------- DELETE SPECIFIC CHECK
var systemCheckDeleteConfirmed bool
var systemCheckDeleteCmd = &cobra.Command{
	Use: "delete [systemID] [checkID]",
	Short: "Delete a specific check from a system",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		//check for valid system ID
		systemID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalln("Not a valid system ID!")
		}

		//check for valid system checkID
		checkID, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalln("Not a valid check ID!")
		}

		Level27Client.SystemCheckDelete(systemID, checkID, systemCheckDeleteConfirmed)
	},
}


//------------------------------------------------- SYSTEM/COOKBOOKS TOPLEVEL (GET / CREATE) ----------------------------------
// ---------------- MAIN COMMAND (checks)
var systemCookbookCmd = &cobra.Command{
	Use:   "cookbook",
	Short: "Manage systems cookbooks",
}

// ---------- GET COOKBOOKS
var systemCookbookGetCmd = &cobra.Command{
	Use:   "get [system ID]",
	Short: "Gets a list of all cookbooks from a system.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalln("Not a valid system ID!")
		}

		outputFormatTable(getSystemCookbooks(id), []string{"ID", "CHECKTYPE", "STATUS"}, []string{"Id", "Checktype", "Status"})
	},
}

func getSystemCookbooks(id int) []types.Cookbook {

	return Level27Client.SystemCookbookGetList(id)

}

// -------- SYSTEM ACTIONS

var systemActionsCmd = &cobra.Command{
	Use: "actions",
	Short: "Actions for systems such as rebooting",
}

var systemActionsStartCmd = &cobra.Command{
	Use: "start",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) { runAction("start", args) },
}

var systemActionsStopCmd = &cobra.Command{
	Use: "stop",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) { runAction("stop", args) },
}

var systemActionsShutdownCmd = &cobra.Command{
	Use: "shutdown",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) { runAction("shutdown", args) },
}

var systemActionsRebootCmd = &cobra.Command{
	Use: "reboot",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) { runAction("reboot", args) },
}

var systemActionsResetCmd = &cobra.Command{
	Use: "reset",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) { runAction("reset", args) },
}

var systemActionsEmergencyPowerOffCmd = &cobra.Command{
	Use: "emergencyPowerOff",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) { runAction("emergencyPowerOff", args) },
}

var systemActionsDeactivateCmd = &cobra.Command{
	Use: "deactivate",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) { runAction("deactivate", args) },
}

var systemActionsActivateCmd = &cobra.Command{
	Use: "activate",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) { runAction("activate", args) },
}

var systemActionsAutoInstallCmd = &cobra.Command{
	Use: "autoInstall",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) { runAction("autoInstall", args) },
}

func runAction(action string, args []string) {
	id, err := convertStringToId(args[0])
	if err != nil {
		log.Fatalln("Invalid system ID")
	}

	Level27Client.SystemAction(id, action)
}