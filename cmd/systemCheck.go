package cmd

import (
	"fmt"

	"github.com/Jeffail/gabs/v2"
	"github.com/level27/l27-go"
	"github.com/level27/lvl/utils"
	"github.com/spf13/cobra"
)

func init() {
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
	flags := systemCheckAddCmd.Flags()
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
	systemCheckUpdateCmd.Flags().StringArrayVarP(&systemCheckUnsetParams, "unset", "u", systemCheckUnsetParams, "Unset an existing parameter on a check, restoring it to its default value")

	// #endregion

	// ------------------------------------ MONITORING ON SPECIFIC SYSTEM ----------------------------------------------
	// ---- MONITORING COMMAND
	systemCmd.AddCommand(systemMonitoringCmd)

	// ---- MONITORING ON
	systemMonitoringCmd.AddCommand(systemMonitoringOnCmd)
	// ---- MONITORING OFF
	systemMonitoringCmd.AddCommand(systemMonitoringOffCmd)
}

func resolveSystemCheck(systemID l27.IntID, arg string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.SystemCheckLookup(systemID, arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"system check",
		func(app l27.SystemCheckGet) string { return fmt.Sprintf("%d", app.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, err
}

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
			return fmt.Errorf("monitoring is currently disabled for system: [NAME:%v - ID: %v]. Use the 'monitoring' command to change monitoring status", system.Name, system.ID)
		}

		checks, err := Level27Client.SystemCheckGetList(id, optGetParameters)
		if err != nil {
			return err
		}

		// Creating readable output
		outputFormatTableFuncs(checks, []string{"ID", "CHECKTYPE", "STATUS", "LAST_STATUS_CHANGE", "INFORMATION"},
			[]interface{}{"ID", "CheckType", "Status", func(s l27.SystemCheckGet) string { return utils.FormatUnixTime(s.DtLastStatusChanged) }, "StatusInformation"})

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

		outputFormatTemplate(check, "templates/entities/systemCheck/add.tmpl")
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
		checkID, err := resolveSystemCheck(systemID, args[1])
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
		checkID, err := resolveSystemCheck(systemID, args[1])
		if err != nil {
			return err
		}

		if !optDeleteConfirmed {
			system, err := Level27Client.SystemGetSingle(systemID)
			if err != nil {
				return err
			}

			if !confirmPrompt(fmt.Sprintf("Delete system check %d on system %s (%d)?", checkID, system.Name, system.ID)) {
				return nil
			}
		}

		err = Level27Client.SystemCheckDelete(systemID, checkID)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/systemCheck/delete.tmpl")
		return nil
	},
}

// -------------- UPDATE SPECIFIC CHECK
var systemCheckUnsetParams []string
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
		checkID, err := resolveSystemCheck(systemID, args[1])
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
		for key, value := range currentData.CheckParameters.Map {
			// put each possible parrameter in array for later
			possibleParameters = append(possibleParameters, key)

			if !value.Default && !sliceContains(systemCheckUnsetParams, key) {
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
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/systemCheck/update.tmpl")
		return nil
	},
}

// #endregion

// ------------------------------------------------- MONITORING ON SPECIFIC SYSTEM ----------------------------------------------
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
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		_, err = Level27Client.SystemAction(systemID, "enable_monitoring")
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/system/monitoringOn.tmpl")
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
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		_, err = Level27Client.SystemAction(systemID, "disable_monitoring")
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/system/monitoringOff.tmpl")
		return nil
	},
}
