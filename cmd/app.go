package cmd

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"bitbucket.org/level27/lvl/types"
	"bitbucket.org/level27/lvl/utils"
	"github.com/spf13/cobra"
)

func init() {
	// ---- MAIN COMMAND APP
	RootCmd.AddCommand(appCmd)

	//------------------------------------------------- APP (GET / CREATE / DELETE / UPDATE / DESCRIBE)-------------------------------------------------
	// #region APP MAIN COMMANDS (GET / CREATE / UPDATE / DELETE / DESCRIBE)

	// ---- GET
	appCmd.AddCommand(appGetCmd)
	addCommonGetFlags(appGetCmd)

	// ---- CREATE
	appCmd.AddCommand(appCreateCmd)
	// flags used for creating app
	flags := appCreateCmd.Flags()
	flags.StringVarP(&appCreateName, "name", "n", "", "Name of the app.")
	flags.StringVarP(&appCreateOrg, "organisation", "", "", "The name of the organisation/owner of the app.")
	flags.IntSliceVar(&appCreateTeams, "autoTeams", appCreateTeams, "A csv list of team ID's.")
	flags.StringVar(&appCreateExtInfo, "externalInfo", "", "ExternalInfo (required when billableItemInfo entities for an organisation exist in DB.)")
	appCreateCmd.MarkFlagRequired("name")
	appCreateCmd.MarkFlagRequired("organisation")

	// ---- DELETE APP
	appCmd.AddCommand(appDeleteCmd)
	//flag to skip confirmation when deleting an app
	appDeleteCmd.Flags().BoolVarP(&isAppDeleteConfirmed, "yes", "y", false, "Set this flag to skip confirmation when deleting an app")

	// ---- UPDATE APP
	appCmd.AddCommand(appUpdateCmd)
	// flags needed for update command
	flags = appUpdateCmd.Flags()
	flags.StringVarP(&appUpdateName, "name", "n", "", "Name of the app.")
	flags.StringVarP(&appUpdateOrg, "organisation", "", "", "The name of the organisation/owner of the app.")
	flags.StringSliceVar(&appUpdateTeams, "autoTeams", appUpdateTeams, "A csv list of team ID's.")

	// ---- DESCRIBE APP
	appCmd.AddCommand(AppDescribeCmd)
	// #endregion

	//-------------------------------------------------  APP ACTIONS (ACTIVATE / DEACTIVATE) -------------------------------------------------
	// #region  APP ACTIONS (ACTIVATE / DEACTIVATE)

	// ACTION COMMAND
	appCmd.AddCommand(AppActionCmd)

	// ACTIVATE APP
	AppActionCmd.AddCommand(AppActionActivateCmd)

	// DEACTIVATE APP
	AppActionCmd.AddCommand(AppActionDeactivateCmd)
	// #endregion

	//------------------------------------------------- APP COMPONENTS (CREATE / GET / UPDATE / DELETE / DESCRIBE)-------------------------------------------------
	// ----	COMPONENT COMMAND
	appCmd.AddCommand(appComponentCmd)

	// ---- GET COMPONENTS
	appComponentCmd.AddCommand(appComponentGetCmd)
	addCommonGetFlags(appComponentGetCmd)

	// ---- CREATE COMPONENT
	appComponentCmd.AddCommand(appComponentCreateCmd)

	// ---- DELETE COMPONENTS
	appComponentCmd.AddCommand(AppComponentDeleteCmd)
	//flag to skip confirmation when deleting an appcomponent
	AppComponentDeleteCmd.Flags().BoolVarP(&isComponentDeleteConfirmed, "yes", "y", false, "Set this flag to skip confirmation when deleting an app")

	//------------------------------------------------- APP COMPONENTS HELPERS (CATEGORY/ COMPONENTTYPES/ PARAMETERS )-------------------------------------------------
	// ---- GET COMPONENT CATEGORIES
	appComponentCmd.AddCommand(appComponentCategoryGetCmd)

	// ---- GET COMPONENTTYPES
	appComponentCmd.AddCommand(appComponentTypeCmd)

	// ---- GET COMPONENTTYPE PARAMATERS
	appComponentCmd.AddCommand(appComponentParametersCmd)
	// flags needed to get parameters of a type
	appComponentParametersCmd.Flags().StringVarP(&appComponentType, "type", "t", "", "The type name to show its parameters.")
	appComponentParametersCmd.MarkFlagRequired("type")

	//-------------------------------------------------  APP ACCESS -------------------------------------------------
	addAccessCmds(appCmd, "apps", resolveApp)

	// ----------- APP SSL CERTIFICATE COMMANDS

	// ---- APP SSL
	appCmd.AddCommand(appSslCmd)

	// ---- APP SSL GET
	appSslCmd.AddCommand(appSslGetCmd)
	addCommonGetFlags(appSslCmd)

	// ---- APP SSL DESCRIBE
	appSslCmd.AddCommand(appSslDescribeCmd)

	// ---- APP SSL CREATE
	appSslCmd.AddCommand(appSslCreateCmd)
	appSslCreateCmd.Flags().StringVarP(&appSslCreateName, "name", "n", "", "Name of this SSL certificate")
	appSslCreateCmd.Flags().StringVarP(&appSslCreateSslType, "type", "t", "", "Type of SSL certificate to use. Options are: letsencrypt, xolphin, own")
	appSslCreateCmd.Flags().StringVar(&appSslCreateAutoSslCertificateUrls, "auto-urls", "", "URL or CSV list of URLs (required for Let's Encrypt)")
	appSslCreateCmd.Flags().StringVar(&appSslCreateSslKey, "ssl-key", "", "SSL key for own certificate. Can be read from a file by specifying @filename.")
	appSslCreateCmd.Flags().StringVar(&appSslCreateSslCabundle, "ssl-cabundle", "", "SSL CA bundle for own certificate. Can be read from a file by specifying @filename.")
	appSslCreateCmd.Flags().StringVar(&appSslCreateSslCrt, "ssl-crt", "", "SSL CRT for own certificate. Can be read from a file by specifying @filename.")
	appSslCreateCmd.Flags().BoolVar(&appSslCreateAutoUrlLink, "auto-link", false, "After creation, automatically link to any URLs without existing certificate")
	appSslCreateCmd.Flags().BoolVar(&appSslCreateSslForce, "ssl-force", false, "Force SSL")
	appSslCreateCmd.MarkFlagRequired("name")
	appSslCreateCmd.MarkFlagRequired("type")

	// ---- APP SSL DELETE
	appSslCmd.AddCommand(appSslDeleteCmd)
	appSslDeleteCmd.Flags().BoolVar(&appSslDeleteForce, "force", false, "Do not ask for confirmation to delete the SSL certificate")

	// ---- APP SSL UPDATE
	appSslCmd.AddCommand(appSslUpdateCmd)
	settingsFileFlag(appSslUpdateCmd)
	settingString(appSslUpdateCmd, updateSettings, "name", "New name for the SSL certificate")

	// ---- APP SSL FIX
	appSslCmd.AddCommand(appSslFixCmd)

	// ---- APP SSL ACTION
	appSslCmd.AddCommand(appSslActionCmd)

	// ---- APP SSL KEY
	appSslCmd.AddCommand(appSslKeyCmd)

	// ---- Actions (Retry and ValidateChallenge)
	appSslActionCmd.AddCommand(appSslActionRetryCmd)
	appSslActionCmd.AddCommand(appSslActionValidateChallengeCmd)

	//-------------------------------------------------  APP ACCESS -------------------------------------------------
	// ---- ACCESS COMMANDS
	addAccessCmds(appCmd, "apps", resolveApp)

	//-------------------------------------------------  APP COMPONENT RESTORE (GET / DESCRIBE / CREATE / UPDATE / DELETE / DOWNLOAD) -------------------------------------------------

	// ---- RESTORE COMMAND
	appComponentCmd.AddCommand(appComponentRestoreCmd)

	// ---- GET LIST OF RESTORES
	appComponentRestoreCmd.AddCommand(appComponentRestoreGetCmd)

	// ---- CREATE RESTORE FOR APPCOMPONENT
	appComponentRestoreCmd.AddCommand(appComponentRestoreCreateCmd)

	// ---- DELETE RESTORE
	appComponentRestoreCmd.AddCommand(appRestoreDeleteCmd)
	//flag to skip confirmation when deleting a restore
	appRestoreDeleteCmd.Flags().BoolVarP(&isAppRestoreDeleteConfirmed, "yes", "y", false, "Set this flag to skip confirmation when deleting a check")

	// ---- DOWNLOAD RESTORE FILE
	appComponentRestoreCmd.AddCommand(appComponentRestoreDownloadCmd)
	// flags needed for downloading the restore
	appComponentRestoreDownloadCmd.Flags().StringVarP(&appComponentRestoreDownloadName, "filename", "f", "", "The name of the downloaded file.")
	//-------------------------------------------------  APP COMPONENT BACKUP (GET) -------------------------------------------------
	// ---- BACKUP COMMAND
	appComponentCmd.AddCommand(appComponentBackupsCmd)
	// ---- GET LIST OF BACKUPS
	appComponentBackupsCmd.AddCommand(appComponentBackupsGetCmd)

	//-------------------------------------------------  APP MIGRATIONS (GET / DESCRIBE / CREATE / UPDATE) -------------------------------------------------
	// ---- MIGRATIONS COMMAND
	appCmd.AddCommand(appMigrationsCmd)

	// ---- GET LIST OF MIGRATIONS
	appMigrationsCmd.AddCommand(appMigrationsGetCmd)
}
//------------------------------------------------- APP HELPER FUNCTIONS -------------------------------------------------
// GET AN APPID BASED ON THE NAME
func resolveApp(arg string) int {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id
	}

	return resolveShared(
		Level27Client.AppLookup(arg),
		arg,
		"app",
		func(app types.App) string { return fmt.Sprintf("%s (%d)", app.Name, app.ID) }).ID
}

