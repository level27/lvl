package cmd

import (
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/level27/l27-go"
	"github.com/level27/lvl/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
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
	flags.StringVarP(&systemCreateManageType, "management", "", "basic", "Managament type (one of basic, professional, enterprise, professional_level27).")
	flags.BoolVarP(&systemCreatePublicNetworking, "publicNetworking", "", true, "For digitalOcean servers always true. (non-editable)")
	flags.StringVarP(&systemCreateImage, "image", "", "", "The ID of a systemimage. (must match selected configuration and zone. non-editable)")
	flags.StringVarP(&systemCreateOrganisation, "organisation", "", "", "The unique ID of an organisation")
	flags.StringVarP(&systemCreateProviderConfig, "config", "", "", "The unique ID of a SystemproviderConfiguration")
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

	// ------------------------------------ MONITORING ON SPECIFIC SYSTEM ----------------------------------------------
	// ---- MONITORING COMMAND
	systemCmd.AddCommand(systemMonitoringCmd)

	// ---- MONITORING ON
	systemMonitoringCmd.AddCommand(systemMonitoringOnCmd)
	// ---- MONITORING OFF
	systemMonitoringCmd.AddCommand(systemMonitoringOffCmd)
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
	settingsFileFlag(systemUpdateCmd)
	settingString(systemUpdateCmd, updateSettings, "name", "New name for this system")
	settingInt(systemUpdateCmd, updateSettings, "cpu", "Set amount of CPU cores of the system")
	settingInt(systemUpdateCmd, updateSettings, "memory", "Set amount of memory in GB of the system")
	settingString(systemUpdateCmd, updateSettings, "managementType", "Set management type of the system")
	settingString(systemUpdateCmd, updateSettings, "organisation", "Set organisation that owns this system. Can be both a name or an ID")
	settingInt(systemUpdateCmd, updateSettings, "publicNetworking", "")
	settingInt(systemUpdateCmd, updateSettings, "limitRiops", "Set read IOPS limit")
	settingInt(systemUpdateCmd, updateSettings, "limitWiops", "Set write IOPS limit")
	settingInt(systemUpdateCmd, updateSettings, "installSecurityUpdates", "Set security updates mode index")
	settingString(systemUpdateCmd, updateSettings, "remarks", "")

	// --- Delete

	systemCmd.AddCommand(systemDeleteCmd)
	systemDeleteCmd.Flags().BoolVar(&systemDeleteForce, "force", false, "")
	addDeleteConfirmFlag(systemDeleteCmd)
	// #endregion

	//-------------------------------------  SYSTEMS/CHECKS PARAMETERS (get parameters) --------------------------------------
	// #region GET CHECK PARAMETERS
	// ---- GET PARAMETERS (for specific checktype)
	systemCheckCmd.AddCommand(systemChecktypeParametersGetCmd)

	// flags needed to get checktype parameters
	systemChecktypeParametersGetCmd.Flags().StringVarP(&systemCheckCreateType, "type", "t", "", "Check type to see all its available parameters")
	// #endregion

	//-------------------------------------  SYSTEMS/CHECKS TOPLEVEL (get / post) --------------------------------------
	// #region SYSTEMS/CHECKS TOPLEVEL
	systemCmd.AddCommand(systemCheckCmd)
	// ---- GET LIST OF ALL CHECKS
	systemCheckCmd.AddCommand(systemCheckGetCmd)
	addCommonGetFlags(systemCheckGetCmd)

	// ---- CREATE NEW CHECK
	systemCheckCmd.AddCommand(systemCheckAddCmd)

	// -- flags needed to create a check
	flags = systemCheckAddCmd.Flags()
	flags.StringVarP(&systemCheckCreateType, "type", "t", "", "Check type (non-editable)")
	systemCheckAddCmd.MarkFlagRequired("type")

	// -- optional flag
	flags.StringArrayVarP(&systemDynamicParams, "parameters", "p", systemDynamicParams, "Add custom parameters for a check. usage -> SINGLE PAR: [ -p waf=true ], MULTIPLE PAR: [ -p waf=true -p timeout=200 ], MULTIPLE VALUES: [ -p versions=''7, 5.4'']")

	// #endregion

	//-------------------------------------  SYSTEMS/CHECKS ACTIONS (get/ delete/ update) --------------------------------------
	// #region SYSTEMS/CHECKS ACTIONS
	// --- DESCRIBE CHECK
	systemCheckCmd.AddCommand(systemCheckGetSingleCmd)
	// --- DELETE CHECK
	systemCheckCmd.AddCommand(systemCheckDeleteCmd)
	addDeleteConfirmFlag(systemCheckDeleteCmd)

	// --- UPDATE CHECK (ONLY FOR HTTP REQUEST)
	systemCheckCmd.AddCommand(systemCheckUpdateCmd)

	// flag needed to update a specific check
	systemCheckUpdateCmd.Flags().StringArrayVarP(&systemDynamicParams, "parameters", "p", systemDynamicParams, "Add custom parameters for a check. Usage -> SINGLE PAR: [ -p waf=true ], MULTIPLE PAR: [ -p waf=true -p timeout=200 ], MULTIPLE VALUES: [ -p versions=''7, 5.4'']")
	systemCheckUpdateCmd.MarkFlagRequired("parameters")

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

	//-------------------------------------  SYSTEMS/COOKBOOKS SPECIFIC (describe / delete / update) --------------------------------------
	// #region SYSTEM/COOKBOOKS SPECIFIC (DESCRIBE / DELETE / UPDATE)

	// --- DESCRIBE
	systemCookbookCmd.AddCommand(systemCookbookDescribeCmd)

	// --- DELETE
	systemCookbookCmd.AddCommand(systemCookbookDeleteCmd)
	addDeleteConfirmFlag(systemCookbookDeleteCmd)

	// --- UPDATE
	systemCookbookCmd.AddCommand(systemCookbookUpdateCmd)
	// flags for update
	systemCookbookUpdateCmd.Flags().StringArrayVarP(&systemDynamicParams, "parameters", "p", systemDynamicParams, "Add custom parameters for a check. Usage -> SINGLE PAR: [ -p waf=true ], MULTIPLE PAR: [ -p waf=true -p timeout=200 ], MULTIPLE VALUES: [ -p versions=''7, 5.4'']")
	systemCookbookUpdateCmd.MarkFlagRequired("parameters")
	// #endregion

	//-------------------------------------  SYSTEMS/INTEGRITYCHECKS (get / post / download) --------------------------------------
	addIntegrityCheckCmds(systemCmd, "systems", resolveSystem)

	//-------------------------------------  SYSTEMS/GROUPS (get/ add / describe / delete) --------------------------------------
	// #region SYSTEMS/GROUPS (get/ add / delete / describe)

	systemCmd.AddCommand(SystemSystemgroupsCmd)

	// --- GET
	SystemSystemgroupsCmd.AddCommand(SystemSystemgroupsGetCmd)

	// --- ADD
	SystemSystemgroupsCmd.AddCommand(SystemSystemgroupsAddCmd)

	// --- DELETE
	SystemSystemgroupsCmd.AddCommand(SystemSystemgroupsRemoveCmd)

	//-------------------------------------  SYSTEMS/SSH KEYS (get/ add / delete) --------------------------------------
	// #region SYSTEMS/SSH KEYS (get/ add / describe / delete)

	// SSH KEYS
	systemCmd.AddCommand(systemSshKeysCmd)

	systemSshKeysCmd.AddCommand(systemSshKeysGetCmd)
	addCommonGetFlags(systemSshKeysGetCmd)

	systemSshKeysCmd.AddCommand(systemSshKeysAddCmd)
	systemSshKeysCmd.AddCommand(systemSshKeysRemoveCmd)

	// #endregion

	// SYSTEM SSH
	systemCmd.AddCommand(systemSshCmd)

	// SYSTEM SCP
	systemCmd.AddCommand(systemScpCommand)

	//------------------------------------- NETWORKS -------------------------------------
	// #region NETWORKS

	systemCmd.AddCommand(systemNetworkCmd)

	systemNetworkCmd.AddCommand(systemNetworkGetCmd)

	systemNetworkCmd.AddCommand(systemNetworkDescribeCmd)

	systemNetworkCmd.AddCommand(systemNetworkAddCmd)

	systemNetworkCmd.AddCommand(systemNetworkRemoveCmd)

	//------------------------------------- NETWORK IPS -------------------------------------

	systemNetworkCmd.AddCommand(systemNetworkIpCmd)

	systemNetworkIpCmd.AddCommand(systemNetworkIpGetCmd)

	systemNetworkIpCmd.AddCommand(systemNetworkIpAddCmd)
	systemNetworkIpAddCmd.Flags().StringVar(&systemNetworkIpAddHostname, "hostname", "", "Hostname for the IP address. If not specified the system hostname is used.")

	systemNetworkIpCmd.AddCommand(systemNetworkIpRemoveCmd)

	systemNetworkIpCmd.AddCommand(systemNetworkIpUpdateCmd)
	settingsFileFlag(systemNetworkIpUpdateCmd)
	settingString(systemNetworkIpUpdateCmd, updateSettings, "hostname", "New hostname for this IP")
	// #endregion

	// SYSTEM VOLUME
	systemCmd.AddCommand(systemVolumeCmd)

	// SYSTEM VOLUME GET
	systemVolumeCmd.AddCommand(systemVolumeGetCmd)
	addCommonGetFlags(systemVolumeGetCmd)

	// SYSTEM VOLUME CREATE
	systemVolumeCmd.AddCommand(systemVolumeCreateCmd)
	systemVolumeCreateCmd.Flags().StringVar(&systemVolumeCreateName, "name", "", "Name of the new volume")
	systemVolumeCreateCmd.Flags().StringVar(&systemVolumeCreateOrganisation, "organisation", "", "Organisation for the new volume")
	systemVolumeCreateCmd.Flags().StringVar(&systemVolumeCreateDeviceName, "deviceName", "", "Device name for the new volume")
	systemVolumeCreateCmd.Flags().BoolVar(&systemVolumeCreateAutoResize, "autoResize", false, "Enable automatic resizing")
	systemVolumeCreateCmd.Flags().IntVar(&systemVolumeCreateSpace, "space", 0, "Space of the new volume (in GB)")

	// SYSTEM VOLUME LINK
	systemVolumeCmd.AddCommand(systemVolumeLinkCmd)

	// SYSTEM VOLUME UNLINK
	systemVolumeCmd.AddCommand(systemVolumeUnlinkCmd)

	// SYSTEM VOLUME DELETE
	systemVolumeCmd.AddCommand(systemVolumeDeleteCmd)
	systemVolumeDeleteCmd.Flags().BoolVar(&systemVolumeDeleteForce, "force", false, "Do not ask for confirmation to delete the volume")

	// SYSTEM VOLUME UPDATE
	systemVolumeCmd.AddCommand(systemVolumeUpdateCmd)
	settingsFileFlag(systemVolumeUpdateCmd)
	settingString(systemVolumeUpdateCmd, updateSettings, "name", "New name for the volume")
	settingBool(systemVolumeUpdateCmd, updateSettings, "autoResize", "New autoResize setting")
	settingInt(systemVolumeUpdateCmd, updateSettings, "space", "New volume space (in GB)")

	// ACCESS
	addAccessCmds(systemCmd, "systems", resolveSystem)

	// BILLING
	addBillingCmds(systemCmd, "systems", resolveSystem)

	// JOBS
	addJobCmds(systemCmd, "system", resolveSystem)
}

