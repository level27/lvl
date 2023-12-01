package cmd

import (
	"fmt"
	"strings"

	"github.com/level27/l27-go"
	"github.com/level27/lvl/utils"
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
	addWaitFlag(systemCreateCmd)
	flags.StringVarP(&systemCreateName, "name", "n", "", "The name you want to give the system")
	flags.StringVarP(&systemCreateFqdn, "Fqdn", "", "", "Valid hostname for the system")
	flags.StringVarP(&systemCreateRemarks, "remarks", "", "", "Remarks (Admin only)")
	flags.Int32VarP(&systemCreateDisk, "disk", "", 0, "Disk (non-editable)")
	flags.Int32VarP(&systemCreateCpu, "cpu", "", 0, "Cpu (Required for Level27 systems)")
	flags.Int32VarP(&systemCreateMemory, "memory", "", 0, "Memory (Required for Level27 systems)")
	flags.StringVarP(&systemCreateManageType, "management", "", "basic", "Managament type (one of basic, professional, enterprise, professional_level27).")
	flags.BoolVarP(&systemCreatePublicNetworking, "publicNetworking", "", true, "For digitalOcean servers always true. (non-editable)")
	flags.StringVarP(&systemCreateImage, "image", "", "", "The ID of a systemimage. (must match selected configuration and zone. non-editable)")
	flags.StringVarP(&systemCreateOrganisation, "organisation", "", "", "The unique ID of an organisation")
	flags.StringVarP(&systemCreateProviderConfig, "config", "", "", "The unique ID of a SystemproviderConfiguration")
	flags.StringVarP(&systemCreateZone, "zone", "", "", "The unique ID of a zone")
	//	flags.StringVarP(&systemCreateSecurityUpdates, "security", "", "", "installSecurityUpdates (default: random POST:1-8, PUT:0-12)") NOT NEEDED FOR CREATE REQUEST
	flags.StringVarP(&systemCreateAutoTeams, "autoTeams", "", "", "A csv list of team ID's")
	flags.StringVarP(&systemCreateExternalInfo, "externalInfo", "", "", "ExternalInfo (required when billableItemInfo entities for an organisation exist in db)")
	flags.Int32VarP(&systemCreateOperatingSystemVersion, "version", "", 0, "The unique ID of an OperatingsystemVersion (non-editable)")
	flags.Int32VarP(&systemCreateParentSystem, "parent", "", 0, "The unique ID of a system (parent system)")
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
	systemActionsCmd.AddCommand(systemActionsHypervisorFailedCmd)
	systemActionsStartMaintenanceCmd.Flags().Int32VarP(&systemActionsStartMaintenanceDuration, "duration", "d", 1440, "How long maintenance should last, in minutes")
	systemActionsCmd.AddCommand(systemActionsStartMaintenanceCmd)
	systemActionsCmd.AddCommand(systemActionsStopMaintenanceCmd)

	// --- UPDATE

	systemCmd.AddCommand(systemUpdateCmd)
	settingsFileFlag(systemUpdateCmd)
	settingString(systemUpdateCmd, updateSettings, "name", "New name for this system")
	settingInt32(systemUpdateCmd, updateSettings, "cpu", "Set amount of CPU cores of the system")
	settingInt32(systemUpdateCmd, updateSettings, "memory", "Set amount of memory in GB of the system")
	settingString(systemUpdateCmd, updateSettings, "managementType", "Set management type of the system")
	settingString(systemUpdateCmd, updateSettings, "organisation", "Set organisation that owns this system. Can be both a name or an ID")
	settingInt32(systemUpdateCmd, updateSettings, "publicNetworking", "")
	settingInt32(systemUpdateCmd, updateSettings, "limitRiops", "Set read IOPS limit")
	settingInt32(systemUpdateCmd, updateSettings, "limitWiops", "Set write IOPS limit")
	settingInt32(systemUpdateCmd, updateSettings, "installSecurityUpdates", "Set security updates mode index")
	settingString(systemUpdateCmd, updateSettings, "remarks", "")
	settingInt32(systemUpdateCmd, updateSettings, "operatingsystemVersion", "")
	settingStringS(systemUpdateCmd, updateSettings, "customerFqdn", "hostname", "")

	// --- Delete

	systemCmd.AddCommand(systemDeleteCmd)
	addWaitFlag(systemDeleteCmd)
	systemDeleteCmd.Flags().BoolVar(&systemDeleteForce, "force", false, "")
	addDeleteConfirmFlag(systemDeleteCmd)
	// #endregion

	//-------------------------------------  SYSTEMS/INTEGRITYCHECKS (get / post / download) --------------------------------------
	addIntegrityCheckCmds(systemCmd, "systems", resolveSystem)

	// ACCESS
	addAccessCmds(systemCmd, "systems", resolveSystem)

	// BILLING
	addBillingCmds(systemCmd, "systems", resolveSystem)

	// JOBS
	addJobCmds(systemCmd, "system", resolveSystem)
}