// GET SSL CERTIFICATE ID BASED ON ID
func resolveAppSslCertificate(appID int, arg string) int {
	// if arg already int, this is the ID
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id
	}

	return resolveShared(
		Level27Client.AppSslCertificatesLookup(appID, arg),
		arg,
		"app SSL certificate",
		func(app types.AppSslCertificate) string { return fmt.Sprintf("%s (%d)", app.Name, app.ID) }).ID
}


// GET AN APPCOMPONENT ID BASED ON THE NAME
func resolveAppComponent(appId int, arg string) int {
	// if arg already int, this is the ID
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id
	}

	return int(resolveShared(
		Level27Client.AppComponentLookup(appId, arg),
		arg,
		"component",
		func(comp types.AppComponent) string { return fmt.Sprintf("%s (%d)", comp.Name, comp.ID) }).ID)
}

// MAIN COMMAND APPS
var appCmd = &cobra.Command{
	Use:     "app",
	Short:   "Commands to manage apps",
	Example: "lvl app get -f searchThisApp\nlvl app action activate",
}

//------------------------------------------------- APP (GET / CREATE  / UPDATE / DELETE / DESCRIBE)-------------------------------------------------
// #region APP MAIN SUBCOMMANDS (GET / CREATE / UPDATE / DELETE / DESCRIBE)

