package cmd

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

func init() {
	// APP COMPONENT
	appCmd.AddCommand(appComponentCmd)

	// APP COMPONENT GET
	appComponentCmd.AddCommand(appComponentGetCmd)
	addCommonGetFlags(appComponentGetCmd)

	// APP COMPONENT CREATE
	appComponentCmd.AddCommand(appComponentCreateCmd)
	addWaitFlag(appComponentCreateCmd)
	appComponentCreateCmd.Flags().StringVarP(&appComponentCreateParamsFile, "params-file", "f", "", "JSON file to read params from. Pass '-' to read from stdin.")
	appComponentCreateCmd.Flags().StringVar(&appComponentCreateName, "name", "", "")
	appComponentCreateCmd.Flags().StringVar(&appComponentCreateType, "type", "", "")
	appComponentCreateCmd.Flags().StringVar(&appComponentCreateSystem, "system", "", "")
	appComponentCreateCmd.Flags().StringVar(&appComponentCreateSystemgroup, "systemgroup", "", "")
	appComponentCreateCmd.Flags().StringVar(&appComponentCreateLimitgroup, "limitgroup", "", "For Agency Hosting applications, which limit group the component will be added to.")
	appComponentCreateCmd.Flags().Int32Var(&appComponentCreateSystemprovider, "systemprovider", 0, "")
	appComponentCreateCmd.Flags().Int32Var(&appComponentCreateAttachment, "attachment", 0, "ID of the attachment to use with the appcomponent. Used for some components, such as solr config upload. Attachments may be managed with the 'lvl app component attachment' set of commands.")
	appComponentCreateParams = appComponentCreateCmd.Flags().StringArray("param", nil, "")
	appComponentCreateCmd.MarkFlagRequired("name")
	appComponentCreateCmd.MarkFlagRequired("type")

	// APP COMPONENT UPDATE
	appComponentCmd.AddCommand(appComponentUpdateCmd)
	settingsFileFlag(appComponentUpdateCmd)
	settingString(appComponentUpdateCmd, updateSettings, "name", "New name for this app component")
	appComponentUpdateParams = appComponentUpdateCmd.Flags().StringArray("param", nil, "")

	// APP COMPONENT DELETE
	appComponentCmd.AddCommand(AppComponentDeleteCmd)
	addDeleteConfirmFlag(AppComponentDeleteCmd)

	// APP COMPONENT CATEGORIES
	appComponentCmd.AddCommand(appComponentCategoryGetCmd)

	// APP COMPONENT TYPES
	appComponentCmd.AddCommand(appComponentTypeCmd)

	// APP COMPONENT PARAMETERS
	appComponentCmd.AddCommand(appComponentParametersCmd)
	appComponentParametersCmd.Flags().StringVarP(&appComponentType, "type", "t", "", "The type name to show its parameters.")
	appComponentParametersCmd.MarkFlagRequired("type")
}