// Resolve an integer or name domain.
// If the domain is a name, a request is made to resolve the integer ID.
func resolveSystem(arg string) (int, error) {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.LookupSystem(arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"system",
		func(app l27.System) string { return fmt.Sprintf("%s (%d)", app.Name, app.Id) })

	if err != nil {
		return 0, err
	}

	return res.Id, err
}

func resolveSystemCookbook(systemID int, arg string) (int, error) {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id, nil
	}

	cookbook, err := Level27Client.SystemCookbookLookup(systemID, arg)
	if err != nil {
		return 0, err
	}

	if cookbook == nil {
		return 0, fmt.Errorf("system (%d) does not have a cookbook of type '%s'", systemID, arg)
	}

	return cookbook.Id, nil
}

func resolveSystemProviderConfiguration(region int, arg string) (int, error) {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id, nil
	}

	cfgs, err := Level27Client.GetSystemProviderConfigurations()
	if err != nil {
		return 0, err
	}

	for _, cfg := range cfgs {
		if cfg.Name == arg {
			return cfg.ID, nil
		}
	}

	return 0, fmt.Errorf("unable to find provider configuration: '%s'", arg)
}

func resolveSystemHasNetwork(systemID int, arg string) (int, error) {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.LookupSystemHasNetworks(systemID, arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"system network",
		func(app l27.SystemHasNetwork) string { return fmt.Sprintf("%s (%d)", app.Network.Name, app.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, err
}

func resolveSystemHasNetworkIP(systemID int, hasNetworkID int, arg string) (int, error) {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.LookupSystemHasNetworkIp(systemID, hasNetworkID, arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"system network IP",
		func(app l27.SystemHasNetworkIp) string { return fmt.Sprintf("%d", app.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, err
}

func resolveSystemVolume(systemID int, arg string) (int, error) {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id, nil
	}

	ip, err := Level27Client.LookupSystemVolumes(systemID, arg)
	if err != nil {
		return 0, err
	}

	if ip == nil {
		return 0, fmt.Errorf("nable to find volume: %s", arg)
	}

	return ip.ID, nil
}

//------------------------------------------------- SYSTEM TOPLEVEL (GET / DESCRIBE CREATE) ----------------------------------
// #region SYSTEM TOPLEVEL (GET / DESCRIBE / CREATE)
//----------------------------------------- GET ---------------------------------------
var systemGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get a list of all curent systems",
	RunE: func(cmd *cobra.Command, args []string) error {
		ids, err := convertStringsToIds(args)
		if err != nil {
			return err
		}

		systems, err := getSystems(ids)
		if err != nil {
			return err
		}

		outputFormatTable(systems, []string{"ID", "NAME", "STATUS"}, []string{"Id", "Name", "Status"})
		return nil
	},
}

func getSystems(ids []int) ([]l27.System, error) {
	if len(ids) == 0 {
		return Level27Client.SystemGetList(optGetParameters)
	}

	systems := make([]l27.System, len(ids))
	for idx, id := range ids {
		var err error
		systems[idx], err = Level27Client.SystemGetSingle(id)
		if err != nil {
			return nil, err
		}
	}
	return systems, nil
}

//----------------------------------------- DESCRIBE ---------------------------------------
var systemDescribeHideJobs = false

var systemDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Get detailed information about a system.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		var system DescribeSystem
		system.System, err = Level27Client.SystemGetSingle(systemID)
		if err != nil {
			return err
		}

		if !systemDescribeHideJobs {
			system.Jobs, err = Level27Client.EntityJobHistoryGet("system", systemID)
			if err != nil {
				return err
			}

			for idx, j := range system.Jobs {
				system.Jobs[idx], err = Level27Client.JobHistoryRootGet(j.Id)

				if err != nil {
					return err
				}
			}
		}

		system.SshKeys, err = Level27Client.SystemGetSshKeys(systemID, l27.CommonGetParams{})
		if err != nil {
			return err
		}

		securityUpdates, err := Level27Client.SecurityUpdateDates()
		if err != nil {
			return err
		}

		system.InstallSecurityUpdatesString = securityUpdates[system.InstallSecurityUpdates]
		system.HasNetworks, err = Level27Client.SystemGetHasNetworks(systemID)
		if err != nil {
			return err
		}

		system.Volumes, err = Level27Client.SystemGetVolumes(systemID, l27.CommonGetParams{})
		if err != nil {
			return err
		}

		if system.System.MonitoringEnabled {
			system.Checks, err = Level27Client.SystemCheckGetList(systemID, l27.CommonGetParams{})
			if err != nil {
				return err
			}
		}

		outputFormatTemplate(system, "templates/system.tmpl")
		return nil
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
	Use:     "create",
	Short:   "Create a new system",
	Example: "lvl system create -n mySystemName --zone hasselt --organisation level27 --image 'Ubuntu 20.04 LTS' --config 'Level27 Small' --management professional_level27",
	RunE: func(cmd *cobra.Command, args []string) error {

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
				return fmt.Errorf("ERROR: given managementType is not valid: '%v'", managementTypeValue)
			}
		}

		zoneID, regionID, err := resolveZoneRegion(systemCreateZone)
		if err != nil {
			return err
		}

		imageID, err := resolveRegionImage(regionID, systemCreateImage)
		if err != nil {
			return err
		}

		orgID, err := resolveOrganisation(systemCreateOrganisation)
		if err != nil {
			return err
		}

		providerConfigID, err := resolveSystemProviderConfiguration(regionID, systemCreateProviderConfig)
		if err != nil {
			return err
		}

		// Using data from the flags to make the right type used for posting a new system. (types systemPost)
		RequestData := l27.SystemPost{
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

		system, err := Level27Client.SystemCreate(RequestData)
		if err != nil {
			return err
		}

		log.Printf("System created! [Fullname: '%v' , ID: '%v']", system.Name, system.Id)
		return nil
	},
}

// #endregion

//------------------------------------------------- SYSTEM SPECIFIC (UPDATE / FORCE DELETE ) ----------------------------------
// #region SYSTEM SPECIFIC (UPDATE / FORCE DELETE)
var systemUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update settings on a system",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		settings, err := loadMergeSettings(updateSettingsFile, updateSettings)
		if err != nil {
			return err
		}

		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		system, err := Level27Client.SystemGetSingle(systemID)
		if err != nil {
			return err
		}

		systemPut := l27.SystemPut{
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
			Preferredparentsystem:       system.Preferredparentsystem.ID,
			Remarks:                     system.Remarks,
			InstallSecurityUpdates:      system.InstallSecurityUpdates,
			LimitRiops:                  system.LimitRiops,
			LimitWiops:                  system.LimitWiops,
		}

		data := utils.RoundTripJson(systemPut).(map[string]interface{})
		data = mergeMaps(data, settings)

		data["organisation"], err = resolveOrganisation(fmt.Sprint(data["organisation"]))
		if err != nil {
			return err
		}

		err = Level27Client.SystemUpdate(systemID, data)
		if err != nil {
			return err
		}

		log.Print("System succesfully updated!")
		return nil
	},
}