// ---- GET apps
var appGetCmd = &cobra.Command{
	Use:     "get",
	Short:   "Shows a list of all available apps.",
	Example: "lvl app get",
	Args:    cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ids, err := convertStringsToIds(args)
		if err != nil {
			log.Fatalln("Invalid app ID")
		}

		outputFormatTable(
			getApps(ids),
			[]string{"ID", "NAME", "STATUS"},
			[]string{"ID", "Name", "Status"})
	},
}

func getApps(ids []int) []types.App {
	c := Level27Client
	if len(ids) == 0 {
		return c.Apps(optGetParameters)
	} else {
		apps := make([]types.App, len(ids))
		for idx, id := range ids {
			apps[idx] = c.App(id)
		}
		return apps
	}
}

// ---- CREATE NEW APP
var appCreateName, appCreateOrg, appCreateExtInfo string
var appCreateTeams []int
var appCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new app.",
	Example: "lvl app create -n myNewApp --organisation level27",
	Run: func(cmd *cobra.Command, args []string) {
		// check if name is valid.
		if appCreateName == "" {
			log.Fatalln("app name cannot be empty.")
		}

		// fill in all the props needed for the post request
		organisation := resolveOrganisation(appCreateOrg)
		request := types.AppPostRequest{
			Name:         appCreateName,
			Organisation: organisation,
			AutoTeams:    appCreateTeams,
			ExternalInfo: appCreateExtInfo,
		}

		// when succesfully creating app. app will be returned
		app := Level27Client.AppCreate(request)
		log.Printf("Succesfully created app! [name: '%v' - ID: '%v']", app.Name, app.ID)

	},
}

// ---- DELETE AN APP
var isAppDeleteConfirmed bool
var appDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete an app",
	Example: "lvl app delete NameOfMyApp",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// try to find appId based on name
		appId := resolveApp(args[0])

		Level27Client.AppDelete(appId, isAppDeleteConfirmed)
	},
}

// ---- UPDATE AN APP
var appUpdateName, appUpdateOrg string
var appUpdateTeams []string
var appUpdateCmd = &cobra.Command{
	Use:     "update [appID]",
	Short:   "Update an app.",
	Example: "lvl app update 2067 --name myUpdatedName",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//check if appId is valid
		appId := resolveApp(args[0])

		//get the current data from the app. if not changed its needed for put request
		currentData := Level27Client.App(appId)

		var currentTeamIds []string
		for _, team := range currentData.Teams {
			currentTeamIds = append(currentTeamIds, strconv.Itoa(team.ID))
		}
		// fill in request with the current data.
		request := types.AppPutRequest{
			Name:         currentData.Name,
			Organisation: currentData.Organisation.ID,
			AutoTeams:    currentTeamIds,
		}

		//when flags have been set. we need the currentdata to be updated.
		if cmd.Flag("name").Changed {
			request.Name = appUpdateName
		}

		if cmd.Flag("organisation").Changed {
			organisationID := resolveOrganisation(appUpdateOrg)
			request.Organisation = organisationID
		}

		if cmd.Flag("autoTeams").Changed {
			request.AutoTeams = appUpdateTeams
		}

		Level27Client.AppUpdate(appId, request)

	},
}