// Resolve an integer or name domain.
// If the domain is a name, a request is made to resolve the integer ID.
func resolveSystem(arg string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
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
		func(app l27.System) string { return fmt.Sprintf("%s (%d)", app.Name, app.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, err
}

func resolveSystemProviderConfiguration(region l27.IntID, arg string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
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

// ------------------------------------------------- SYSTEM TOPLEVEL (GET / DESCRIBE CREATE) ----------------------------------
// #region SYSTEM TOPLEVEL (GET / DESCRIBE / CREATE)
// ----------------------------------------- GET ---------------------------------------
var systemGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get a list of all curent systems",
	RunE: func(cmd *cobra.Command, args []string) error {
		systems, err := resolveGets(
			args,
			Level27Client.LookupSystem,
			Level27Client.SystemGetSingle,
			Level27Client.SystemGetList)

		if err != nil {
			return err
		}

		outputFormatTable(systems, []string{"ID", "NAME", "STATUS"}, []string{"ID", "Name", "Status"})
		return nil
	},
}

// ----------------------------------------- DESCRIBE ---------------------------------------
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
			jobs, err := Level27Client.EntityJobHistoryGet("system", systemID, l27.PageableParams{})
			if err != nil {
				return err
			}

			system.Jobs = make([]l27.Job, len(jobs))

			for idx, j := range jobs {
				system.Jobs[idx], err = Level27Client.JobHistoryRootGet(
					j.ID,
					l27.JobHistoryGetParams{})

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

// ----------------------------------------- CREATE ---------------------------------------
// vars needed to save flag data.
var systemCreateName, systemCreateFqdn, systemCreateRemarks string
var systemCreateDisk, systemCreateCpu, systemCreateMemory int32
var systemCreateManageType string
var systemCreatePublicNetworking bool
var systemCreateImage, systemCreateOrganisation, systemCreateProviderConfig, systemCreateZone string

var systemCreateAutoTeams, systemCreateExternalInfo string
var systemCreateOperatingSystemVersion, systemCreateParentSystem l27.IntID
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

		if optWait {
			system, err = waitForStatus(
				func() (l27.System, error) { return Level27Client.SystemGetSingle(system.ID) },
				func(s l27.System) string { return s.Status },
				"allocated",
				[]string{"to_create", "creating"},
			)

			if err != nil {
				return fmt.Errorf("waiting on system status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(system, "templates/entities/system/create.tmpl")
		return nil
	},
}

// #endregion

// ------------------------------------------------- SYSTEM SPECIFIC (UPDATE / FORCE DELETE ) ----------------------------------
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
			ID:                          system.ID,
			Name:                        system.Name,
			Type:                        system.Type,
			Cpu:                         system.Cpu,
			Memory:                      system.Memory,
			Disk:                        system.Disk,
			ManagementType:              system.ManagementType,
			Organisation:                system.Organisation.ID,
			SystemImage:                 nil,
			OperatingsystemVersion:      nil,
			SystemProviderConfiguration: nil,
			Zone:                        nil,
			PublicNetworking:            system.PublicNetworking,
			Preferredparentsystem:       nil,
			Remarks:                     system.Remarks,
			InstallSecurityUpdates:      system.InstallSecurityUpdates,
			LimitRiops:                  system.LimitRiops,
			LimitWiops:                  system.LimitWiops,
			CustomerFqdn:                system.Hostname,
		}

		if system.SystemImage != nil {
			systemPut.SystemImage = &system.SystemImage.ID
		}

		if system.Preferredparentsystem != nil {
			systemPut.Preferredparentsystem = &system.Preferredparentsystem.ID
		}

		if system.SystemProviderConfiguration != nil {
			systemPut.SystemProviderConfiguration = &system.SystemProviderConfiguration.ID
		}

		if system.Zone != nil {
			systemPut.Zone = &system.Zone.ID
		}

		if system.OperatingSystemVersion != nil {
			systemPut.OperatingsystemVersion = &system.OperatingSystemVersion.ID
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

		outputFormatTemplate(nil, "templates/entities/system/update.tmpl")
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

			if !confirmPrompt(fmt.Sprintf("Delete system %s (%d)?", system.Name, system.ID)) {
				return nil
			}
		}

		if systemDeleteForce {
			err = Level27Client.SystemDeleteForce(systemID)
		} else {
			err = Level27Client.SystemDelete(systemID)
		}

		if err != nil {
			return err
		}

		if optWait {
			err = waitForDelete(
				func() (l27.System, error) { return Level27Client.SystemGetSingle(systemID) },
				func(s l27.System) string { return s.Status },
				[]string{"to_delete"},
			)

			if err != nil {
				return fmt.Errorf("waiting on system status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(nil, "templates/entities/system/delete.tmpl")
		return nil
	},
}

// #endregion

//------------------------------------------------- ACTIONS ON SPECIFIC SYSTEM ----------------------------------------------
// #region SYSTEM ACTIONS

var systemActionsCmd = &cobra.Command{
	Use:   "actions",
	Short: "Actions for systems such as rebooting",
}

var systemActionsStartCmd = &cobra.Command{
	Use:  "start",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error { return runAction("start", args, false) },
}

var systemActionsStopCmd = &cobra.Command{
	Use:  "stop",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error { return runAction("stop", args, false) },
}

var systemActionsShutdownCmd = &cobra.Command{
	Use:  "shutdown",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error { return runAction("shutdown", args, false) },
}

var systemActionsRebootCmd = &cobra.Command{
	Use:  "reboot",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error { return runAction("reboot", args, false) },
}

var systemActionsResetCmd = &cobra.Command{
	Use:  "reset",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error { return runAction("reset", args, false) },
}

var systemActionsEmergencyPowerOffCmd = &cobra.Command{
	Use:  "emergencyPowerOff",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error { return runAction("emergencyPowerOff", args, false) },
}

var systemActionsDeactivateCmd = &cobra.Command{
	Use:  "deactivate",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error { return runAction("deactivate", args, false) },
}

var systemActionsActivateCmd = &cobra.Command{
	Use:  "activate",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error { return runAction("activate", args, false) },
}

var systemActionsAutoInstallCmd = &cobra.Command{
	Use:  "autoInstall",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error { return runAction("autoInstall", args, false) },
}

var systemActionsHypervisorFailedCmd = &cobra.Command{
	Use:  "hypervisorFailed",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error { return runAction("hypervisorFailed", args, false) },
}

var systemActionsStartMaintenanceDuration int32
var systemActionsStartMaintenanceCmd = &cobra.Command{
	Use:   "startMaintenance <system>",
	Short: "Mark the system as being in maintenance",
	Long: `Mark the system as being in maintenance

Marking a system as maintenance silences alerts on it.
`,
	Example: `Put a system in maintenance:
  lvl system actions startMaintenance web1
Put a system in maintenance for one hour:
  lvl system actions startMaintenance web1 --duration 60
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		system, err := Level27Client.SystemActionStartMaintenance(id, systemActionsStartMaintenanceDuration)
		if err != nil {
			return err
		}

		outputFormatTemplate(system, "templates/entities/system/actions/startMaintenance.tmpl")
		return nil
	},
}

var systemActionsStopMaintenanceCmd = &cobra.Command{
	Use:   "stopMaintenance <system>",
	Short: "Mark a system as no longer being in maintenance",
	Args:  cobra.ExactArgs(1),
	RunE:  func(cmd *cobra.Command, args []string) error { return runAction("stopMaintenance", args, true) },
}

func runAction(action string, args []string, templateResponse bool) error {
	id, err := resolveSystem(args[0])
	if err != nil {
		return err
	}

	system, err := Level27Client.SystemAction(id, action)
	if err != nil {
		return err
	}

	template := "templates/entities/system/action.tmpl"
	if templateResponse {
		template = fmt.Sprintf("templates/entities/system/actions/%s.tmpl", action)
	}
	outputFormatTemplate(system, template)
	return nil
}

// #endregion

type DescribeSystem struct {
	l27.System
	SshKeys                      []l27.SystemSshkey     `json:"sshKeys"`
	InstallSecurityUpdatesString string                 `json:"installSecurityUpdatesString"`
	HasNetworks                  []l27.SystemHasNetwork `json:"hasNetworks"`
	Volumes                      []l27.SystemVolume     `json:"volumes"`
	Checks                       []l27.SystemCheckGet   `json:"checks"`
}
