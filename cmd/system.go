package cmd

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"bitbucket.org/level27/lvl/types"
	"bitbucket.org/level27/lvl/utils"
	"github.com/Jeffail/gabs/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// MAIN COMMAND
var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Commands for managing systems",
}

func init() {
	//adding main command to root
	RootCmd.AddCommand(systemCmd)

	//-------------------------------------  Toplevel SYSTEM COMMANDS (get/post) --------------------------------------
	// #region Toplevel SYSTEM COMMANDS (get/post)

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
	flags.StringVarP(&systemCreateImage, "image", "", "", "The ID of a systemimage. (must match selected configuration and zone. non-editable)")
	flags.StringVarP(&systemCreateOrganisation, "organisation", "", "", "The unique ID of an organisation")
	flags.StringVarP(&systemCreateProviderConfig, "provider", "", "", "The unique ID of a SystemproviderConfiguration")
	flags.StringVarP(&systemCreateZone, "zone", "", "", "The unique ID of a zone")
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
	// #endregion

	// ------------------------------------ ACTIONS ON SPECIFIC SYSTEM ----------------------------------------------

	// #region ACTIONS ON SPECIFIC SYSTEM
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

	// --- UPDATE

	systemCmd.AddCommand(systemUpdateCmd)
	systemUpdateCmd.Flags().StringVarP(&systemUpdateSettingsFile, "settings-file", "s", "", "JSON file to read settings from. Pass '-' to read from stdin.")
	settingString(systemUpdateCmd, systemUpdateSettings, "name", "New name for this system")
	settingInt(systemUpdateCmd, systemUpdateSettings, "cpu", "Set amount of CPU cores of the system")
	settingInt(systemUpdateCmd, systemUpdateSettings, "memory", "Set amount of memory in GB of the system")
	settingString(systemUpdateCmd, systemUpdateSettings, "managementType", "Set management type of the system")
	settingString(systemUpdateCmd, systemUpdateSettings, "organisation", "Set organisation that owns this system. Can be both a name or an ID")
	settingInt(systemUpdateCmd, systemUpdateSettings, "publicNetworking", "")
	settingInt(systemUpdateCmd, systemUpdateSettings, "limitRiops", "Set read IOPS limit")
	settingInt(systemUpdateCmd, systemUpdateSettings, "limitWiops", "Set write IOPS limit")
	settingInt(systemUpdateCmd, systemUpdateSettings, "installSecurityUpdates", "Set security updates mode index")
	settingString(systemUpdateCmd, systemUpdateSettings, "remarks", "")

	// --- Delete

	systemCmd.AddCommand(systemDeleteCmd)
	systemDeleteCmd.Flags().BoolVar(&systemDeleteForce, "force", false, "")
	// #endregion

	//-------------------------------------  SYSTEMS/CHECKS TOPLEVEL (get/post) --------------------------------------

	// #region SYSTEMS/CHECKS TOPLEVEL
	systemCmd.AddCommand(systemCheckCmd)
	// ---- GET LIST OF ALL CHECKS
	systemCheckCmd.AddCommand(systemCheckGetCmd)
	addCommonGetFlags(systemCheckGetCmd)

	// ---- CREATE NEW CHECK
	systemCheckCmd.AddCommand(systemCheckCreateCmd)

	// -- flags needed to create a check
	flags = systemCheckCreateCmd.Flags()
	flags.StringVarP(&systemCheckCreateType, "type", "t", "", "Check type (non-editable)")
	systemCheckCreateCmd.MarkFlagRequired("type")

	// -- optional flags, only for creating a http check
	flags.StringArrayVarP(&systemDynamicParams, "parameters", "p", systemDynamicParams, "Add custom parameters for cookbook. SINGLE PAR: [ -p waf=true ], MULTIPLE PAR: [ -p waf=true -p timeout=200 ], MULTIPLE VALUES: [ -p versions=''7, 5.4'']")

	// ---- GET PARAMETERS (for specific checktype)
	systemCheckCmd.AddCommand(systemChecktypeParametersGetCmd)

	// flags needed to get checktype parameters
	systemChecktypeParametersGetCmd.Flags().StringVarP(&systemCheckCreateType, "type", "t", "", "Check type to see all its available parameters")

	//-------------------------------------  SYSTEMS/CHECKS ACTIONS (get/ delete/ update) --------------------------------------
	// #region SYSTEMS/CHECKS ACTIONS
	// --- DESCRIBE CHECK
	systemCheckCmd.AddCommand(systemCheckGetSingleCmd)
	// --- DELETE CHECK
	systemCheckCmd.AddCommand(systemCheckDeleteCmd)

	//flag to skip confirmation when deleting a check
	systemCheckDeleteCmd.Flags().BoolVarP(&systemCheckDeleteConfirmed, "yes", "y", false, "Set this flag to skip confirmation when deleting a check")

	// --- UPDATE CHECK (ONLY FOR HTTP REQUEST)
	systemCheckCmd.AddCommand(systemCheckUpdateCmd)

	// #endregion

	//-------------------------------------  SYSTEMS/COOKBOOKS TOPLEVEL (get/post) --------------------------------------

	// #region SYSTEMS/COOKBOOKS TOPLEVEL (get/post)

	// adding cookbook subcommand to system command
	systemCmd.AddCommand(systemCookbookCmd)

	// ---- GET cookbooks
	systemCookbookCmd.AddCommand(systemCookbookGetCmd)

	// ---- ADD cookbook (to system)
	systemCookbookCmd.AddCommand(systemCookbookAddCmd)

	// flags needed to add new cookbook to a system
	flags = systemCookbookAddCmd.Flags()
	flags.StringVarP(&systemCreateCookbookType, "type", "t", "", "Cookbook type (non-editable). Cookbook types can't repeat for one system")
	flags.StringArrayVarP(&systemDynamicParams, "parameters", "p", systemDynamicParams, "Add custom parameters for cookbook. SINGLE PAR: [ -p waf=true ], MULTIPLE PAR: [ -p waf=true -p timeout=200 ], MULTIPLE VALUES: [ -p versions=''7, 5.4'']")

	systemCookbookAddCmd.MarkFlagRequired("type")
	// #endregion

	//-------------------------------------  SYSTEMS/COOKBOOKS PARAMETERS (get) --------------------------------------

	// #region SYSTEMS/COOKBOOKS PARAMETERS (get)

	// ---- GET COOKBOOKTYPES PARAMETERS
	systemCookbookCmd.AddCommand(SystemCookbookTypesGetCmd)

	//flags needed to get specific parameters info
	SystemCookbookTypesGetCmd.Flags().StringVarP(&systemCreateCookbookType, "type", "t", "", "Cookbook type (non-editable). Cookbook types can't repeat for one system")
	SystemCookbookTypesGetCmd.MarkFlagRequired("type")
	// #endregion

	//-------------------------------------  SYSTEMS/SSH KEYS (get/ add / delete) --------------------------------------

	// #region SYSTEMS/SSH KEYS (get/ add / delete)

	// SSH KEYS
	systemCmd.AddCommand(systemSshKeysCmd)

	systemSshKeysCmd.AddCommand(systemSshKeysGetCmd)
	addCommonGetFlags(systemSshKeysGetCmd)

	systemSshKeysCmd.AddCommand(systemSshKeysAddCmd)
	systemSshKeysCmd.AddCommand(systemSshKeysRemoveCmd)
	// #endregion
}