// ---- DESCRIBE APP
var AppDescribeCmd = &cobra.Command{
	Use:     "describe",
	Short:   "Get detailed info about an app.",
	Example: "lvl app describe 2077",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//check for valid appId
		appId := resolveApp(args[0])
		// get all data from app by appId
		app := Level27Client.App(appId)
		outputFormatTemplate(app, "templates/app.tmpl")
	},
}

// #endregion

//------------------------------------------------- APP ACTIONS (ACTIVATE / DEACTIVATE)-------------------------------------------------
// #region APP ACTIONS (ACTIVATE / DEACTIVATE)

// ---- ACTION COMMAND
var AppActionCmd = &cobra.Command{
	Use:     "action",
	Short:   "commands to run specific actions on an app",
	Example: "lvl app action [subcommand]",
}

// ---- ACTIVATE APP
var AppActionActivateCmd = &cobra.Command{
	Use:     "activate",
	Short:   "Activate an app",
	Example: "lvl app action activate 2077",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// check for valid appId
		appId := resolveApp(args[0])

		Level27Client.AppAction(appId, "activate")
	},
}

// ---- DEACTIVATE APP
var AppActionDeactivateCmd = &cobra.Command{
	Use:     "deactivate",
	Short:   "Deactivate an app",
	Example: "lvl app action deactivate 2077",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// check for valid appId
		appId := resolveApp(args[0])

		Level27Client.AppAction(appId, "deactivate")
	},
}

// APP SSL CERTIFICATES

// APP SSL
var appSslCmd = &cobra.Command{
	Use:     "ssl",
	Short:   "Commands for managing SSL certificates on apps",
	Example: "lvl app ssl get forum\nlvl app ssl describe forum forum.example.com",

	Aliases: []string{"sslcert"},
}

// APP SSL GET
var appSslGetType string
var appSslGetStatus string
var appSslGetCmd = &cobra.Command{
	Use:     "get [app]",
	Short:   "Get a list of SSL certificates for an app",
	Example: "lvl app ssl get forum\nlvl app ssl get forum -f admin",

	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appID := resolveApp(args[0])

		certs := resolveGets(
			// First arg is app ID.
			args[1:],
			func(name string) []types.AppSslCertificate {
				return Level27Client.AppSslCertificatesLookup(appID, name)
			},
			func(certID int) types.AppSslCertificate {
				return Level27Client.AppSslCertificatesGetSingle(appID, certID)
			},
			func(get types.CommonGetParams) []types.AppSslCertificate {
				return Level27Client.AppSslCertificatesGetList(appID, appSslGetType, appSslGetStatus, get)
			},
		)

		outputFormatTableFuncs(
			certs,
			[]string{"ID", "Name", "Type", "Status", "SSL Status", "Expiry Date"},
			[]interface{}{"ID", "Name", "SslType", "Status", "SslStatus", "DtExpires", func(c types.AppSslCertificate) string { return utils.FormatUnixTime(c.DtExpires) }})
	},
}

// APP SSL DESCRIBE

var appSslDescribeCmd = &cobra.Command{
	Use:     "describe [app] [SSL cert]",
	Short:   "Get detailed information of an SSL certificate",
	Example: "lvl app ssl describe forum forum.example.com",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		appID := resolveApp(args[0])
		certID := resolveAppSslCertificate(appID, args[1])

		cert := Level27Client.AppSslCertificatesGetSingle(appID, certID)

		outputFormatTemplate(cert, "templates/appSslCertificate.tmpl")
	},
}

// APP SSL CREATE
var appSslCreateName string
var appSslCreateSslType string
var appSslCreateAutoSslCertificateUrls string
var appSslCreateSslKey string
var appSslCreateSslCrt string
var appSslCreateSslCabundle string
var appSslCreateAutoUrlLink bool
var appSslCreateSslForce bool

var appSslCreateCmd = &cobra.Command{
	Use:     "create [app]",
	Short:   "Create a new SSL certificate on an app",
	Example: "lvl app ssl create forum --name forum.example.com --auto-urls forum.example.com --auto-link --type letsencrypt\nlvl app ssl create forum --name forum.example.com --type own --ssl-cabundle '@cert.ca-bundle' --ssl-key '@key.pem' --ssl-crt '@cert.crt'",

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appID := resolveApp(args[0])

		create := types.AppSslCertificateCreate{
			Name:                   appSslCreateName,
			SslType:                appSslCreateSslType,
			AutoSslCertificateUrls: appSslCreateAutoSslCertificateUrls,
			SslForce:               appSslCreateSslForce,
			AutoUrlLink:            appSslCreateAutoUrlLink,
		}

		var certificate types.AppSslCertificate

		switch appSslCreateSslType {
		case "own":
			createOwn := types.AppSslCertificateCreateOwn{
				AppSslCertificateCreate: create,
				SslKey:                  readArgFileSupported(appSslCreateSslKey),
				SslCrt:                  readArgFileSupported(appSslCreateSslCrt),
				SslCabundle:             readArgFileSupported(appSslCreateSslCabundle),
			}
			certificate = Level27Client.AppSslCertificatesCreateOwn(appID, createOwn)

		case "letsencrypt", "xolphin":
			certificate = Level27Client.AppSslCertificatesCreate(appID, create)

		default:
			cobra.CheckErr(fmt.Sprintf("Invalid SSL type: %s", appSslCreateSslType))
		}

		fmt.Printf("sslCertificate created: [name: '%s' - ID: '%d'].", certificate.Name, certificate.ID)
	},
}