var systemDeleteForce bool
var systemDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a system",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		if !optDeleteConfirmed {
			system, err := Level27Client.SystemGetSingle(systemID)
			if err != nil {
				return err
			}

			if !confirmPrompt(fmt.Sprintf("Delete system %s (%d)?", system.Name, system.Id)) {
				return nil
			}
		}

		if systemDeleteForce {
			err = Level27Client.SystemDeleteForce(systemID)
		} else {
			err = Level27Client.SystemDelete(systemID)
		}

		return err
	},
}

// #endregion

//------------------------------------------------- SYSTEM/CHECKS TOPLEVEL (GET / CREATE) ----------------------------------

// ---------------- MAIN COMMAND (checks)
var systemCheckCmd = &cobra.Command{
	Use:   "checks",
	Short: "Manage systems checks",
}

// #region SYSTEM/CHECKS (GET / CREATE)

// ---------------- GET
var systemCheckGetCmd = &cobra.Command{
	Use:   "get [system ID]",
	Short: "Get a list of all checks from a system",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// check for valid system ID
		id, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		system, err := Level27Client.SystemGetSingle(id)
		if err != nil {
			return err
		}

		// when monitoring is disabled on system -> checks dont need to be visible
		if system.MonitoringEnabled {
			return fmt.Errorf("monitoring is currently disabled for system: [NAME:%v - ID: %v]. Use the 'monitoring' command to change monitoring status", system.Name, system.Id)
		}

		checks, err := Level27Client.SystemCheckGetList(id, optGetParameters)
		if err != nil {
			return err
		}

		// Creating readable output
		outputFormatTableFuncs(checks, []string{"ID", "CHECKTYPE", "STATUS", "LAST_STATUS_CHANGE", "INFORMATION"},
			[]interface{}{"Id", "CheckType", "Status", func(s l27.SystemCheckGet) string { return utils.FormatUnixTime(s.DtLastStatusChanged) }, "StatusInformation"})

		return nil
	},
}

// ---------------- CREATE CHECK
var systemCheckCreateType string
var systemCheckAddCmd = &cobra.Command{
	Use:   "add [system ID] [parameters]",
	Short: "add a new check to a specific system",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// check for valid system ID
		id, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		// get the value of the flag type set by user
		checkTypeInput := cmd.Flag("type").Value.String()

		//get all data from the chosen checktype returned as Systemchecktype struct
		checktypeResult, err := Level27Client.SystemCheckTypeGet(checkTypeInput)
		if err != nil {
			return err
		}

		possibleParameters := checktypeResult.ServiceType.Parameters

		// create base of json container, will be used for post request and eventually filled with custom parameters
		jsonObjCheckPost := gabs.New()
		jsonObjCheckPost.Set(checkTypeInput, "checktype")

		// if user wants to use custom parameters
		if cmd.Flag("parameters").Changed {
			// check if given parameters and usage of -p flag is correct
			customParameterDict, err := SplitCustomParameters(systemDynamicParams)
			if err != nil {
				return err
			}

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
					return fmt.Errorf("given parameter name is not valid: '%v'", customParameterName)
				}
			}

		}

		check, err := Level27Client.SystemCheckAdd(id, jsonObjCheckPost)
		if err != nil {
			return err
		}

		log.Printf("System check added! [Checktype: '%v' , ID: '%v']", check.CheckType, check.Id)
		return nil
	},
}

// #endregion

//------------------------------------------------- SYSTEM/CHECKS PARAMETERS (GET) ----------------------------------
// #region SYSTEM/CHECKS PARAMETERS (GET)

// ------------- GET CHECK PARAMETERS (for specific checktype)
var systemChecktypeParametersGetCmd = &cobra.Command{
	Use:   "parameters",
	Short: "Show all default parameters for a specific checktype.",
	RunE: func(cmd *cobra.Command, args []string) error {
		chosenType := cmd.Flag("type").Value.String()

		checktypeResult, err := Level27Client.SystemCheckTypeGet(chosenType)
		if err != nil {
			return err
		}

		outputFormatTable(checktypeResult.ServiceType.Parameters, []string{"NAME", "DESCRIPTION", "DEFAULT_VALUE"}, []string{"Name", "Description", "DefaultValue"})
		return nil
	},
}

// #endregion

//------------------------------------------------- SYSTEM/CHECKS SPECIFIC (DESCRIBE / DELETE / UPDATE) ----------------------------------
// #region SYSTEM/CHECKS (DESCRIBE / DELETE / UPDATE)

// -------------- GET DETAILS FROM A CHECK
var systemCheckGetSingleCmd = &cobra.Command{
	Use:   "describe [systemID] [checkID]",
	Short: "Get detailed info about a specific check.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		//check for valid system ID
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		//check for valid system checkID
		checkID, err := checkSingleIntID(args[1], "check")
		if err != nil {
			return err
		}

		check, err := Level27Client.SystemCheckDescribe(systemID, checkID)
		if err != nil {
			return err
		}

		outputFormatTemplate(check, "templates/systemCheck.tmpl")
		return nil
	},
}

// -------------- DELETE SPECIFIC CHECK
var systemCheckDeleteCmd = &cobra.Command{
	Use:   "delete [systemID] [checkID]",
	Short: "Delete a specific check from a system",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		//check for valid system ID
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		//check for valid system checkID
		checkID, err := checkSingleIntID(args[1], "check")
		if err != nil {
			return err
		}

		if !optDeleteConfirmed {
			system, err := Level27Client.SystemGetSingle(systemID)
			if err != nil {
				return err
			}

			if !confirmPrompt(fmt.Sprintf("Delete system check %d on system %s (%d)?", checkID, system.Name, system.Id)) {
				return nil
			}
		}

		err = Level27Client.SystemCheckDelete(systemID, checkID)
		return err
	},
}