// Resolve an integer or name domain.
// If the domain is a name, a request is made to resolve the integer ID.
func resolveSystem(arg string) int {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id
	}

	system := Level27Client.LookupSystem(arg)
	if system == nil {
		cobra.CheckErr(fmt.Sprintf("Unable to find system: %s", arg))
		return 0
	}
	return system.Id
}
func resolveSystemProviderConfiguration(region int, arg string) int {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id
	}

	cfgs := Level27Client.GetSystemProviderConfigurations()
	for _, cfg := range cfgs {
		if cfg.Name == arg {
			return cfg.ID
		}
	}

	cobra.CheckErr(fmt.Sprintf("Unable to find provider configuration: %s", arg))
	return 0
}

//------------------------------------------------- SYSTEM TOPLEVEL (GET / DESCRIBE CREATE) ----------------------------------
// #region SYSTEM TOPLEVEL (GET / DESCRIBE / CREATE)
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
		systemID := resolveSystem(args[0])

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
var systemCreateImage, systemCreateOrganisation, systemCreateProviderConfig, systemCreateZone string

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

		zoneID, regionID := resolveZoneRegion(systemCreateZone)
		imageID := resolveRegionImage(regionID, systemCreateImage)
		orgID := resolveOrganisation(systemCreateOrganisation)
		providerConfigID := resolveSystemProviderConfiguration(regionID, systemCreateProviderConfig)

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
			SystemImage:                 imageID,
			Organisation:                orgID,
			SystemProviderConfiguration: providerConfigID,
			Zone:                        zoneID,
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