// APP SSL DELETE

var appSslDeleteForce bool
var appSslDeleteCmd = &cobra.Command{
	Use:     "delete [app] [SSL cert]",
	Short:   "Delete an SSL certificate from an app",
	Example: "lvl app ssl delete forum forum.example.com",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		appID := resolveApp(args[0])
		certID := resolveAppSslCertificate(appID, args[1])

		if !appSslDeleteForce {
			app := Level27Client.App(appID)
			cert := Level27Client.AppSslCertificatesGetSingle(appID, certID)
			if !confirmPrompt(fmt.Sprintf("Delete SSL certificate %s (%d) on app %s (%d)?", cert.Name, certID, app.Name, appID)) {
				return
			}
		}

		Level27Client.AppSslCertificatesDelete(appID, certID)
	},
}

// APP SSL UPDATE
var appSslUpdateCmd = &cobra.Command{
	Use: "update [app] [SSL cert]",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		settings := loadMergeSettings(updateSettingsFile, updateSettings)

		appID := resolveApp(args[0])
		certID := resolveAppSslCertificate(appID, args[1])

		cert := Level27Client.AppSslCertificatesGetSingle(appID, certID)
		put := types.AppSslCertificatePut{
			Name:    cert.Name,
			SslType: cert.SslType,
		}

		data := utils.RoundTripJson(put).(map[string]interface{})
		data = mergeMaps(data, settings)

		Level27Client.AppSslCertificatesUpdate(appID, certID, data)
	},
}

// APP SSL FIX
var appSslFixCmd = &cobra.Command{
	Use:     "fix [app] [SSL cert]",
	Short:   "Fix an invalid SSL certificate",
	Example: "lvl app ssl fix forum forum.example.com",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		appID := resolveApp(args[0])
		certID := resolveAppSslCertificate(appID, args[1])

		Level27Client.AppSslCertificatesFix(appID, certID)
	},
}

// APP SSL ACTION

var appSslActionCmd = &cobra.Command{
	Use: "action",
}

var appSslActionRetryCmd = &cobra.Command{
	Use: "retry [app] [SSL cert]",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		appID := resolveApp(args[0])
		certID := resolveAppSslCertificate(appID, args[1])

		Level27Client.AppSslCertificatesActions(appID, certID, "retry")
	},
}

var appSslActionValidateChallengeCmd = &cobra.Command{
	Use: "validateChallenge [app] [SSL cert]",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		appID := resolveApp(args[0])
		certID := resolveAppSslCertificate(appID, args[1])

		Level27Client.AppSslCertificatesActions(appID, certID, "validateChallenge")
	},
}

// APP SSL KEY
var appSslKeyCmd = &cobra.Command{
	Use:     "key",
	Short:   "Return a private key for type 'own' sslCertificate.",
	Example: "lvl app ssl key MyAppName",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		appId := resolveApp(args[0])
		certID := resolveAppSslCertificate(appId, args[1])

		key := Level27Client.AppSslCertificatesKey(appId, certID)

		fmt.Print(key.SslKey)
	},
}

// #endregion

//------------------------------------------------- APP COMPONENTS (CREATE / GET / UPDATE / DELETE / DESCRIBE)-------------------------------------------------
// ---- COMPONENT COMMAND
var appComponentCmd = &cobra.Command{
	Use:     "component",
	Short:   "Commands for managing appcomponents.",
	Example: "lvl app component get",
}

// ---- GET COMPONENTS
var appComponentGetCmd = &cobra.Command{
	Use:     "get [App]",
	Short:   "Show list of all available components on an app.",
	Example: "lvl app component get MyAppName",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		//search for appId based on Appname
		appId := resolveApp(args[0])
		ids, err := convertStringsToIds(args[1:])
		if err != nil {
			log.Fatalln("Invalid component ID")
		}

		outputFormatTable(
			getComponents(appId, ids),
			[]string{"ID", "NAME", "STATUS"},
			[]string{"ID", "Name", "Status"})
	},
}

func getComponents(appId int, ids []int) []types.AppComponent {
	c := Level27Client
	if len(ids) == 0 {
		return c.AppComponentsGet(appId, optGetParameters)
	} else {
		components := make([]types.AppComponent, len(ids))
		for idx, id := range ids {
			components[idx] = c.AppComponentGetSingle(appId, id)
		}
		return components
	}
}