// -------------- UPDATE SPECIFIC CHECK
var systemCheckUpdateCmd = &cobra.Command{
	Use:   "update [SystemID] [CheckID]",
	Short: "update a specific check from a system",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// check for valid system ID
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		// check for valid check ID
		checkID, err := checkSingleIntID(args[1], "check")
		if err != nil {
			return err
		}

		// get the current data from the check
		currentData, err := Level27Client.SystemCheckDescribe(systemID, checkID)
		if err != nil {
			return err
		}

		// create base of PUT request in JSON (checktype required and cannot be changed)
		updateCheckJson := gabs.New()
		updateCheckJson.Set(currentData.CheckType, "checktype")

		// keep track of possbile parameters for current checktype
		var possibleParameters []string
		// loop over current parameters for the check.
		// if parameter value is not default -> it needs to be sent in put request again.
		for key, value := range currentData.CheckParameters {
			// put each possible parrameter in array for later
			possibleParameters = append(possibleParameters, key)

			if !value.Default {
				updateCheckJson.Set(value.Value, key)
			}
		}

		// check wich parameters the user gave in.
		// also check if way of using parameter flag is correct
		customParamaterDict, err := SplitCustomParameters(systemDynamicParams)
		if err != nil {
			return err
		}

		// check for each given parameter if its one of the possible parameters
		// if parameter = valid -> add key/value to json object for put request
		for givenParameter, givenValue := range customParamaterDict {
			var isValidParameter bool = false
			for i := range possibleParameters {
				if givenParameter == possibleParameters[i] {
					isValidParameter = true
					updateCheckJson.Set(givenValue, givenParameter)
				}
			}

			if !isValidParameter {
				return fmt.Errorf("given parameter key: '%v' is not valid for checktype %v", givenParameter, currentData.CheckType)
			}
		}

		//log.Print(updateCheckJson.StringIndent(""," "))
		err = Level27Client.SystemCheckUpdate(systemID, checkID, updateCheckJson)
		return err
	},
}

// #endregion

//------------------------------------------------- MONITORING ON SPECIFIC SYSTEM ----------------------------------------------
// ---- MONITORING COMMAND
var systemMonitoringCmd = &cobra.Command{
	Use:     "monitoring",
	Short:   "Turn the monitoring for a system on or off.",
	Example: "lvl system monitoring on MySystemName",
}

// ---- MONITORING ON
var systemMonitoringOnCmd = &cobra.Command{
	Use:     "on",
	Short:   "Turn on the monitoring for a system.",
	Example: "lvl system monitoring on MySystemName",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		//search for sytsemID based on name
		systemId, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		_, err = Level27Client.SystemAction(systemId, "enable_monitoring")
		if err != nil {
			return err
		}

		log.Printf("Monitoring is turned on for system: '%v'", args[0])
		return nil
	},
}

// ---- MONITORING OFF
var systemMonitoringOffCmd = &cobra.Command{
	Use:     "off",
	Short:   "Turn off the monitoring for a system.",
	Example: "lvl system monitoring off MySystemName",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		//search for sytsemID based on name
		systemId, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		_, err = Level27Client.SystemAction(systemId, "disable_monitoring")
		if err != nil {
			return err
		}

		log.Printf("Monitoring is turned off for system: '%v'", args[0])
		return nil
	},
}

//------------------------------------------------- ACTIONS ON SPECIFIC SYSTEM ----------------------------------------------
// #region SYSTEM ACTIONS

var systemActionsCmd = &cobra.Command{
	Use:   "actions",
	Short: "Actions for systems such as rebooting",
}

var systemActionsStartCmd = &cobra.Command{
	Use:  "start",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error { return runAction("start", args) },
}

var systemActionsStopCmd = &cobra.Command{
	Use:  "stop",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error { return runAction("stop", args) },
}

var systemActionsShutdownCmd = &cobra.Command{
	Use:  "shutdown",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error { return runAction("shutdown", args) },
}

var systemActionsRebootCmd = &cobra.Command{
	Use:  "reboot",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error { return runAction("reboot", args) },
}

var systemActionsResetCmd = &cobra.Command{
	Use:  "reset",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error { return runAction("reset", args) },
}

var systemActionsEmergencyPowerOffCmd = &cobra.Command{
	Use:  "emergencyPowerOff",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error { return runAction("emergencyPowerOff", args) },
}

var systemActionsDeactivateCmd = &cobra.Command{
	Use:  "deactivate",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error { return runAction("deactivate", args) },
}

var systemActionsActivateCmd = &cobra.Command{
	Use:  "activate",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error { return runAction("activate", args) },
}

var systemActionsAutoInstallCmd = &cobra.Command{
	Use:  "autoInstall",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error { return runAction("autoInstall", args) },
}

func runAction(action string, args []string) error {
	id, err := resolveSystem(args[0])
	if err != nil {
		return err
	}

	_, err = Level27Client.SystemAction(id, action)
	return err
}

// #endregion

//------------------------------------------------- SYSTEM/COOKBOOKS TOPLEVEL (GET error / CREATE) ----------------------------------
// ---------------- MAIN return COMMAND (cookbooks)
var systemCookbookCmd = &cobra.Command{
	Use:   "cookbooks",
	Short: "Manage systems cookbooks",
}

// #region SYSTEM/COOKBOOKS TOPLEVEL (GET / ADD )

// ---------- GET COOKBOOKS
var systemCookbookGetCmd = &cobra.Command{
	Use:   "get [system ID]",
	Short: "Gets a list of all cookbooks from a system.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// check for valid system ID
		id, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		cookbooks, err := Level27Client.SystemCookbookGetList(id)
		if err != nil {
			return err
		}

		outputFormatTable(cookbooks, []string{"ID", "COOKBOOKTYPE", "STATUS"}, []string{"Id", "CookbookType", "Status"})
		return nil
	},
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
	RunE: func(cmd *cobra.Command, args []string) error {
		//checking for valid system ID
		systemId, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		// get information about the current chosen system [systemID]
		currentSystem, err := Level27Client.SystemGetSingle(systemId)
		if err != nil {
			return err
		}

		currentSystemOS := fmt.Sprintf("%v %v", currentSystem.OperatingSystemVersion.OsName, currentSystem.OperatingSystemVersion.OsVersion)

		// get the user input from the type flag (cookbooktype)
		inputType := cmd.Flag("type").Value.String()

		// get all current data for the chosen cookbooktype
		cookbooktypeData, _, err := Level27Client.SystemCookbookTypeGet(inputType)
		if err != nil {
			return err
		}

		// // create base of json container, will be used for post request and eventually filled with custom parameters
		cookbookRequest := l27.CookbookRequest{
			Cookbooktype:       inputType,
			Cookbookparameters: map[string]interface{}{},
		}

		// when user wants to use custom parameters
		if cmd.Flag("parameters").Changed {

			// split the slice of customparameters set by user into key/value pairs. also check if declaration method is used correctly (-p key=value).
			customParameterDict, err := SplitCustomParameters(systemDynamicParams)
			if err != nil {
				return err
			}

			checkForValidCookbookParameter(customParameterDict, cookbooktypeData, currentSystemOS, &cookbookRequest)
		}

		cookbook, err := Level27Client.SystemCookbookAdd(systemId, &cookbookRequest)
		if err != nil {
			return err
		}

		log.Printf("Cookbook: '%v' succesfully added!", cookbook.CookbookType)

		//apply changes to cookbooks
		err = Level27Client.SystemCookbookChangesApply(systemId)
		return err
	},
}

// #endregion

//------------------------------------------------- SYSTEM/COOKBOOKS PARAMETERS GET ----------------------------------
// #region SYSTEM/COOKBOOKS PARAMETERS (GET)