// #endregion

var systemUpdateSettings = map[string]interface{}{}
var systemUpdateSettingsFile string

var systemUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update settings on a system",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		settings := loadMergeSettings(systemUpdateSettingsFile, systemUpdateSettings)

		systemID := resolveSystem(args[0])

		system := Level27Client.SystemGetSingle(systemID)

		systemPut := types.SystemPut{
			Id:                          system.Id,
			Name:                        system.Name,
			Type:                        system.Type,
			Cpu:                         system.Cpu,
			Memory:                      system.Memory,
			Disk:                        system.Disk,
			ManagementType:              system.ManagementType,
			Organisation:                system.Organisation.ID,
			SystemImage:                 system.SystemImage.Id,
			OperatingsystemVersion:      system.OperatingSystemVersion.Id,
			SystemProviderConfiguration: system.SystemProviderConfiguration.ID,
			Zone:                        system.Zone.Id,
			PublicNetworking:            system.PublicNetworking,
			Preferredparentsystem:       system.Preferredparentsystem,
			Remarks:                     system.Remarks,
			InstallSecurityUpdates:      system.InstallSecurityUpdates,
			LimitRiops:                  system.LimitRiops,
			LimitWiops:                  system.LimitWiops,
		}

		data := roundTripJson(systemPut).(map[string]interface{})
		data = mergeMaps(data, settings)

		data["organisation"] = resolveOrganisation(fmt.Sprint(data["organisation"]))
		data["organisation"] = resolveOrganisation(fmt.Sprint(data["organisation"]))
		data["organisation"] = resolveOrganisation(fmt.Sprint(data["organisation"]))

		Level27Client.SystemUpdate(systemID, data)
	},
}

var systemDeleteForce bool
var systemDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a system",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		systemID := resolveSystem(args[0])

		if systemDeleteForce {
			Level27Client.SystemDeleteForce(systemID)
		} else {
			Level27Client.SystemDelete(systemID)
		}
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
	Short: "Get a list of all checks from a system",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// check for valid system ID
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalln("Not a valid system ID!")
		}

		// Creating readable output
		outputFormatTableFuncs(getSystemChecks(id), []string{"ID", "CHECKTYPE", "STATUS", "LAST_STATUS_CHANGE", "INFORMATION"},
			[]interface{}{"Id", "CheckType", "Status", func(s types.SystemCheck) string { return utils.FormatUnixTime(s.DtLastStatusChanged) }, "StatusInformation"})

	},
}

func getSystemChecks(id int) []types.SystemCheck {

	return Level27Client.SystemCheckGetList(id, optGetParameters)

}

// ---------------- CREATE CHECK
var systemCheckCreateType string
var systemCheckCreateCmd = &cobra.Command{
	Use:   "add [system ID] [parameters]",
	Short: "add a new check to a specific system",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// check for valid system ID
		id := checkSingleIntID(args, "check")

		var err error
		// get the value of the flag type set by user
		checkTypeInput := cmd.Flag("type").Value.String()

		//get all data from the chosen checktype returned as Systemchecktype struct
		checktypeResult := Level27Client.SystemCheckTypeGet(checkTypeInput)
		possibleParameters := checktypeResult.ServiceType.Parameters

		// create base of json container, will be used for post request and eventually filled with custom parameters
		jsonObjCheckPost := gabs.New()
		jsonObjCheckPost.Set(checkTypeInput, "checktype")

		// if user wants to use custom parameters
		if cmd.Flag("parameters").Changed {
			// check if given parameters and usage of -p flag is correct
			customParameterDict := SplitCustomParameters(systemDynamicParams)

			// loop over all given custom parameters by user
			for customParameterName, customParameterValue := range customParameterDict {
				var isCustomParameterValid bool = false
				// loop over all possible parameters we got back form the API
				for i := range possibleParameters {
					possibleParName := possibleParameters[i].Name

					//when match found between custom paramater and possible parameter
					if possibleParName == customParameterName {
						isCustomParameterValid = true
						jsonObjCheckPost.Set(customParameterValue, customParameterName)
					}
				}
				if !isCustomParameterValid {
					err = fmt.Errorf("given parameter name is not valid: '%v'", customParameterName)
					log.Fatal(err)
				}
			}

		}
		//log.Print(jsonObjCookbookPost.StringIndent("", " "))
		Level27Client.SystemCheckCreate(id, jsonObjCheckPost)
	},
}