// ---- CREATE COMPONENT
var appComponentCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new appcomponent.",
	Example: "lvl app component create -n myComponentName -c docker -ctype mysql",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// ---- DELETE COMPONENT
var isComponentDeleteConfirmed bool
var AppComponentDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete component from an app.",
	Example: "lvl app component delete MyAppName MyComponentName",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// search for appId based on appName
		appId := resolveApp(args[0])
		// search for component based on name
		appComponentId := resolveAppComponent(appId, args[1])

		Level27Client.AppComponentsDelete(appId, appComponentId, isComponentDeleteConfirmed)
	},
}

//------------------------------------------------- APP COMPONENTS HELPERS (CATEGORY/ TYPES/ PARAMETERS )-------------------------------------------------

// current possible current categories for appcomponents
var currentComponentCategories = []string{"web-apps", "databases", "extensions"}

// #region APP COMPONENTS HELPERS (CATEGORY/ TYPES/ PARAMETERS )

// ---- (CATEGORY) GET COMPONENT CATEGORIES
var appComponentCategoryGetCmd = &cobra.Command{
	Use:     "categories",
	Short:   "shows a list of all current appcomponent categories.",
	Example: "lvl app component categories",
	Run: func(cmd *cobra.Command, args []string) {

		// type to convert string into category type
		var AppcomponentCategories struct {
			Data []types.AppcomponentCategory
		}

		for _, category := range currentComponentCategories {
			cat := types.AppcomponentCategory{Name: category}
			AppcomponentCategories.Data = append(AppcomponentCategories.Data, cat)
		}
		// display output in readable table
		outputFormatTable(AppcomponentCategories.Data, []string{"CATEGORY"}, []string{"Name"})
	},
}

// ---- (TYPES) GET LIST OF APPCOMPONENT TYPES
var appComponentTypeCmd = &cobra.Command{
	Use:     "types",
	Short:   "Shows a list of all current componenttypes.",
	Example: "lvl app component types",
	Run: func(cmd *cobra.Command, args []string) {

		// get map of all types back from API (API doesnt give slice back in this case.)
		types := Level27Client.AppComponenttypesGet()

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

	},
}

// ---- (PARAMETERS) GET LIST OF PARAMETERS FOR A SPECIFIC APPCOMPONENT TYPE
var appComponentType string
var appComponentParametersCmd = &cobra.Command{
	Use:     "parameters",
	Short:   "Show list of all possible parameters with their default values of a specific componenttype.",
	Example: "lvl app component parameters -t python",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		// get map of all types (and their parameters) back from API (API doesnt give slice back in this case.)
		types := Level27Client.AppComponenttypesGet()

		// check if chosen componenttype is found
		componenttype, isTypeFound := types[appComponentType]

		if isTypeFound {
			outputFormatTable(componenttype.Servicetype.Parameters,
				[]string{"NAME", "DESCRIPTION", "TYPE", "DEFAULT_VALUE", "REQUIRED"},
				[]string{"Name", "Description", "Type", "DefaultValue", "Required"})
		} else {
			log.Fatalf("Given componenttype: '%v' NOT found!", appComponentType)
		}

	},
}

// #endregion

//-------------------------------------------------  APP SSL CERTIFICATES (GET / CREATE/ DELETE/ FIX) -------------------------------------------------
// ---- SSL COMMAND
var appCertificateCmd = &cobra.Command{
	Use:     "ssl",
	Short:   "Commands for managing ssl certificates.",
	Example: "lvl app ssl get",
}

// #region APP SSL CERTIFICATES (GET / CREATE/ DELETE/ FIX)

// ---- GET LIST OF SSL CERTIFICATES
var appCertificateGetCmd = &cobra.Command{
	Use:     "get",
	Short:   "Show list of all available ssl certificates on an app.",
	Example: "lvl app ssl get",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Search appId based on name
		appId := resolveApp(args[0])

		certificates := Level27Client.AppCertificateGet(appId)
		// Display output in readable table
		outputFormatTableFuncs(certificates, []string{"ID", "NAME", "TYPE", "STATUS", "EXPIRING DATE"},
			[]interface{}{"ID", "Name", "SslType", "Status", func(s types.SslCertificate) string { return utils.FormatUnixTime(s.DtExpires) }})

	},
}