// ----------- GET COOKBOOKTYPE PARAMETERS
// seperate command used to see wich parameters can be used for a specific cookbooktype. also shows the description and default values
var SystemCookbookTypesGetCmd = &cobra.Command{
	Use:   "parameters",
	Short: "Show all default parameters for a specific cookbooktype.",
	RunE: func(cmd *cobra.Command, args []string) error {

		// get the user input from the type flag
		inputType := cmd.Flag("type").Value.String()

		// Get request to get all cookbooktypes data
		validCookbooktype, _, err := Level27Client.SystemCookbookTypeGet(inputType)
		if err != nil {
			return err
		}

		outputFormatTable(validCookbooktype.CookbookType.Parameters, []string{"NAME", "DESCRIPTION", "DEFAULT_VALUE"}, []string{"Name", "Description", "DefaultValue"})
		return nil
	},
}

// #endregion

//------------------------------------------------- SYSTEM/COOKBOOKS SPECIFIC (DESCRIBE / DELETE / UPDATE) ----------------------------------
// #region SYSTEM/COOKBOOKS SPECIFIC (DESCRIBE / DELETE / UPDATE)

// ---------------- DESCRIBE
var systemCookbookDescribeCmd = &cobra.Command{
	Use:   "describe <system> <cookbook>",
	Short: "show detailed info about a cookbook on a system",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// check for valid system id
		systemId, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		cookbookId, err := resolveSystemCookbook(systemId, args[1])
		if err != nil {
			return err
		}

		result, err := Level27Client.SystemCookbookDescribe(systemId, cookbookId)
		if err != nil {
			return err
		}

		outputFormatTemplate(result, "templates/systemCookbook.tmpl")
		return nil
	},
}

// ---------------- DELETE
var systemCookbookDeleteCmd = &cobra.Command{
	Use:   "delete [systemID] [cookbookID]",
	Short: "delete a cookbook from a system.",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemId, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		cookbookId, err := resolveSystemCookbook(systemId, args[1])
		if err != nil {
			return err
		}

		if !optDeleteConfirmed {
			cookbook, err := Level27Client.SystemCookbookDescribe(systemId, cookbookId)
			if err != nil {
				return err
			}

			if !confirmPrompt(fmt.Sprintf("Delete system cookbook %s (%d) on system %s (%d)?", cookbook.CookbookType, cookbook.Id, cookbook.System.Name, cookbook.System.Id)) {
				return nil
			}
		}

		err = Level27Client.SystemCookbookDelete(systemId, cookbookId)
		if err != nil {
			return err
		}

		//apply changes
		err = Level27Client.SystemCookbookChangesApply(systemId)
		return err
	},
}

// ---------------- UPDATE
var systemCookbookUpdateCmd = &cobra.Command{
	Use:   "update [systemID] [cookbookID]",
	Short: "update existing cookbook from a system",
	Example: "lvl system cookbooks update [systemID] [cookbookID] {-p}.\nSINGLE PARAMETER:		-p waf=true  \nMULTIPLE PARAMETERS:		-p waf=true -p timeout=200  \nMULTIPLE VALUES:		-p versions=''7, 5.4'' OR -p versions=7,5.4 (seperated by comma)",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// check for valid system id
		systemId, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		cookbookId, err := resolveSystemCookbook(systemId, args[1])
		if err != nil {
			return err
		}

		// get current data from the current installed cookbooktype
		currentCookbookData, err := Level27Client.SystemCookbookDescribe(systemId, cookbookId)
		if err != nil {
			return err
		}

		// get current data from the chosen system
		currentSystemData, err := Level27Client.SystemGetSingle(systemId)
		if err != nil {
			return err
		}

		currentSystem := fmt.Sprintf("%v %v", currentSystemData.OperatingSystemVersion.OsName, currentSystemData.OperatingSystemVersion.OsVersion)

		// get all standard data that belongs to this cookbooktype in general (parameters / parameteroptions..).
		cookbookData, _, err := Level27Client.SystemCookbookTypeGet(currentCookbookData.CookbookType)
		if err != nil {
			return err
		}

		// create base form of json for PUT request (cookbooktype is non-editable)
		baseRequestData := l27.CookbookRequest{
			Cookbooktype:       currentCookbookData.CookbookType,
			Cookbookparameters: map[string]interface{}{},
		}

		// loop over current data and check if values are default. (default values dont need to be in put request)
		for key, value := range currentCookbookData.CookbookParameters.Map {
			if !value.Default {
				baseRequestData.Cookbookparameters[key] = value.Value
			}
		}

		//check if parameter flag is used correctly
		//split key/value pairs from parameter flag
		customParameterDict, err := SplitCustomParameters(systemDynamicParams)
		if err != nil {
			return err
		}

		// check for each set parameter if its one of the possible parameters for this cookbooktype
		// als checks if values are valid in case of selectable parameter
		checkForValidCookbookParameter(customParameterDict, cookbookData, currentSystem, &baseRequestData)

		err = Level27Client.SystemCookbookUpdate(systemId, cookbookId, &baseRequestData)
		if err != nil {
			return err
		}

		// aplly changes to cookbooks
		err = Level27Client.SystemCookbookChangesApply(systemId)
		return err
	},
}

// #endregion

//-------------------------------------------------// COOKBOOK HELPER FUNCTIONS //-------------------------------------------------
// #region COOKBOOK HELPER FUNCTIONS

// checks if a given parameter is valid for the specific cookbooktype.
// also checks if given values are valid for chosen parameter or compatible with current system
func checkForValidCookbookParameter(customParameters map[string]interface{}, allCookbookData l27.CookbookType, currenSystemOs string, currenRequestData *l27.CookbookRequest) error {

	// for each custom set parameter, check if its one of the possible parameters for the current cookbooktype
	for givenParameter, givenValue := range customParameters {
		var isValidParameter bool = false

		for _, possibleParameter := range allCookbookData.CookbookType.Parameters {
			if givenParameter == possibleParameter.Name {
				isValidParameter = true

				//check if chosen parameter is of type "select" (needs extra validation)
				if possibleParameter.Type == "select" {

					AllParameterOptions := allCookbookData.CookbookType.ParameterOptions

					// check type of given value (selectable parameters needs to be post in array type)
					givenValueType := reflect.ValueOf(givenValue)

					// if type of given interface value is array or slice we need to convert the interface to a go slice
					if givenValueType.Kind() == reflect.Array || givenValueType.Kind() == reflect.Slice {
						rawValues := reflect.ValueOf(givenValue)

						// need to convert interface to a go slice
						givenValuesSlice := make([]interface{}, rawValues.Len())
						for i := 0; i < rawValues.Len(); i++ {
							givenValuesSlice[i] = rawValues.Index(i).Interface()
						}

						for _, value := range givenValuesSlice {

							valueString := fmt.Sprintf("%v", value)
							// is value valid for given parameter
							isExclusive, err := CheckCBValueForParameter(valueString, AllParameterOptions[givenParameter], givenParameter, currenSystemOs)
							if err != nil {
								return err
							}

							if isExclusive {
								return fmt.Errorf("value '%v' is not possible for multiselect", value)
							}
						}

						currenRequestData.Cookbookparameters[givenParameter] = givenValuesSlice

					} else {
						// only a single value was given by the user for the parameter
						valueString := fmt.Sprintf("%v", givenValue)
						CheckCBValueForParameter(valueString, AllParameterOptions[givenParameter], givenParameter, currenSystemOs)
						//key has one value but needs to be sent in array type
						var values []interface{}
						values = append(values, valueString)

						currenRequestData.Cookbookparameters[givenParameter] = values
					}
				} else {
					currenRequestData.Cookbookparameters[givenParameter] = givenValue
				}

			}
		}

		// when parameter is not valid for cookbooktype
		if !isValidParameter {
			return fmt.Errorf("given parameter key: '%v' NOT valid for cookbooktype %v", givenParameter, allCookbookData.CookbookType.Name)
		}
	}

	return nil
}