// Resolve the ID of an app component based on user-provided name or ID.
func resolveAppComponent(appID l27.IntID, arg string) (l27.IntID, error) {
	// if arg already int, this is the ID
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.AppComponentLookup(appID, arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"component",
		func(comp l27.AppComponent) string { return fmt.Sprintf("%s (%d)", comp.Name, comp.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, nil
}

// APP COMPONENT CREATE
var appComponentCreateParamsFile string
var appComponentCreateName string
var appComponentCreateType string
var appComponentCreateSystem string
var appComponentCreateSystemgroup string
var appComponentCreateLimitgroup string
var appComponentCreateSystemprovider l27.IntID
var appComponentCreateAttachment l27.IntID
var appComponentCreateParams *[]string
var appComponentCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new appcomponent.",
	Example: "lvl app component create --name myComponentName --type docker",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if appComponentCreateSystem == "" && appComponentCreateSystemgroup == "" && appComponentCreateLimitgroup == "" {
			return errors.New("must specify either a system or a system group")
		}

		if appComponentCreateSystem != "" && appComponentCreateSystemgroup != "" {
			return errors.New("cannot specify both a system and a system group")
		}

		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		componentTypes, err := Level27Client.AppComponenttypesGet()
		if err != nil {
			return err
		}

		val, ok := componentTypes[appComponentCreateType]
		if !ok {
			return fmt.Errorf("unknown component type: %s", appComponentCreateType)
		}

		addSshKeyParameter(&val)

		paramsPassed, err := loadSettings(appComponentCreateParamsFile)
		if err != nil {
			return err
		}

		// Parse params from command line
		for _, param := range *appComponentCreateParams {
			split := strings.SplitN(param, "=", 2)
			if len(split) != 2 {
				return fmt.Errorf("expected key=value pair to --param: %s", param)
			}

			paramsPassed[split[0]], err = readArgFileSupported(split[1])
			if err != nil {
				return err
			}
		}

		create := map[string]interface{}{}
		create["name"] = appComponentCreateName
		create["appcomponenttype"] = appComponentCreateType

		if appComponentCreateLimitgroup != "" {
			create["limitGroup"] = appComponentCreateLimitgroup
		}

		if appComponentCreateSystem != "" {
			create["system"], err = resolveSystem(appComponentCreateSystem)
			if err != nil {
				return err
			}
		}

		if appComponentCreateSystemgroup != "" {
			create["systemgroup"], err = checkSingleIntID(appComponentCreateSystemgroup, "systemgroup")
			if err != nil {
				return err
			}
		}

		if appComponentCreateSystemprovider != 0 {
			create["systemprovider"] = appComponentCreateSystemprovider
		}

		if appComponentCreateAttachment != 0 {
			create["attachment"] = appComponentCreateAttachment
		}

		// Go over specified commands in app component types to validate and map data.

		paramNames := map[string]bool{}
		for _, param := range val.Servicetype.Parameters {
			paramName := param.Name
			paramNames[paramName] = true
			paramValue, hasValue := paramsPassed[paramName]
			if hasValue {
				if (param.Readonly || param.DisableEdit) && param.Type != "dependentSelect" {
					return fmt.Errorf("param cannot be changed: %s", paramName)
				}

				res, err := parseComponentParameter(param, paramValue)
				if err != nil {
					return err
				}

				create[paramName] = res
			} else if param.Required && param.DefaultValue == nil && !param.Received {
				return fmt.Errorf("required parameter not given: %s", paramName)
			}
		}

		// Check that there aren't any params given that don't exist.
		for k := range paramsPassed {
			if !paramNames[k] {
				return fmt.Errorf("unknown parameter given: %s", k)
			}
		}

		component, err := Level27Client.AppComponentCreate(appID, create)
		if err != nil {
			return err
		}

		if optWait {
			component, err = waitForStatus(
				func() (l27.AppComponent, error) { return Level27Client.AppComponentGetSingle(appID, component.ID) },
				func(ac l27.AppComponent) string { return ac.Status },
				"ok",
				[]string{"creating", "to_create"},
			)

			if err != nil {
				return fmt.Errorf("waiting on component status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(component, "templates/entities/appComponent/create.tmpl")
		return nil
	},
}

// APP COMPONENT UPDATE
var appComponentUpdateParams *[]string
var appComponentUpdateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a new appcomponent.",
	Example: "",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		settings, err := loadMergeSettings(updateSettingsFile, updateSettings)
		if err != nil {
			return err
		}

		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		appComponentID, err := resolveAppComponent(appID, args[1])
		if err != nil {
			return err
		}

		componentTypes, err := Level27Client.AppComponenttypesGet()
		if err != nil {
			return err
		}

		appComponent, err := Level27Client.AppComponentGetSingle(appID, appComponentID)
		if err != nil {
			return err
		}

		serviceType := componentTypes[appComponent.Appcomponenttype]

		addSshKeyParameter(&serviceType)

		parameterTypes := make(map[string]l27.AppComponentTypeParameter)
		for _, param := range serviceType.Servicetype.Parameters {
			parameterTypes[param.Name] = param
		}

		data := make(map[string]interface{})
		data["appcomponenttype"] = appComponent.Appcomponenttype
		data["name"] = appComponent.Name

		// Limit group implies Agency Hosting.
		isAgency := appComponent.LimitGroup != nil

		for k, v := range appComponent.Appcomponentparameters {
			param, ok := parameterTypes[k]
			if !ok {
				// Maybe should be a panic instead?
				return fmt.Errorf("API returned unknown parameter in component data: '%s'", k)
			}

			paramType := parameterTypes[k].Type

			if param.Readonly || param.DisableEdit || (isAgency && param.DisableEditAgency) || (!isAgency && param.DisableEditClassic) {
				continue
			}

			switch paramType {
			case "password-sha512", "password-plain", "password-sha1", "passowrd-sha256-scram":
				// Passwords are sent as "******". Skip them to avoid getting API errors.
				continue
			case "sshkey[]":
				// Need to map SSH keys -> IDs
				sshKeys := v.([]interface{})
				ids := make([]l27.IntID, len(sshKeys))
				for i, sshKey := range sshKeys {
					keyCast := sshKey.(map[string]interface{})
					ids[i] = l27.IntID(keyCast["id"].(float64))
				}
				v = ids
			default:
				// Just use value as-is
			}

			data[k] = v
		}

		data = mergeMaps(data, settings)

		// Parse params from command line
		for _, param := range *appComponentUpdateParams {
			split := strings.SplitN(param, "=", 2)
			if len(split) != 2 {
				return fmt.Errorf("expected key=value pair to --param: %s", param)
			}

			paramName := split[0]
			paramValue, err := readArgFileSupported(split[1])
			if err != nil {
				return err
			}

			paramType, ok := parameterTypes[paramName]
			if !ok {
				return fmt.Errorf("unknown parameter: %s", paramName)
			}

			res, err := parseComponentParameter(paramType, paramValue)
			if err != nil {
				return err
			}

			data[paramName] = res
		}

		err = Level27Client.AppComponentUpdate(appID, appComponentID, data)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/appComponent/update.tmpl")
		return nil
	},
}

func parseComponentParameter(param l27.AppComponentTypeParameter, paramValue interface{}) (interface{}, error) {
	// Convert parameters to the correct types in-JSON.
	var str string
	var ok bool
	if str, ok = paramValue.(string); !ok {
		// Value isn't a string. This means it must have come from a JSON input file or something (i.e. not command line arg)
		// So assume it's the correct type and let the API complain if it isn't.
		return paramValue, nil
	}

	switch param.Type {
	case "sshkey[]":
		keys := []l27.IntID{}
		for _, key := range strings.Split(str, ",") {
			// TODO: Resolve SSH key
			res, err := checkSingleIntID(key, "SSH key")
			if err != nil {
				return nil, err
			}

			keys = append(keys, res)
		}
		return keys, nil
	case "integer":
		return strconv.Atoi(str)
	case "boolean":
		return strings.EqualFold(str, "true"), nil
	case "array":
		found := false
		for _, possibleValue := range param.PossibleValues {
			if str == possibleValue {
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf(
				"parameter %s: value '%s' not in range of possible values: %s",
				param.Name,
				str,
				strings.Join(param.PossibleValues, ", "))
		}

		return str, nil

	default:
		// Pass as string
		return str, nil
	}
}

func addSshKeyParameter(appComponentType *l27.AppcomponenttypeServicetype) {
	// Older versions of the API used to expose SSH keys as a component parameter with the "sshkey[]" type.
	// This was changed to .SSHKeyPossible on the service type,
	// which would make the existing parameter parsing code far less elegant.
	// I've decided the best solution is to add an SSH key parameter back to the list locally.


	type keys struct {
		En string `json:"en"`
		Nl string `json:"nl"`
	}


	if appComponentType.Servicetype.SSHKeyPossible {
		appComponentType.Servicetype.Parameters = append(
			appComponentType.Servicetype.Parameters,
			l27.AppComponentTypeParameter{
				Name:           "sshkeys",
				DisplayName:  keys{
				En:	"SSH Keys",
				Nl:	"SSH Sleutel",
				},
				Description:   keys{En: "The SSH keys that can be used to log into the component", Nl: ""},
				Type:           "sshkey[]",
				DefaultValue:   nil,
				Readonly:       false,
				DisableEdit:    false,
				Required:       false,
				Received:       false,
				Category:       "credential",
				PossibleValues: nil,
			})
	}
}

// APP COMPONENT DELETE
var AppComponentDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete component from an app.",
	Example: "lvl app component delete MyAppName MyComponentName",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// search for appID based on appName
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		// search for component based on name
		appComponentID, err := resolveAppComponent(appID, args[1])
		if err != nil {
			return err
		}

		if !optDeleteConfirmed {
			app, err := Level27Client.App(appID)
			if err != nil {
				return err
			}

			appComponent, err := Level27Client.AppComponentGetSingle(appID, appComponentID)
			if err != nil {
				return err
			}

			if !confirmPrompt(fmt.Sprintf("Delete app component %s (%d) on app %s (%d)?", appComponent.Name, appComponent.ID, app.Name, app.ID)) {
				return nil
			}
		}

		err = Level27Client.AppComponentsDelete(appID, appComponentID)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/appComponent/delete.tmpl")
		return nil
	},
}

// APP COMPONENT CATEGORIES
var currentComponentCategories = []string{"web-apps", "databases", "extensions"}
var appComponentCategoryGetCmd = &cobra.Command{
	Use:     "categories",
	Short:   "shows a list of all current appcomponent categories.",
	Example: "lvl app component categories",
	Run: func(cmd *cobra.Command, args []string) {

		// type to convert string into category type
		var AppcomponentCategories struct {
			Data []l27.AppcomponentCategory
		}

		for _, category := range currentComponentCategories {
			cat := l27.AppcomponentCategory{Name: category}
			AppcomponentCategories.Data = append(AppcomponentCategories.Data, cat)
		}

		// display output in readable table
		outputFormatTable(AppcomponentCategories.Data, []string{"CATEGORY"}, []string{"Name"})
	},
}

// APP COMPONENT TYPES
var appComponentTypeCmd = &cobra.Command{
	Use:     "types",
	Short:   "Shows a list of all current componenttypes.",
	Example: "lvl app component types",
	RunE: func(cmd *cobra.Command, args []string) error {
		// get map of all types back from API (API doesnt give slice back in this case.)
		types, err := Level27Client.AppComponenttypesGet()
		if err != nil {
			return err
		}

		//create a type that contains an appcomponenttype name and category
		type typeInfo struct {
			Name     string
			Category string
		}

		//create slice of type typeInfo -> used to generate readable output for user
		allTypes := []typeInfo{}

		// loop over result and filter out the types Name and category into the right format.
		for key, value := range types {
			allTypes = append(allTypes, typeInfo{Name: key, Category: value.Servicetype.Category})
		}

		// sort slice based on category
		sort.Slice(allTypes, func(i, j int) bool {
			return allTypes[i].Category < allTypes[j].Category
		})

		// print result for user
		outputFormatTable(allTypes, []string{"NAME", "CATEGORY"}, []string{"Name", "Category"})
		return nil
	},
}

// APP COMPONENT PARAMETERS
var appComponentType string
var appComponentParametersCmd = &cobra.Command{
	Use:     "parameters",
	Short:   "Show list of all possible parameters with their default values of a specific componenttype.",
	Example: "lvl app component parameters -t python",
	Args:    cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		// get map of all types (and their parameters) back from API (API doesnt give slice back in this case.)
		types, err := Level27Client.AppComponenttypesGet()
		if err != nil {
			return err
		}

		// check if chosen componenttype is found
		componenttype, isTypeFound := types[appComponentType]

		if !isTypeFound {
			return fmt.Errorf("given componenttype: '%v' NOT found", appComponentType)
		}

		addSshKeyParameter(&componenttype)

		outputFormatTable(componenttype.Servicetype.Parameters,
			[]string{"NAME", "DESCRIPTION", "TYPE", "DEFAULT_VALUE", "REQUIRED"},
			[]string{"Name", "Description", "Type", "DefaultValue", "Required"})

		return nil
	},
}