// ---- ADD CERTIFICATE
// vars to hold properties set by user.
var appAddSslName, appAddSslType, appAddSslKey, appAddSslCrt, appAddSslCabundle, appAddSslAutoUrl string
var appAddSslAutoUrlLink, appAddSslForce bool
var PossibleSslTypes = []string{"letsencrypt", "xolphin", "own"}
var appCertificateAddCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add new ssl certificate to an app.",
	Example: "lvl app ssl add MyAppName -n mySslCertificateName",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Search appId based on name
		appId := resolveApp(args[0])

		var certificate types.SslCertificate
		// checking if the chosen type is one of the valid options.
		var isSslTypeValid bool = false
		for _, sslType := range PossibleSslTypes {
			// when type is valid.
			// check if type is 'own' or not. -> own -> needs special request data
			if appAddSslType == sslType {
				isSslTypeValid = true
				// type own -> different properties
				if appAddSslType == "own" {
					request := types.AppSslCertificateTypeOwnRequest{
						Name:                   appAddSslName,
						SslType:                appAddSslType,
						AutoSslCertificateUrls: appAddSslAutoUrl,
						SslKey:                 appAddSslKey,
						SslCrt:                 appAddSslCrt,
						SslCabundle:            appAddSslCabundle,
						AutoUrlLink:            appAddSslAutoUrlLink,
						SslForce:               appAddSslForce,
					}
					certificate = Level27Client.AppCertificateAdd(appId, request)
					// type letsencrypt or xolphin
				} else {
					request := types.AppSslCertificateRequest{
						Name:                   appAddSslName,
						SslType:                appAddSslType,
						AutoSslCertificateUrls: appAddSslAutoUrl,
						AutoUrlLink:            appAddSslAutoUrlLink,
						SslForce:               appAddSslForce,
					}
					certificate = Level27Client.AppCertificateAdd(appId, request)
				}

			}
		}

		// error when given type is not valid
		if !isSslTypeValid {
			log.Fatal(fmt.Sprintf("Given sslType: '%v' is NOT valid.", appAddSslType))
		}

		message := fmt.Sprintf("sslCertificate created: [name: '%v' - ID: '%v'].", certificate.Name, certificate.ID)
		fmt.Println(message)

	},
}

// ---- DELETE SSL CERTIFICATE
var appCertificateDeleteConfirmed bool
var appCertificateDeleteCmd = &cobra.Command{
	Use:     "delete [appID] [CertificateID]",
	Short:   "Delete a ssl certificate from an app.",
	Example: "lvl app ssl delete MyAppName --yes",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Search appId based on name
		appId := resolveApp(args[0])
		//check for valid certificateID
		certificateID := checkSingleIntID(args[1], "appCertificate")

		Level27Client.AppCertificateDelete(appId, certificateID, appCertificateDeleteConfirmed)
	},
}

// ---- FIX SSL CERTIFICATE
var appCertificateFixCmd = &cobra.Command{
	Use:     "fix",
	Short:   "Fix an invalid certificate.",
	Example: "lvl app ssl fix MyAppName 3022",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Search appId based on name
		appId := resolveApp(args[0])
		//check for valid certificateID
		certificateID := checkSingleIntID(args[1], "appCertificate")

		Level27Client.AppCertificateFix(appId, certificateID)
	},
}

// ---- GET PRIVATE KEY FOR TYPE 'OWN' SSL CERTIFICATES
var appCertificateKeyCmd = &cobra.Command{
	Use:     "key",
	Short:   "Return a private key for type 'own' sslCertificate.",
	Example: "lvl app ssl key MyAppName",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Search appId based on name
		appId := resolveApp(args[0])
		//check for valid certificateID
		certificateID := checkSingleIntID(args[1], "appCertificate")

		Level27Client.AppCertificateKey(appId, certificateID)
	},
}

// #endregion
//-------------------------------------------------  APP SSL CERTIFICATES (ACTIONS) -------------------------------------------------
// ---- ACTIONS (CREATE JOB FOR CERTIFICATE)
var appCertificateActionCmd = &cobra.Command{
	Use:   "action",
	Short: "commands to create a job for a ssl certificate.",
}

// #region APP SSL CERTIFICATES (ACTIONS)

// ---- ACTION RETRY
var appCertificateActionRetryCmd = &cobra.Command{
	Use:     "retry [app] [certificateID]",
	Short:   "Create 'retry' job for ssl certificate.",
	Example: "lvl app ssl action retry MyAppName 3023",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Search appId based on name
		appId := resolveApp(args[0])
		//check for valid certificateID
		certificateID := checkSingleIntID(args[1], "appCertificate")

		Level27Client.AppCertificateAction(appId, certificateID, "retry")
	},
}

// ---- ACTION VALIDATECHALLENGE
var appCertificateActionValidateCmd = &cobra.Command{
	Use:     "validateChallenge",
	Short:   "Create 'validateChallenge' job for ssl certificate.",
	Example: "lvl app ssl action validateChallenge MyAppName 1603",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Search appId based on name
		appId := resolveApp(args[0])
		//check for valid certificateID
		certificateID := checkSingleIntID(args[1], "appCertificate")

		Level27Client.AppCertificateAction(appId, certificateID, "validateChallenge")
	},
}

// #endregion

//-------------------------------------------------  APP COMPONENT RESTORE (GET / CREATE / UPDATE / DELETE / DOWNLOAD) -------------------------------------------------
// ---- RESTORE COMMAND
var appComponentRestoreCmd = &cobra.Command{
	Use:     "restores",
	Short:   "Command to manage restores on an app.",
	Example: "lvl app restore [subcommand]",
}