// check a value if its a valid option for the given parameter for the cookbook.
// also do checks on compatibility with system and exlusivity
func CheckCBValueForParameter(value string, options l27.CookbookParameterOptionValue, givenParameter string, currentSystemOs string) (bool, error) {
	parameterOptionValue, found := options[value]

	// check if given value is one of the options for the chosen selectable parameter
	if !found {
		return false, fmt.Errorf("given value: '%v' NOT a valid option for parameter '%v'", value, givenParameter)
	}

	//  loop over all possible OS version and check if the chosen value is compatible with current system
	var isCompatibleWithSystem bool = false
	for _, osVersion := range parameterOptionValue.OperatingSystemVersions {

		if osVersion.Name == currentSystemOs {
			isCompatibleWithSystem = true

		}
	}

	// error when value required OS version doesnt equal current system OS version
	if !isCompatibleWithSystem {
		return false, fmt.Errorf("given %v: '%v' NOT compatible with current system: %v", givenParameter, value, currentSystemOs)
	}

	return parameterOptionValue.Exclusive, nil
}

// #endregion

// #endregion

//------------------------------------------------- SYSTEMS/GROUPS (GET / ADD  / DELETE)-------------------------------------------------
// ---------------- MAIN COMMAND (groups)
var SystemSystemgroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "Manage a system's groups.",
}

// #region SYSTEMS/GROUPS (GET / ADD  / DELETE)

// ---------------- GET GROUPS
var SystemSystemgroupsGetCmd = &cobra.Command{
	Use:   "get [systemID]",
	Short: "Show list of all groups from a system.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		//check for valid systemID
		systemId, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		groups, err := Level27Client.SystemSystemgroupsGet(systemId)
		if err != nil {
			return err
		}

		outputFormatTable(groups, []string{"ID", "NAME"}, []string{"ID", "Name"})
		return nil
	},
}

// ---------------- LINK SYSTEM TO A GROUP (ADD)
var SystemSystemgroupsAddCmd = &cobra.Command{
	Use:   "add [systemID] [systemgroupID]",
	Short: "Link a system with a systemgroup.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// check for valid systemID
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		// check for valid groupID type (int)
		groupId, err := resolveSystemgroup(args[1])
		if err != nil {
			return err
		}

		jsonRequest := gabs.New()
		jsonRequest.Set(groupId, "systemgroup")
		err = Level27Client.SystemSystemgroupsAdd(systemID, jsonRequest)
		if err != nil {
			return err
		}

		log.Printf("System succesfully linked to systemgroup!")
		return nil
	},
}

// ---------------- UNLINK SYSTEM FROM A GROUP (DELETE)
var SystemSystemgroupsRemoveCmd = &cobra.Command{
	Use:   "remove [systemID] [systemgroupID]",
	Short: "Unlink a system from a systemgroup.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// check for valid systemId
		systemId, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		// check for valid systemgroupId
		groupId, err := resolveSystemgroup(args[1])
		if err != nil {
			return err
		}

		err = Level27Client.SystemSystemgroupsRemove(systemId, groupId)
		if err != nil {
			return err
		}

		log.Printf("System succesfully removed from systemgroup!")
		return nil
	},
}

// #endregion

//------------------------------------------------- SYSTEMS / SSH KEYS (GET / ADD / DELETE)

var systemSshKeysCmd = &cobra.Command{
	Use: "sshkeys",
}

// #region SYSTEMS/SHH KEYS (GET / ADD / DELETE)

// --- GET
var systemSshKeysGetCmd = &cobra.Command{
	Use: "get",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		keys, err := Level27Client.SystemGetSshKeys(id, optGetParameters)
		if err != nil {
			return err
		}

		outputFormatTable(keys, []string{"ID", "DESCRIPTION", "STATUS", "FINGERPRINT"}, []string{"ID", "Description", "ShsStatus", "Fingerprint"})
		return nil
	},
}

// --- ADD
var systemSshKeysAddCmd = &cobra.Command{
	Use: "add",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		keyName := args[1]
		keyID, err := strconv.Atoi(keyName)
		if err != nil {
			user := viper.GetInt("user_id")
			org := viper.GetInt("org_id")
			system, err := Level27Client.LookupSystemNonAddedSshkey(systemID, org, user, keyName)
			if err != nil {
				return err
			}

			if system == nil {
				existing, err := Level27Client.LookupSystemSshkey(systemID, keyName)
				if err != nil {
					return err
				}

				if existing != nil {
					return errors.New("SSH key already exists on system")
				}

				return fmt.Errorf("unable to find SSH key to add: '%s'", keyName)
			}

			keyID = system.Id
		}

		_, err = Level27Client.SystemAddSshKey(systemID, keyID)
		if err != nil {
			return err
		}

		log.Printf("SSH key added succesfully!")
		return nil
	},
}

// --- DELETE
var systemSshKeysRemoveCmd = &cobra.Command{
	Use: "remove",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		keyName := args[1]
		keyID, err := strconv.Atoi(keyName)
		if err != nil {
			existing, err := Level27Client.LookupSystemSshkey(systemID, keyName)
			if err != nil {
				return err
			}

			if existing == nil {
				return fmt.Errorf("unable to find SSH key to remove: %s", keyName)
			}

			keyID = existing.ID
		}

		err = Level27Client.SystemRemoveSshKey(systemID, keyID)
		return err
	},
}

// #endregion

// SYSTEM SSH
var systemSshCmd = &cobra.Command{
	Use:   "ssh [system] [ssh args]",
	Short: "Connect to a system via SSH, automatically adding SSH keys to the system if necessary",

	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		favoriteKeyID := viper.GetInt("ssh_favoritekey")
		if favoriteKeyID == 0 {
			return fmt.Errorf("no favorite SSH key configured. Use 'lvl sshkey favorite' to configure one")
		}

		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		// We need to do two things:
		// 1. Make sure we have an SSH key on the system.
		// 2. Fetch the host to pass in the ssh command.
		// We send these as concurrent tasks to reduce latency on the command.

		taskSshKey := taskRunVoid(func() error {
			return waitEnsureSshKey(systemID, favoriteKeyID)
		})

		taskSshHost := taskRun(func() (string, error) {
			return sshResolveHost(systemID)
		})

		sshHost := <-taskSshHost
		if sshHost.Error != nil {
			return sshHost.Error
		}

		err = <-taskSshKey
		if err != nil {
			return err
		}

		sshArgs := []string{fmt.Sprintf("root@%s", sshHost.Result)}
		sshArgs = append(sshArgs, args[1:]...)

		return tailExecProcess("ssh", sshArgs)
	},
}

// Ensure the given SSH key is available and 'ok' on a system.
func waitEnsureSshKey(systemID int, sshKeyId int) error {
	_, err := Level27Client.SystemSshKeysGetSingle(systemID, sshKeyId)
	if err == nil {
		// No error, so key exists.
		return nil
	}

	// Error, might indicate SSH key doesn't exist yet.
	_, ok := err.(l27.ErrorResponse)
	if !ok {
		// Not an API error, could be network failure or something instead, abort.
		return err
	}

	// TODO: check error code above, isn't currently correct thanks to PL-7611
	// For now we assume it's just a 404, so try to add the SSH key.

	err = waitAddSshKey(systemID, sshKeyId)
	return err
}

// Add an SSH key to a system, waiting for the status to change to 'ok'.
func waitAddSshKey(systemID int, sshKeyID int) error {
	key, err := Level27Client.SystemAddSshKey(systemID, sshKeyID)
	if err != nil {
		return err
	}

	// Use polling to wait for the SSH key to be fully set-up on the system.

	// Growing wait times on further iterations.
	waitTimes := []int{1, 1, 2, 3, 5, 8, 13, 21, 34}
	for _, wait := range waitTimes {
		time.Sleep(time.Duration(wait * int(time.Second)))

		key, err = Level27Client.SystemSshKeysGetSingle(systemID, key.ID)
		if err != nil {
			return err
		}

		if key.ShsStatus == "ok" {
			return nil
		}
	}

	return fmt.Errorf("timeout waiting for added key to change to 'ok' status")
}