// ------------- GET CHECK PARAMETERS (for specific checktype)
var systemChecktypeParametersGetCmd = &cobra.Command{
	Use:   "parameters",
	Short: "Show all default parameters for a specific checktype.",
	Run: func(cmd *cobra.Command, args []string) {
		chosenType := cmd.Flag("type").Value.String()

		checktypeResult := Level27Client.SystemCheckTypeGet(chosenType)

		outputFormatTable(checktypeResult.ServiceType.Parameters, []string{"NAME", "DESCRIPTION", "DEFAULT_VALUE"}, []string{"Name", "Description", "DefaultValue"})

	},
}

//------------------------------------------------- SYSTEM/CHECKS ACTIONS (GET / DELETE / UPDATE) ----------------------------------
// -------------- GET DETAILS FROM A CHECK
var systemCheckGetSingleCmd = &cobra.Command{
	Use:   "describe [systemID] [checkID]",
	Short: "Get detailed info about a specific check.",
	Args:  cobra.ExactArgs(2),
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
	Use:   "delete [systemID] [checkID]",
	Short: "Delete a specific check from a system",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		//check for valid system ID
		systemID := checkSingleIntID(args, "system")

		//check for valid system checkID
		checkID, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalln("Not a valid check ID!")
		}

		Level27Client.SystemCheckDelete(systemID, checkID, systemCheckDeleteConfirmed)
	},
}

// -------------- UPDATE SPECIFIC CHECK
var systemCheckUpdateCmd = &cobra.Command{
	Use:   "update [SystemID] [CheckID]",
	Short: "update a specific check from a system",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// //check for valid system ID
		// systemID, err := strconv.Atoi(args[0])
		// if err != nil {
		// 	log.Fatalln("Not a valid system ID!")
		// }

		// //check for valid system checkID
		// checkID, err := strconv.Atoi(args[1])
		// if err != nil {
		// 	log.Fatalln("Not a valid check ID!")
		// }
		// // get the current data from the check
		// currentData := Level27Client.SystemCheckDescribe(systemID, checkID)

		// request := types.SystemCheckRequestHttp{
		// 	Checktype: currentData.CheckType,
		// 	Port: ,
		// }
		// if cmd.Flag("port").Changed {
		// 	currentData.CheckParameters.p
		// }
		// Level27Client.SystemCheckUpdate(systemID, checkID, nil)
	},
}

// ---------------------------------------------- ACTIONS ON SPECIFIC SYSTEM ----------------------------------------------

// #region SYSTEM ACTIONS

var systemActionsCmd = &cobra.Command{
	Use:   "actions",
	Short: "Actions for systems such as rebooting",
}

var systemActionsStartCmd = &cobra.Command{
	Use:  "start",
	Args: cobra.ExactArgs(1),
	Run:  func(cmd *cobra.Command, args []string) { runAction("start", args) },
}

var systemActionsStopCmd = &cobra.Command{
	Use:  "stop",
	Args: cobra.ExactArgs(1),
	Run:  func(cmd *cobra.Command, args []string) { runAction("stop", args) },
}

var systemActionsShutdownCmd = &cobra.Command{
	Use:  "shutdown",
	Args: cobra.ExactArgs(1),
	Run:  func(cmd *cobra.Command, args []string) { runAction("shutdown", args) },
}

var systemActionsRebootCmd = &cobra.Command{
	Use:  "reboot",
	Args: cobra.ExactArgs(1),
	Run:  func(cmd *cobra.Command, args []string) { runAction("reboot", args) },
}

var systemActionsResetCmd = &cobra.Command{
	Use:  "reset",
	Args: cobra.ExactArgs(1),
	Run:  func(cmd *cobra.Command, args []string) { runAction("reset", args) },
}

var systemActionsEmergencyPowerOffCmd = &cobra.Command{
	Use:  "emergencyPowerOff",
	Args: cobra.ExactArgs(1),
	Run:  func(cmd *cobra.Command, args []string) { runAction("emergencyPowerOff", args) },
}

var systemActionsDeactivateCmd = &cobra.Command{
	Use:  "deactivate",
	Args: cobra.ExactArgs(1),
	Run:  func(cmd *cobra.Command, args []string) { runAction("deactivate", args) },
}

var systemActionsActivateCmd = &cobra.Command{
	Use:  "activate",
	Args: cobra.ExactArgs(1),
	Run:  func(cmd *cobra.Command, args []string) { runAction("activate", args) },
}

