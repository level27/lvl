package cmd

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

func init() {

	//-------------------------------------  SYSTEMS/COOKBOOKS TOPLEVEL (get/post) --------------------------------------
	// #region SYSTEMS/COOKBOOKS TOPLEVEL (get/post)

	// adding cookbook subcommand to system command
	systemCmd.AddCommand(systemCookbookCmd)

	// ---- GET cookbooks
	systemCookbookCmd.AddCommand(systemCookbookGetCmd)

	// ---- ADD cookbook (to system)
	systemCookbookCmd.AddCommand(systemCookbookAddCmd)
	addWaitFlag(systemCookbookAddCmd)
	flags := systemCookbookAddCmd.Flags()
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
	addWaitFlag(systemCookbookDeleteCmd)

	// --- UPDATE
	systemCookbookCmd.AddCommand(systemCookbookUpdateCmd)
	addWaitFlag(systemCookbookUpdateCmd)
	systemCookbookUpdateCmd.Flags().StringArrayVarP(&systemDynamicParams, "parameters", "p", systemDynamicParams, "Add custom parameters for a check. Usage -> SINGLE PAR: [ -p waf=true ], MULTIPLE PAR: [ -p waf=true -p timeout=200 ], MULTIPLE VALUES: [ -p versions=''7, 5.4'']")
	systemCookbookUpdateCmd.MarkFlagRequired("parameters")
	// #endregion
}

func resolveSystemCookbook(systemID l27.IntID, arg string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	cookbook, err := Level27Client.SystemCookbookLookup(systemID, arg)
	if err != nil {
		return 0, err
	}

	if cookbook == nil {
		// Try system settings instead.
		cookbook, err = Level27Client.SystemSettingsLookup(systemID, arg)
		if err != nil {
			return 0, err
		}

		if cookbook == nil {
			return 0, fmt.Errorf("system (%d) does not have a cookbook of type '%s'", systemID, arg)
		}
	}

	return cookbook.ID, nil
}

// ------------------------------------------------- SYSTEM/COOKBOOKS TOPLEVEL (GET error / CREATE) ----------------------------------
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

		cookbooks, err := Level27Client.SystemCookbookGetList(id, l27.CommonGetParams{})
		if err != nil {
			return err
		}

		settings, err := Level27Client.SystemSettingsGetList(id, l27.CommonGetParams{})
		if err != nil {
			return err
		}

		cookbooks = append(cookbooks, settings...)

		outputFormatTable(cookbooks, []string{"ID", "COOKBOOKTYPE", "STATUS"}, []string{"ID", "CookbookType", "Status"})
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
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		// get information about the current chosen system [systemID]
		currentSystem, err := Level27Client.SystemGetSingle(systemID)
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

		cookbook, err := Level27Client.SystemCookbookAdd(systemID, &cookbookRequest)
		if err != nil {
			return err
		}

		//apply changes to cookbooks
		err = Level27Client.SystemCookbookChangesApply(systemID)
		if err != nil {
			return err
		}

		if optWait {
			cookbook, err = waitForStatus(
				func() (l27.Cookbook, error) { return Level27Client.SystemCookbookDescribe(systemID, cookbook.ID) },
				func(s l27.Cookbook) string { return s.Status },
				"ok",
				[]string{"updating"},
			)

			if err != nil {
				return fmt.Errorf("waiting on cookbook status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(cookbook, "templates/entities/systemCookbook/add.tmpl")

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
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		cookbookID, err := resolveSystemCookbook(systemID, args[1])
		if err != nil {
			return err
		}

		result, err := Level27Client.SystemCookbookDescribe(systemID, cookbookID)
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
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		cookbookID, err := resolveSystemCookbook(systemID, args[1])
		if err != nil {
			return err
		}

		if !optDeleteConfirmed {
			cookbook, err := Level27Client.SystemCookbookDescribe(systemID, cookbookID)
			if err != nil {
				return err
			}

			if !confirmPrompt(fmt.Sprintf("Delete system cookbook %s (%d) on system %s (%d)?", cookbook.CookbookType, cookbook.ID, cookbook.System.Name, cookbook.System.ID)) {
				return nil
			}
		}

		err = Level27Client.SystemCookbookDelete(systemID, cookbookID)
		if err != nil {
			return err
		}

		//apply changes
		err = Level27Client.SystemCookbookChangesApply(systemID)
		if err != nil {
			return err
		}

		if optWait {
			err = waitForDelete(
				func() (l27.Cookbook, error) { return Level27Client.SystemCookbookDescribe(systemID, cookbookID) },
				func(s l27.Cookbook) string { return s.Status },
				[]string{"deleting"},
			)

			if err != nil {
				return fmt.Errorf("waiting on system cookbook status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(nil, "templates/entities/systemCookbook/delete.tmpl")

		return err
	},
}

// ---------------- UPDATE
var systemCookbookUpdateCmd = &cobra.Command{
	Use:     "update [systemID] [cookbookID]",
	Short:   "update existing cookbook from a system",
	Example: "lvl system cookbooks update [systemID] [cookbookID] {-p}.\nSINGLE PARAMETER:		-p waf=true  \nMULTIPLE PARAMETERS:		-p waf=true -p timeout=200  \nMULTIPLE VALUES:		-p versions=''7, 5.4'' OR -p versions=7,5.4 (seperated by comma)",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// check for valid system id
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		cookbookID, err := resolveSystemCookbook(systemID, args[1])
		if err != nil {
			return err
		}

		// get current data from the current installed cookbooktype
		currentCookbookData, err := Level27Client.SystemCookbookDescribe(systemID, cookbookID)
		if err != nil {
			return err
		}

		// get current data from the chosen system
		currentSystemData, err := Level27Client.SystemGetSingle(systemID)
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
		err = checkForValidCookbookParameter(customParameterDict, cookbookData, currentSystem, &baseRequestData)
		if err != nil {
			return err
		}

		err = Level27Client.SystemCookbookUpdate(systemID, cookbookID, &baseRequestData)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/systemCookbook/update.tmpl")

		// aplly changes to cookbooks
		err = Level27Client.SystemCookbookChangesApply(systemID)
		if err != nil {
			return err
		}

		if optWait {
			_, err = waitForStatus(
				func() (l27.Cookbook, error) { return Level27Client.SystemCookbookDescribe(systemID, cookbookID) },
				func(s l27.Cookbook) string { return s.Status },
				"ok",
				[]string{"updating"},
			)

			if err != nil {
				return fmt.Errorf("waiting on cookbook status failed: %s", err.Error())
			}
		}

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