// ---- GET LIST OF RESTORES
var appComponentRestoreGetCmd = &cobra.Command{
	Use:     "get",
	Short:   "Show a list of al available restores on an app.",
	Example: "lvl app restore get NameOfMyApp",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// search for appId based on name
		appId := resolveApp(args[0])

		Restores := Level27Client.AppComponentRestoresGet(appId)

		outputFormatTableFuncs(Restores,
			[]string{"ID", "FILENAME", "STATUS", "DATE", "APPCOMPONENT_ID", "APPCOMPONENT_NAME"},
			[]interface{}{"ID", "Filename", "Status", func(r types.AppComponentRestore) string { return utils.FormatUnixTime(r.AvailableBackup.Date) }, "Appcomponent.ID", "Appcomponent.Name"})
	},
}

// ---- CREATE A NEW RESTORE
var appRestoreCreateComponent string
var appRestoreCreateBackup int
var appComponentRestoreCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new restore for an app.",
	Example: "lvl app restore create MyAppName MyComponentName 453",
	Args:    cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		//search AppId based on appname
		appId := resolveApp(args[0])
		// search componentId based on name
		componentId := resolveAppComponent(appId, args[1])
		backupId, err := strconv.Atoi(args[2])

		if err != nil {
			log.Fatalf("BackupID is NOT valid! '%v'.", args[2])
		}

		request := types.AppComponentRestoreRequest{
			Appcomponent:    componentId,
			AvailableBackup: backupId,
		}
		restore := Level27Client.AppComponentRestoreCreate(appId, request)

		log.Printf("Restore created. [ID: %v].", restore.ID)
	},
}

// ---- DELETE A RESTORE
var isAppRestoreDeleteConfirmed bool
var appRestoreDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a specific restore from an app.",
	Example: "lvl app component restore delete MyAppName 4532",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// search for appId based on name
		appId := resolveApp(args[0])

		// check if restoreId is valid type
		restoreId := checkSingleIntID(args[1], "restore")

		Level27Client.AppComponentRestoresDelete(appId, restoreId, isAppRestoreDeleteConfirmed)

	},
}

// ---- DOWNLOAD RESTORE FILE
var appComponentRestoreDownloadName string
var appComponentRestoreDownloadCmd = &cobra.Command{
	Use:     "download [appname] [restoreID]",
	Short:   "Download the restore file.",
	Example: "lvl app component restore download MyAppName 4123",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// search appId based on name
		appId := resolveApp(args[0])
		// check if restoreId is valid type
		restoreId := checkSingleIntID(args[1], "Restore")

		Level27Client.AppComponentRestoreDownload(appId, restoreId)
	},
}

//-------------------------------------------------  APP COMPONENT BACKUPS (GET) -------------------------------------------------
var appComponentNameBackup string
var appComponentBackupsCmd = &cobra.Command{
	Use:     "backup",
	Short:   "Commands for managing availableBackups.",
	Example: "lvl app component backup get MyAppName MyComponentName",
}

var appComponentBackupsGetCmd = &cobra.Command{
	Use:     "get",
	Short:   "Show list of available backups.",
	Example: "lvl app component backup get MyAppName MyComponentName",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// search appId based on appname
		appId := resolveApp(args[0])
		// search componentId based on name
		componentId := resolveAppComponent(appId, args[1])

		availableBackups := Level27Client.AppComponentbackupsGet(appId, componentId)

		outputFormatTableFuncs(availableBackups, 
			[]string{"ID", "SNAPSHOTNAME", "DATE"}, 
			[]interface{}{"ID", "SnapshotName", func(a types.AppComponentAvailableBackup) string { return utils.FormatUnixTime(a.Date) 
			}})
	},
}


//-------------------------------------------------  APP MIGRATIONS (GET / DESCRIBE / CREATE / UPDATE) -------------------------------------------------
// ---- MIGRATION COMMAND
var appMigrationsCmd = &cobra.Command{
	Use: "migrations",
	Short: "Commands to manage app migrations.",
	Example: "lvl app migrations get MyAppName\nlvl app migrations describe MyAppName 1513",
}

// ---- GET LIST OF MIGRATIONS 
var appMigrationsGetCmd = &cobra.Command{
	Use: "get [appName]",
	Short: "Show a list of all available migrations.",
	Example: "lvl app migrations get MyAppName",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//search for AppId based on name
		appId := resolveApp(args[0])

		migrations := Level27Client.AppMigrationsGet(appId)

		outputFormatTableFuncs(migrations, 
			[]string{"ID","MIGRATION_TYPE", "STATUS", "DATE_PLANNED"}, 
			[]interface{}{"ID", "MigrationType", "Status", func(m types.AppMigration) string {return utils.FormatUnixTime(m.DtPlanned)
			}})
	},

}