var systemActionsAutoInstallCmd = &cobra.Command{
	Use:  "autoInstall",
	Args: cobra.ExactArgs(1),
	Run:  func(cmd *cobra.Command, args []string) { runAction("autoInstall", args) },
}

// #endregion

func runAction(action string, args []string) {
	id := resolveSystem(args[0])

	Level27Client.SystemAction(id, action)
}

//------------------------------------------------- SYSTEM/COOKBOOKS PARAMETERS GET ----------------------------------

// ----------- GET COOKBOOKTYPE PARAMETERS
// seperate command used to see wich parameters can be used for a specific cookbooktype. also shows the description and default values
var SystemCookbookTypesGetCmd = &cobra.Command{
	Use:   "parameters",
	Short: "Show all default parameters for a specific cookbooktype.",
	Run: func(cmd *cobra.Command, args []string) {

		// get the user input from the type flag
		inputType := cmd.Flag("type").Value.String()

		// Get request to get all cookbooktypes data
		validCookbooktype, _ := Level27Client.SystemCookbookTypeGet(inputType)

		outputFormatTable(validCookbooktype.CookbookType.Parameters, []string{"NAME", "DESCRIPTION", "DEFAULT_VALUE"}, []string{"Name", "Description", "DefaultValue"})

	},
}

//------------------------------------------------- SYSTEM/COOKBOOKS TOPLEVEL (GET / CREATE) ----------------------------------
// ---------------- MAIN COMMAND (checks)
var systemCookbookCmd = &cobra.Command{
	Use:   "cookbooks",
	Short: "Manage systems cookbooks",
}

// ---------- GET COOKBOOKS
var systemCookbookGetCmd = &cobra.Command{
	Use:   "get [system ID]",
	Short: "Gets a list of all cookbooks from a system.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// check for valid system ID
		id := checkSingleIntID(args, "system")

		outputFormatTable(getSystemCookbooks(id), []string{"ID", "COOKBOOKTYPE", "STATUS"}, []string{"Id", "CookbookType", "Status"})
	},
}

func getSystemCookbooks(id int) []types.Cookbook {

	return Level27Client.SystemCookbookGetList(id)
}

func CheckforValidType(input string, validTypes []string) (string, bool) {
	var isTypeValid bool
	// check if given cookbooktype is 1 of valid options
	for _, cookbooktype := range validTypes {
		if strings.ToLower(input) == cookbooktype {
			input = cookbooktype
			isTypeValid = true
			return input, isTypeValid
		}
	}
	return "", isTypeValid
}

