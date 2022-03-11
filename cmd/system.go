package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"bitbucket.org/level27/lvl/types"
	"bitbucket.org/level27/lvl/utils"
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

	// SSH KEYS
	systemCmd.AddCommand(systemSshKeysCmd)

	systemSshKeysCmd.AddCommand(systemSshKeysGetCmd)
	addCommonGetFlags(systemSshKeysGetCmd)

	systemSshKeysCmd.AddCommand(systemSshKeysAddCmd)
	systemSshKeysCmd.AddCommand(systemSshKeysRemoveCmd)

	// NETWORKS

	systemCmd.AddCommand(systemNetworksCmd)

	systemNetworksCmd.AddCommand(systemNetworksGetCmd)

	systemNetworksCmd.AddCommand(systemNetworksDescribeCmd)

	systemNetworksCmd.AddCommand(systemNetworksAddCmd)

	systemNetworksCmd.AddCommand(systemNetworksRemoveCmd)

	// NETWORK IPS

	systemNetworksCmd.AddCommand(systemNetworkIpCmd)

	systemNetworkIpCmd.AddCommand(systemNetworkIpGetCmd)

	systemNetworkIpCmd.AddCommand(systemNetworkIpAddCmd)
	systemNetworkIpAddCmd.Flags().StringVar(&systemNetworkIpAddHostname, "hostname", "", "Hostname for the IP address. If not specified the system hostname is used.")
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

func resolveSystemHasNetwork(systemID int, arg string) int {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id
	}

	network := Level27Client.LookupSystemHasNetworks(systemID, arg)
	if network == nil {
		cobra.CheckErr(fmt.Sprintf("Unable to find network: %s", arg))
		return 0
	}

	return network.ID
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

var systemUpdateSettings = map[string]interface{}{}
var systemUpdateSettingsFile string

var systemUpdateCmd = &cobra.Command{
	Use: "update",
	Short: "Update settings on a system",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		settings := loadMergeSettings(systemUpdateSettingsFile, systemUpdateSettings)

		systemID := resolveSystem(args[0])

		system := Level27Client.SystemGetSingle(systemID)

		systemPut := types.SystemPut{
			Id: system.Id,
			Name: system.Name,
			Type: system.Type,
			Cpu: system.Cpu,
			Memory: system.Memory,
			Disk: system.Disk,
			ManagementType: system.ManagementType,
			Organisation: system.Organisation.ID,
			SystemImage: system.SystemImage.Id,
			OperatingsystemVersion: system.OperatingSystemVersion.Id,
			SystemProviderConfiguration: system.SystemProviderConfiguration.ID,
			Zone: system.Zone.Id,
			PublicNetworking: system.PublicNetworking,
			Preferredparentsystem: system.Preferredparentsystem,
			Remarks: system.Remarks,
			InstallSecurityUpdates: system.InstallSecurityUpdates,
			LimitRiops: system.LimitRiops,
			LimitWiops: system.LimitWiops,
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
	Use: "delete",
	Short: "Delete a system",
	Args: cobra.ExactArgs(1),
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
	Short: "Get a list of all checks for a system",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// check for valid system ID
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalln("Not a valid system ID!")
		}

		// Creating readable output
		outputFormatTableFuncs(getSystemChecks(id), []string{"ID", "CHECKTYPE", "STATUS", "LAST_STATUS_CHANGE", "INFORMATION"},
			[]interface{}{"Id", "CheckType", "Status",func(s types.SystemCheck) string {return utils.FormatUnixTime(s.DtLastStatusChanged)}, "StatusInformation"})



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

func runAction(action string, args []string) {
	id := resolveSystem(args[0])

	Level27Client.SystemAction(id, action)
}

// SSH KEYS

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

// NETWORKS

var systemNetworksCmd = &cobra.Command{
	Use: "networks",
}

var systemNetworksGetCmd = &cobra.Command{
	Use: "get [system]",
	Short: "Get list of networks on a system",

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		systemID := resolveSystem(args[0])
		system := Level27Client.SystemGetSingle(systemID)

		outputFormatTableFuncs(system.Networks, []string{"ID", "Network ID", "Type", "Name", "MAC", "IPs"}, []interface{}{"ID", "NetworkID", func(net types.SystemNetwork) string {
			if net.NetPublic { return "public" }
			if net.NetCustomer { return "customer" }
			if net.NetInternal { return "internal" }
			return ""
		}, "Name", "Mac", func(net types.SystemNetwork) string {
			return strconv.Itoa(len(net.Ips))
		}})
	},
}