// Resolve the hostname to SSH into a system.
func sshResolveHost(systemID int) (string, error) {
	system, err := Level27Client.SystemGetSingle(systemID)
	if err != nil {
		return "", err
	}

	ips, err := net.LookupIP(system.Fqdn)
	if err == nil && len(ips) > 0 {
		// FQDN resolves, pass it to the ssh command.
		return system.Fqdn, nil
	}

	sort.Slice(system.Networks, func(i int, j int) bool {
		netA := system.Networks[i]
		netB := system.Networks[j]

		return netA.NetPublic && netB.NetInternal
	})

	for _, net := range system.Networks {
		for _, ip := range net.Ips {
			if ip.PublicIpv4 != "" {
				return ip.PublicIpv4, nil
			}

			if ip.Ipv4 != "" {
				return ip.Ipv4, nil
			}
		}
	}

	// Couldn't find anything.
	return "", fmt.Errorf("unable to find a suitable address to connect to on system '%s' (%d)", system.Name, system.Id)
}

// SYSTEM SCP
var systemScpCommand = &cobra.Command{
	Use:     "scp [system1:]file1 ... [system2:]file2",
	Short:   "Copy files to/from the system using scp",
	Long:    "Uses the same syntax as regular scp. Arguments are passed through, but host names (before the :) are interpreted as system names/IDs and resolved. To pass flags through to scp, put them after a --",
	Example: "lvl system scp foo.txt mySystem:~/foo.txt",

	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		favoriteKeyID := viper.GetInt("ssh_favoritekey")
		if favoriteKeyID == 0 {
			return fmt.Errorf("no favorite SSH key configured. Use 'lvl sshkey favorite' to configure one")
		}

		// This code is quite complex to be able to be as optimally concurrent (and fast) as possible.
		// Basically how it works:
		// For each argument, we need to resolve the system if it's a system:file argument
		// These resolves are done concurrently. They send back the new value of the arg when they're done.
		// We also need to make sure SSH keys are added, this goes via another set of channels to also be concurrent.

		// Goroutine to asynchronously add SSH keys to systems while we go through resolving systems down below.
		// We need this to avoid trying to add an SSH key to the same system twice, causing race conditions.
		keyAddChannel := make(chan int)
		keyDone := taskRunVoid(func() error {
			var group errgroup.Group
			// Map of systems we're already handling SSH keys on, to avoid running them twice.
			systemsEnsured := map[int]bool{}
			for systemID := range keyAddChannel {
				sysID := systemID
				if _, ok := systemsEnsured[systemID]; ok {
					continue
				}

				systemsEnsured[systemID] = true
				group.Go(func() error { return waitEnsureSshKey(sysID, favoriteKeyID) })
			}

			return group.Wait()
		})

		var systemArgsGroup errgroup.Group
		systemArgsChannel := make(chan tuple2[int, string])

		// Copy input arguments to pass them through to scp.
		// We will modify the ones that are remote files to replace system hosts with the real IP/domain
		scpArgs := append([]string{}, args...)

		argTask := taskRunVoid(func() error {
			// Go over remote arguments, and resolve the system.
			for i, arg := range args {
				split := strings.SplitN(arg, ":", 2)
				if len(split) == 1 {
					// No host specified, so local file or flag or something.
					// TODO: this means of parsing mostly works, but it means that any flag parameters with a colon in them
					// will be interpreted as a remote file.
					// It might be a good idea to manually pass-through flag args for well-known flags.
					continue
				}

				// Do this all in parallel with goroutines to avoid chaining latency, nice and fast.
				ii := i
				systemArgsGroup.Go(func() error {
					system := split[0]
					file := split[1]

					systemID, err := resolveSystem(system)
					if err != nil {
						return err
					}

					// Send ID to SSH key channel so the SSH key gets added.
					keyAddChannel <- systemID

					host, err := sshResolveHost(systemID)
					if err != nil {
						return err
					}

					// Send arg index and new value so the value gets replaced.
					systemArgsChannel <- makeTuple2(ii, fmt.Sprintf("root@%s:%s", host, file))
					return nil
				})
			}

			// Wait for all systems to finish resolving.
			err := systemArgsGroup.Wait()
			close(systemArgsChannel)
			return err
		})

		// Update all the args that need updating from the above loop.
		for tuple := range systemArgsChannel {
			scpArgs[tuple.Item1] = tuple.Item2
		}

		// Handle errors from argument processing.
		if err := <-argTask; err != nil {
			return err
		}

		// All args have been processed, so also all SSH keys have been dispatched at least.
		close(keyAddChannel)

		// Wait for all SSH keys to be available on systems.
		err := <-keyDone
		if err != nil {
			return err
		}

		return tailExecProcess("scp", scpArgs)
	},
}

// NETWORKS

var systemNetworkCmd = &cobra.Command{
	Use: "network",
}

var systemNetworkGetCmd = &cobra.Command{
	Use:   "get [system]",
	Short: "Get list of networks on a system",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		system, err := Level27Client.SystemGetSingle(systemID)
		if err != nil {
			return err
		}

		outputFormatTableFuncs(system.Networks, []string{"ID", "Network ID", "Type", "Name", "MAC", "IPs"}, []interface{}{"ID", "NetworkID", func(net l27.SystemNetwork) string {
			if net.NetPublic {
				return "public"
			}
			if net.NetCustomer {
				return "customer"
			}
			if net.NetInternal {
				return "internal"
			}
			return ""
		}, "Name", "Mac", func(net l27.SystemNetwork) string {
			return strconv.Itoa(len(net.Ips))
		}})

		return nil
	},
}

var systemNetworkDescribeCmd = &cobra.Command{
	Use:   "describe [system]",
	Short: "Display detailed information about all networks on a system",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		system, err := Level27Client.SystemGetSingle(systemID)
		if err != nil {
			return err
		}

		networks, err := Level27Client.SystemGetHasNetworks(systemID)
		if err != nil {
			return err
		}

		outputFormatTemplate(DescribeSystemNetworks{
			Networks:    system.Networks,
			HasNetworks: networks,
		}, "templates/systemNetworks.tmpl")

		return nil
	},
}

var systemNetworkAddCmd = &cobra.Command{
	Use:   "add [system] [network]",
	Short: "Add a network to a system",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		networkID, err := resolveNetwork(args[1])
		if err != nil {
			return err
		}

		_, err = Level27Client.SystemAddHasNetwork(systemID, networkID)
		if err != nil {
			return err
		}

		log.Printf("Network succesfully added to system!")
		return nil
	},
}

var systemNetworkRemoveCmd = &cobra.Command{
	Use:   "remove [system] [network]",
	Short: "Remove a network from a system",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		networkID, err := resolveSystemHasNetwork(systemID, args[1])
		if err != nil {
			return err
		}

		err = Level27Client.SystemRemoveHasNetwork(systemID, networkID)
		if err != nil {
			return err
		}

		log.Printf("Network succesfully removed from network!")
		return nil
	},
}

var systemNetworkIpCmd = &cobra.Command{
	Use:   "ip",
	Short: "Manage IP addresses on network connections",
}

var systemNetworkIpGetCmd = &cobra.Command{
	Use:   "get [system] [network]",
	Short: "Get all IP addresses for a system network",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		networkID, err := resolveSystemHasNetwork(systemID, args[1])
		if err != nil {
			return err
		}

		ips, err := Level27Client.SystemGetHasNetworkIps(systemID, networkID)
		if err != nil {
			return err
		}

		outputFormatTableFuncs(ips, []string{"ID", "Public IP", "IP", "Hostname", "Status"}, []interface{}{"ID", func(i l27.SystemHasNetworkIp) string {
			if i.PublicIpv4 != "" {
				i, _ := strconv.Atoi(i.PublicIpv4)
				if i == 0 {
					return ""
				} else {
					return ipv4IntToString(i)
				}
			} else if i.PublicIpv6 != "" {
				ip := net.ParseIP(i.PublicIpv6)
				return fmt.Sprint(ip)
			} else {
				return ""
			}
		},
			func(i l27.SystemHasNetworkIp) string {
				if i.Ipv4 != "" {
					i, _ := strconv.Atoi(i.Ipv4)
					if i == 0 {
						return ""
					} else {
						return ipv4IntToString(i)
					}
				} else if i.Ipv6 != "" {
					ip := net.ParseIP(i.Ipv6)
					return fmt.Sprint(ip)
				} else {
					return ""
				}
			}, "Hostname", "Status"})

		return nil
	},
}