// ----------- ADD COOKBOOK TO SPECIFIC SYSTEM
var systemDynamicParams []string
var systemCreateCookbookType string
var systemCookbookAddCmd = &cobra.Command{
	Use:   "add [systemID] [flags]",
	Short: "add a cookbook to a system",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//checking for valid system ID
		id := checkSingleIntID(args, "system")

		var err error
		// get information about the current chosen system [systemID]
		currentSystem := Level27Client.SystemGetSingle(id)
		currentSystemOS := fmt.Sprintf("%v %v", currentSystem.OperatingSystemVersion.OsName, currentSystem.OperatingSystemVersion.OsVersion)

		// get the user input from the type flag (cookbooktype)
		inputType := cmd.Flag("type").Value.String()

		// get all data from the chossen cookbooktype and its parameteroptions if there are any
		validCookbooktype, parameterOptions := Level27Client.SystemCookbookTypeGet(inputType)

		// // create base of json container, will be used for post request and eventually filled with custom parameters
		jsonObjCookbookPost := gabs.New()
		jsonObjCookbookPost.Set(inputType, "cookbooktype")

		// // when user wants to use custom parameters
		if cmd.Flag("parameters").Changed {

			// split the slice of customparameters set by user into key/value pairs. also check if declaration method is used correctly (-p key=value).
			customParameterDict := SplitCustomParameters(systemDynamicParams)

			// loop over the filtered parameters set by the user
			for key, value := range customParameterDict {

				var isValidParameter bool = false

				//loop over all possible parameters for the chosen type
				for _, parameter := range validCookbooktype.CookbookType.Parameters {
					if parameter.Name == key {
						isValidParameter = true
						// when parameter type is select -> value can only be one of the selectable options + value has specific rules
						if parameter.Type == "select" {

							// var isValidValue bool = false
							// check in json obj with all selectable options for the cookbooktype if parameteroption exists
							if parameterOptions.Exists(key) {

								// check wich type the value currently is. value needs to be of type array for selectable parameter + see if value is 1 or multiple values for 1 key
								valueType := reflect.TypeOf(value)

								// when value = array or slice -> key contains multiple values (PHP versions -> 7.2 + 7.3)
								if valueType.Kind() == reflect.Array || valueType.Kind() == reflect.Slice {
									rawValues := reflect.ValueOf(value)

									// need to convert interface to a go slice
									values := make([]interface{}, rawValues.Len())
									for i := 0; i < rawValues.Len(); i++ {
										values[i] = rawValues.Index(i).Interface()
									}

									// loop over each value
									for i := range values {

										// check if value can be installed on the os and if the value needs to be exclusive
										_, isExclusive := isValueValidForParameter(*parameterOptions.Search(key), values[i], currentSystemOS)

										// when value is exclusive -> cannot be installed with other values of same sort
										if isExclusive {
											message := fmt.Sprintf("Given Value: '%v' NOT possible for multiselect.", values[i])
											err := errors.New(message)
											log.Fatal(err)
										}

									}
									jsonObjCookbookPost.SetP(value, key)

								} else {
									// --- SINGLE VALUE in this case we have a single value asigned to a key
									valueAsString := fmt.Sprintf("%v", value)
									// if given value is one of the possible options for the given parameter key
									if parameterOptions.Search(key).Exists(valueAsString) {
										isAvailable, _ := isValueValidForParameter(*parameterOptions.Search(key), value, currentSystemOS)
										if isAvailable {
											//key has one value but needs to be sent in array type
											var values []interface{}
											values = append(values, value)
											jsonObjCookbookPost.SetP(values, key)
										}

									} else {

										message := fmt.Sprintf("Given Value: '%v' NOT an option for given key: '%v'.", value, key)
										err = errors.New(message)
									}

								}
							} else {
								message := fmt.Sprintf("No parameter options found for: '%v'.", value)
								err = errors.New(message)
							}
						} else {

							jsonObjCookbookPost.SetP(value, key)
						}
					}
				}

				if !isValidParameter {
					message := fmt.Sprintf("Parameter: '%v' NOT known for cookbooktype %v.", key, inputType)
					err = errors.New(message)
				}

			}

			// 	// when parameters or values are not valid -> error, close command
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("custom")
			log.Print(jsonObjCookbookPost.StringIndent("", " "))
			//Level27Client.SystemCookbookAdd(id, jsonObjCookbookPost)
		} else {
			log.Println("standaard")
			log.Print(jsonObjCookbookPost.StringIndent("", " "))
			//Level27Client.SystemCookbookAdd(id, jsonObjCookbookPost)
		}

	},
}

// // ------------------------------------------------------ SSH KEYS

var systemSshKeysCmd = &cobra.Command{
	Use: "sshkeys",
}

var systemSshKeysGetCmd = &cobra.Command{
	Use: "get",

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := resolveSystem(args[0])

		outputFormatTable(Level27Client.SystemGetSshKeys(id, optGetParameters), []string{"ID", "DESCRIPTION", "STATUS", "FINGERPRINT"}, []string{"ID", "Description", "ShsStatus", "Fingerprint"})
	},
}

var systemSshKeysAddCmd = &cobra.Command{
	Use: "add",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		systemID := resolveSystem(args[0])

		keyName := args[1]
		keyID, err := strconv.Atoi(keyName)
		if err != nil {
			user := viper.GetInt("user_id")
			org := viper.GetInt("org_id")
			system := Level27Client.LookupSystemNonAddedSshkey(systemID, org, user, keyName)
			if system == nil {
				existing := Level27Client.LookupSystemSshkey(systemID, keyName)
				if existing != nil {
					fmt.Println("SSH key already exists on system!")
					return
				} else {
					cobra.CheckErr("Unable to find SSH key to add")
					return
				}
			}
			keyID = system.Id
		}

		Level27Client.SystemAddSshKey(systemID, keyID)
	},
}

var systemSshKeysRemoveCmd = &cobra.Command{
	Use: "remove",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		systemID := resolveSystem(args[0])

		keyName := args[1]
		keyID, err := strconv.Atoi(keyName)
		if err != nil {
			existing := Level27Client.LookupSystemSshkey(systemID, keyName)
			if existing == nil {
				cobra.CheckErr("Unable to find SSH key to remove!")
				return
			}

			keyID = existing.ID
		}

		Level27Client.SystemRemoveSshKey(systemID, keyID)
	},
}