var systemNetworksDescribeCmd = &cobra.Command{
	Use: "describe [system]",
	Short: "Display detailed information about all networks on a system",

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		systemID := resolveSystem(args[0])
		system := Level27Client.SystemGetSingle(systemID)
		networks := Level27Client.SystemGetHasNetworks(systemID)

		outputFormatTemplate(types.DescribeSystemNetworks{
			Networks: system.Networks,
			HasNetworks: networks,
		}, "templates/systemNetworks.tmpl")
	},
}

var systemNetworksAddCmd = &cobra.Command{
	Use: "add [system] [network]",
	Short: "Add a network to a system",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		systemID := resolveSystem(args[0])
		networkID := resolveNetwork(args[1])

		Level27Client.SystemAddHasNetwork(systemID, networkID)
	},
}

var systemNetworksRemoveCmd = &cobra.Command{
	Use: "remove [system] [network]",
	Short: "Remove a network from a system",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		systemID := resolveSystem(args[0])
		networkID := resolveSystemHasNetwork(systemID, args[1])

		Level27Client.SystemRemoveHasNetwork(systemID, networkID)
	},
}

var systemNetworkIpCmd = &cobra.Command{
	Use: "ip",
	Short: "Manage IP addresses on network connections",
}

var systemNetworkIpGetCmd = &cobra.Command{
	Use: "get [system] [network]",
	Short: "Get all IP addresses for a system network",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		systemID := resolveSystem(args[0])
		networkID := resolveSystemHasNetwork(systemID, args[1])

		ips := Level27Client.SystemGetHasNetworkIps(systemID, networkID)
		outputFormatTableFuncs(ips, []string{"ID", "Public IP", "IP", "Hostname", "Status"}, []interface{}{"ID", func(i types.SystemHasNetworkIp) string {
				if i.Ipv4 != "" {
					i, _ := strconv.Atoi(i.PublicIpv4)
					if i == 0 {
						return ""
					} else {
						return utils.Ipv4IntToString(i)
					}
				} else {
					return fmt.Sprint(i.Ipv6)
				}
			},
			func(i types.SystemHasNetworkIp) string {
				if i.Ipv4 != "" {
					i, _ := strconv.Atoi(i.Ipv4)
					if i == 0 {
						return ""
					} else {
						return utils.Ipv4IntToString(i)
					}
				} else {
					return fmt.Sprint(i.Ipv6)
				}
		}, "Hostname", "Status"})
	},
}

var systemNetworkIpAddHostname string

var systemNetworkIpAddCmd = &cobra.Command{
	Use: "add [system] [network] [address]",
	Short: "Get all IP addresses for a system network",
	Long: "Adds an IP address to a system network. Address can be either IPv4 or IPv6. The special values 'auto' and 'auto-v6' automatically fetch an unused address to use.",

	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		systemID := resolveSystem(args[0])
		system := Level27Client.SystemGetSingle(systemID)
		hasNetworkID := resolveSystemHasNetwork(systemID, args[1])
		network := Level27Client.GetSystemHasNetwork(systemID, hasNetworkID)
		networkID := network.Network.ID
		address := args[2]

		if address == "auto" || address == "auto-v6" {
			located := Level27Client.NetworkLocate(networkID)

			if address == "auto" {
				address = located.Ipv4[0]
			} else {
				address = located.Ipv6[0]
			}

			if address == "" {
				cobra.CheckErr("Unable to find a free IP address")
			}
		}

		var data types.SystemHasNetworkIpAdd
		public := network.Network.Public

		if strings.Contains(address, ":") {
			// IPv6
			if public {
				data.PublicIpv6 = address
			} else {
				data.Ipv6 = address
			}
		} else {
			// IPv4
			if public {
				data.PublicIpv4 = address
			} else {
				data.Ipv4 = address
			}
		}

		data.Hostname = system.Hostname
		if systemNetworkIpAddHostname != "" {
			data.Hostname = systemNetworkIpAddHostname
		}

		Level27Client.SystemAddHasNetworkIps(systemID, hasNetworkID, data)
	},
}