var systemNetworkIpAddHostname string

var systemNetworkIpAddCmd = &cobra.Command{
	Use:   "add [system] [network] [address]",
	Short: "Add IP address to a system network",
	Long:  "Adds an IP address to a system network. Address can be either IPv4 or IPv6. The special values 'auto' and 'auto-v6' automatically fetch an unused address to use.",

	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		system, err := Level27Client.SystemGetSingle(systemID)
		if err != nil {
			return err
		}

		hasNetworkID, err := resolveSystemHasNetwork(systemID, args[1])
		if err != nil {
			return err
		}

		network, err := Level27Client.GetSystemHasNetwork(systemID, hasNetworkID)
		if err != nil {
			return err
		}

		networkID := network.Network.ID
		address := args[2]

		if address == "auto" || address == "auto-v6" {
			located, err := Level27Client.NetworkLocate(networkID)
			if err != nil {
				return err
			}

			var choices []string
			if address == "auto" {
				choices = located.Ipv4
			} else {
				choices = located.Ipv6
			}

			if len(choices) == 0 {
				return errors.New("unable to find a free IP address")
			}

			address = choices[0]
		}

		var data l27.SystemHasNetworkIpAdd
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

		_, err = Level27Client.SystemAddHasNetworkIps(systemID, hasNetworkID, data)
		return err
	},
}

var systemNetworkIpRemoveCmd = &cobra.Command{
	Use:   "remove [system] [network] [address | id]",
	Short: "Remove IP address from a system network",

	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		hasNetworkID, err := resolveSystemHasNetwork(systemID, args[1])
		if err != nil {
			return err
		}

		ipID, err := resolveSystemHasNetworkIP(systemID, hasNetworkID, args[2])
		if err != nil {
			return err
		}

		err = Level27Client.SystemRemoveHasNetworkIps(systemID, hasNetworkID, ipID)
		return err
	},
}

var systemNetworkIpUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update settings on a system network IP",

	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		settings, err := loadMergeSettings(updateSettingsFile, updateSettings)
		if err != nil {
			return err
		}

		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		hasNetworkID, err := resolveSystemHasNetwork(systemID, args[1])
		if err != nil {
			return err
		}

		ipID, err := resolveSystemHasNetworkIP(systemID, hasNetworkID, args[2])
		if err != nil {
			return err
		}

		ip, err := Level27Client.SystemGetHasNetworkIp(systemID, hasNetworkID, ipID)
		if err != nil {
			return err
		}

		ipPut := l27.SystemHasNetworkIpPut{
			Hostname: ip.Hostname,
		}

		data := mergeSettingsWithEntity(ipPut, settings)

		err = Level27Client.SystemHasNetworkIpUpdate(systemID, hasNetworkID, ipID, data)
		return err
	},
}

// VOLUMES

// SYSTEM VOLUME
var systemVolumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Commands to manage volumes",
}

// SYSTEM VOLUME GET
var systemVolumeGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get all volumes on a system",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		volumes, err := Level27Client.SystemGetVolumes(systemID, optGetParameters)
		if err != nil {
			return err
		}

		outputFormatTable(
			volumes,
			[]string{"ID", "Name", "Status", "Space", "UID", "AutoResize", "DeviceName"},
			[]string{"ID", "Name", "Status", "Space", "UID", "AutoResize", "DeviceName"})

		return nil
	},
}

// SYSTEM VOLUME CREATE
var systemVolumeCreateName string
var systemVolumeCreateSpace int
var systemVolumeCreateOrganisation string
var systemVolumeCreateAutoResize bool
var systemVolumeCreateDeviceName string

var systemVolumeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new volume for a system",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		organisationID, err := resolveOrganisation(systemVolumeCreateOrganisation)
		if err != nil {
			return err
		}

		create := l27.VolumeCreate{
			Name:         systemVolumeCreateName,
			Space:        systemVolumeCreateSpace,
			Organisation: organisationID,
			System:       systemID,
			AutoResize:   systemVolumeCreateAutoResize,
			DeviceName:   systemVolumeCreateDeviceName,
		}

		_, err = Level27Client.VolumeCreate(create)
		return err
	},
}

// SYSTEM VOLUME UNLINK
var systemVolumeUnlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "Unlink a volume from a system",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		volumeID, err := resolveSystemVolume(systemID, args[1])
		if err != nil {
			return err
		}

		_, err = Level27Client.VolumeUnlink(volumeID, systemID)
		return err
	},
}

// SYSTEM VOLUME LINK
var systemVolumeLinkCmd = &cobra.Command{
	Use:   "link [system] [volume] [device name]",
	Short: "Link a volume to a system",

	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		// To resolve from name -> ID we need the volume group
		// Easiest way to get that is by getting the volume group ID from the first volume on the system.
		volumes, err := Level27Client.SystemGetVolumes(systemID, l27.CommonGetParams{})
		if err != nil {
			return err
		}

		volumeGroupID := volumes[0].Volumegroup.ID

		volumeID, err := resolveVolumegroupVolume(volumeGroupID, args[1])
		if err != nil {
			return err
		}

		deviceName := args[2]

		_, err = Level27Client.VolumeLink(volumeID, systemID, deviceName)
		return err
	},
}

// SYSTEM VOLUME DELETE
var systemVolumeDeleteForce bool
var systemVolumeDeleteCmd = &cobra.Command{
	Use:   "delete [system] [volume]",
	Short: "Unlink and delete a volume on a system",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		volumeID, err := resolveSystemVolume(systemID, args[1])
		if err != nil {
			return err
		}

		if !systemVolumeDeleteForce {
			volume, err := Level27Client.VolumeGetSingle(volumeID)
			if err != nil {
				return err
			}

			if !confirmPrompt(fmt.Sprintf("Delete volume %s (%d)?", volume.Name, volume.ID)) {
				return nil
			}
		}

		err = Level27Client.VolumeDelete(volumeID)
		return err
	},
}

// SYSTEM VOLUME UPDATE
var systemVolumeUpdateCmd = &cobra.Command{
	Use:   "update [system] [volume]",
	Short: "Update settings on a volume",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		settings, err := loadMergeSettings(updateSettingsFile, updateSettings)
		if err != nil {
			return err
		}

		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		volumeID, err := resolveSystemVolume(systemID, args[1])
		if err != nil {
			return err
		}

		volume, err := Level27Client.VolumeGetSingle(volumeID)
		if err != nil {
			return err
		}

		volumePut := l27.VolumePut{
			Name:         volume.Name,
			DeviceName:   volume.DeviceName,
			Space:        volume.Space,
			Organisation: volume.Organisation.ID,
			AutoResize:   volume.AutoResize,
			Remarks:      volume.Remarks,
			System:       volume.System.Id,
			Volumegroup:  volume.Volumegroup.ID,
		}

		data := utils.RoundTripJson(volumePut).(map[string]interface{})
		data = mergeMaps(data, settings)

		data["organisation"], err = resolveOrganisation(fmt.Sprint(data["organisation"]))
		if err != nil {
			return err
		}

		err = Level27Client.VolumeUpdate(volumeID, data)
		return err
	},
}

type DescribeSystem struct {
	l27.System
	SshKeys                      []l27.SystemSshkey     `json:"sshKeys"`
	InstallSecurityUpdatesString string                 `json:"installSecurityUpdatesString"`
	HasNetworks                  []l27.SystemHasNetwork `json:"hasNetworks"`
	Volumes                      []l27.SystemVolume     `json:"volumes"`
	Checks                       []l27.SystemCheckGet   `json:"checks"`
}

type DescribeSystemNetworks struct {
	Networks    []l27.SystemNetwork    `json:"networks"`
	HasNetworks []l27.SystemHasNetwork `json:"hasNetworks"`
}
