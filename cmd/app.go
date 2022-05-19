package cmd

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

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
	appComponentCreateCmd.Flags().StringVarP(&appComponentCreateParamsFile, "params-file", "f", "", "JSON file to read params from. Pass '-' to read from stdin.")
	appComponentCreateCmd.Flags().StringVar(&appComponentCreateName, "name", "", "")
	appComponentCreateCmd.Flags().StringVar(&appComponentCreateType, "type", "", "")
	appComponentCreateCmd.Flags().StringVar(&appComponentCreateSystem, "system", "", "")
	appComponentCreateCmd.Flags().StringVar(&appComponentCreateSystemgroup, "systemgroup", "", "")
	appComponentCreateCmd.Flags().IntVar(&appComponentCreateSystemprovider, "systemprovider", 0, "")
	appComponentCreateParams = appComponentCreateCmd.Flags().StringArray("param", nil, "")
	appComponentCreateCmd.MarkFlagRequired("name")
	appComponentCreateCmd.MarkFlagRequired("type")

	// ---- UPDATE COMPONENT
	appComponentCmd.AddCommand(appComponentUpdateCmd)
	settingsFileFlag(appComponentUpdateCmd)
	settingString(appComponentUpdateCmd, updateSettings, "name", "New name for this app component")
	appComponentUpdateParams = appComponentUpdateCmd.Flags().StringArray("param", nil, "")


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

	// APP INTEGRITY
	addIntegrityCheckCmds(appCmd, "apps", resolveApp)

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


	// APP COMPONENT URL
	appComponentCmd.AddCommand(appComponentUrlCmd)

	// APP COMPONENT URL GET
	appComponentUrlCmd.AddCommand(appComponentUrlGetCmd)
	addCommonGetFlags(appComponentUrlGetCmd)

	// APP COMPONENT URL CREATE
	appComponentUrlCmd.AddCommand(appComponentUrlCreateCmd)
	appComponentUrlCreateCmd.Flags().BoolVar(&appComponentUrlCreateAuthentication, "authentication", false, "Require HTTP Basic authentication on the URL")
	appComponentUrlCreateCmd.Flags().StringVarP(&appComponentUrlCreateContent, "content", "c", "", "Content for the new URL")
	appComponentUrlCreateCmd.Flags().BoolVar(&appComponentUrlCreateSslForce, "force-ssl", false, "Force usage of SSL on the URL")
	appComponentUrlCreateCmd.Flags().IntVar(&appComponentUrlCreateSslCertificate, "ssl-certificate", 0, "SSL certificate to use.")
	appComponentUrlCreateCmd.Flags().BoolVar(&appComponentUrlCreateHandleDns, "handle-dns", false, "Automatically create DNS records")
	appComponentUrlCreateCmd.Flags().BoolVar(&appComponentUrlCreateAutoSslCertificate, "auto-ssl-certificate", false, "Automatically create SSL certificate with Let's Encrypt")

	// APP COMPONENT URL DELETE
	appComponentUrlCmd.AddCommand(appComponentUrlDeleteCmd)
	appComponentUrlDeleteCmd.Flags().BoolVar(&appComponentUrlDeleteForce, "force", false, "Do not ask for confirmation to delete the URL")

	//-------------------------------------------------  APP MIGRATIONS (GET / DESCRIBE / CREATE / UPDATE) -------------------------------------------------
	// ---- MIGRATIONS COMMAND
	appCmd.AddCommand(appMigrationsCmd)

	// ---- GET LIST OF MIGRATIONS
	appMigrationsCmd.AddCommand(appMigrationsGetCmd)

	// ---- CREATE NEW APP MIGRATION
	appMigrationsCmd.AddCommand(appMigrationsCreateCmd)
	// flags needed to create new migration
	flags = appMigrationsCreateCmd.Flags()
	flags.StringVarP(&appMigrationCreatePlanned, "planned", "", "", "DateTime - timestamp.")
	flags.StringArrayVarP(&appMigrationCreateItems, "migration-item", "", []string{}, "Migration items. each item should contain at least a 'source' (the component to migrate) and a 'destSystem' or 'destGroup' to migrate to.")

	// ---- UPDATE MIGRATION
	appMigrationsCmd.AddCommand(appMigrationsUpdateCmd)
	flags = appMigrationsUpdateCmd.Flags()
	flags.StringVarP(&appMigrationsUpdateDtPlanned, "planned", "", "", "DateTime - timestamp.")
	flags.StringVarP(&appMigrationsUpdateType, "type", "t", "", "Migration type. (one of automatic (all migration steps are done automatically), confirmed (a user has to confirm each migration step)).")
	appMigrationsUpdateCmd.MarkFlagRequired("type")
	appMigrationsUpdateCmd.MarkFlagRequired("planned")

	// ---- DESCRIBE MIGRATION
	appMigrationsCmd.AddCommand(appMigrationDescribeCmd)
	//-------------------------------------------------  APP MIGRATIONS ACTIONS (CONFIRM / DENY / RESTART) -------------------------------------------------
	// ---- MIGRATION ACTION COMMAND
	appMigrationsCmd.AddCommand(appMigrationsActionCmd)

	// ---- MIGRATION ACTION CONFIRM
	appMigrationsActionCmd.AddCommand(appMigrationsActionConfirmCmd)
	// ---- MIGRATION ACTION DENY
	appMigrationsActionCmd.AddCommand(appMigrationsActionDenyCmd)
	// ---- MIGRATION ACTION RETRY
	appMigrationsActionCmd.AddCommand(appMigrationsActionRetryCmd)
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

// GET AN APP URL ID BASED ON THE NAME
func resolveAppComponentUrl(appId int, appComponentId int, arg string) int {
	// if arg already int, this is the ID
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id
	}

	return resolveShared(
		Level27Client.AppComponentUrlLookup(appId, appComponentId, arg),
		arg,
		"url",
		func(url types.AppComponentUrlShort) string { return fmt.Sprintf("%s (%d)", url.Content, url.ID) }).ID
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
	Example: "lvl app get -f FilterByName",
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
var appComponentCreateParamsFile string
var appComponentCreateName string
var appComponentCreateType string
var appComponentCreateSystem string
var appComponentCreateSystemgroup string
var appComponentCreateSystemprovider int
var appComponentCreateParams *[]string
var appComponentCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new appcomponent.",
	Example: "lvl app component create --name myComponentName --type docker",

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if appComponentCreateSystem == "" && appComponentCreateSystemgroup == "" {
			cobra.CheckErr("Must specify either a system or a system group")
		}

		if appComponentCreateSystem != "" && appComponentCreateSystemgroup != "" {
			cobra.CheckErr("Cannot specify both a system and a system group")
		}

		appID := resolveApp(args[0])
		componentTypes := Level27Client.AppComponenttypesGet();

		val, ok := componentTypes[appComponentCreateType]
		if !ok {
			cobra.CheckErr(fmt.Sprintf("Unknown component type: %s", appComponentCreateType))
		}

		paramsPassed := loadSettings(appComponentCreateParamsFile)

		// Parse params from command line
		for _, param := range *appComponentCreateParams {
			split := strings.SplitN(param, "=", 2)
			if len(split) != 2 {
				cobra.CheckErr(fmt.Sprintf("Expected key=value pair to --param: %s", param))
			}

			paramsPassed[split[0]] = readArgFileSupported(split[1])
		}

		create := map[string]interface{}{}
		create["name"] = appComponentCreateName
		create["category"] = "config"
		create["appcomponenttype"] = appComponentCreateType

		if appComponentCreateSystem != "" {
			create["system"] = resolveSystem(appComponentCreateSystem)
		}

		if appComponentCreateSystemgroup != "" {
			create["systemgroup"] = checkSingleIntID(appComponentCreateSystemgroup, "systemgroup")
		}

		if appComponentCreateSystemprovider != 0 {
			create["systemprovider"] = appComponentCreateSystemprovider
		}

		// Go over specified commands in app component types to validate and map data.

		paramNames := map[string]bool{}
		for _, param := range val.Servicetype.Parameters {
			paramName := param.Name
			paramNames[paramName] = true
			paramValue, hasValue := paramsPassed[paramName]
			if hasValue {
				if param.Readonly || param.DisableEdit {
					cobra.CheckErr(fmt.Sprintf("Param cannot be changed: %s", paramName))
				}
				create[paramName] = parseComponentParameter(param, paramValue)
			} else if param.Required && param.DefaultValue == nil {
				cobra.CheckErr(fmt.Sprintf("Required parameter not given: %s", paramName))
			}
		}

		// Check that there aren't any params given that don't exist.
		for k := range paramsPassed {
			if !paramNames[k] {
				cobra.CheckErr(fmt.Sprintf("Unknown parameter given: %s", k))
			}
		}

		Level27Client.AppComponentCreate(appID, create)
	},
}

var appComponentUpdateParams *[]string
var appComponentUpdateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a new appcomponent.",
	Example: "",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		settings := loadMergeSettings(updateSettingsFile, updateSettings)

		appID := resolveApp(args[0])
		appComponentID := resolveAppComponent(appID, args[1])
		componentTypes := Level27Client.AppComponenttypesGet();

		appComponent := Level27Client.AppComponentGetSingle(appID, appComponentID)

		parameterTypes := make(map[string]types.AppComponentTypeParameter)
		for _, param := range componentTypes[appComponent.Appcomponenttype].Servicetype.Parameters {
			parameterTypes[param.Name] = param
		}

		data := make(map[string]interface{})
		data["appcomponenttype"] = appComponent.Appcomponenttype
		data["name"] = appComponent.Name
		data["category"] = appComponent.Category
		data["systemgroup"] = appComponent.Systemgroup
		data["system"] = nil
		if appComponent.Systemgroup == nil {
			data["system"] = appComponent.Systems[0].ID
		}

		for k, v := range appComponent.Appcomponentparameters {
			param := parameterTypes[k]
			paramType := parameterTypes[k].Type

			if param.Readonly {
				continue
			}

			switch (paramType) {
			case "password-sha512", "password-plain", "password-sha1", "passowrd-sha256-scram":
			// Passwords are sent as "******". Skip them to avoid getting API errors.
				continue
			case "sshkey[]":
				// Need to map SSH keys -> IDs
				sshKeys := v.([]interface{})
				ids := make([]int, len(sshKeys))
				for i, sshKey := range sshKeys {
					keyCast := sshKey.(map[string]interface{})
					ids[i] = int(keyCast["id"].(float64))
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
				cobra.CheckErr(fmt.Sprintf("Expected key=value pair to --param: %s", param))
			}

			paramName := split[0]
			paramValue := readArgFileSupported(split[1])
			paramType, ok := parameterTypes[paramName]
			if !ok {
				cobra.CheckErr(fmt.Sprintf("Unknown parameter: %s", paramName))
			}

			data[paramName] = parseComponentParameter(paramType, paramValue)
		}

		Level27Client.AppComponentUpdate(appID, appComponentID, data)
	},
}


func parseComponentParameter(param types.AppComponentTypeParameter, paramValue interface{}) interface{} {
	// Convert parameters to the correct types in-JSON.
	var str string
	var ok bool
	if str, ok = paramValue.(string); !ok {
		// Value isn't a string. This means it must have come from a JSON input file or something (i.e. not command line arg)
		// So assume it's the correct type and let the API complain if it isn't.
		return paramValue
	}

	switch param.Type {
	case "sshkey[]":
		keys := []int{}
		for _, key := range strings.Split(str, ",") {
			// TODO: Resolve SSH key
			keys = append(keys, checkSingleIntID(key, "SSH key"))
		}
		return keys
	case "integer":
		intVal, err := strconv.Atoi(str)
		cobra.CheckErr(err)
		return intVal
	case "boolean":
		return strings.EqualFold(str, "true")
	case "array":
		found := false
		for _, possibleValue := range param.PossibleValues {
			if str == possibleValue {
				found = true
				break
			}
		}

		if !found {
			cobra.CheckErr(
				fmt.Sprintf(
					"Parameter %s: value '%s' not in range of possible values: %s",
					param.Name,
					str,
					strings.Join(param.PossibleValues, ", ")))
		}

		return str
	default:
		// Pass as string
		return str
	}
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

		Level27Client.AppComponentRestoreDownload(appId, restoreId, appComponentRestoreDownloadName)
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
			[]interface{}{"ID", "SnapshotName", func(a types.AppComponentAvailableBackup) string {
				return utils.FormatUnixTime(a.Date)
			}})
	},
}

//-------------------------------------------------  APP MIGRATIONS (GET / DESCRIBE / CREATE / UPDATE) -------------------------------------------------
// ---- MIGRATION COMMAND
var appMigrationsCmd = &cobra.Command{
	Use:     "migrations",
	Short:   "Commands to manage app migrations.",
	Example: "lvl app migrations get MyAppName\nlvl app migrations describe MyAppName 1513",
}

// ---- GET LIST OF MIGRATIONS
var appMigrationsGetCmd = &cobra.Command{
	Use:     "get [appName]",
	Short:   "Show a list of all available migrations.",
	Example: "lvl app migrations get MyAppName",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//search for AppId based on name
		appId := resolveApp(args[0])

		migrations := Level27Client.AppMigrationsGet(appId)

		outputFormatTableFuncs(migrations,
			[]string{"ID", "MIGRATION_TYPE", "STATUS", "DATE_PLANNED"},
			[]interface{}{"ID", "MigrationType", "Status", func(m types.AppMigration) string {
				return utils.FormatUnixTime(m.DtPlanned)
			}})
	},
}

// --- CREATE MIGRATION
var appMigrationCreatePlanned string
var appMigrationCreateItems []string
var appMigrationsCreateCmd = &cobra.Command{
	Use:     "create [appName] [flags]",
	Short:   "Create a new app migration.",
	Long:    `Items to migrate are specified with --migration-item, taking a parameter in a comma-separated key=value format. Multiple items can be migrated at once by specifying --migration-item multiple times.`,
	Example: "lvl app migrations create MyAppName --migration-item 'source=forum, destSystem=newForumSystem' --migration-item 'source=database, destGroup=newDbGroup, ord=2'",
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		//search for appid based on appName
		appId := resolveApp(args[0])

		items := []types.AppMigrationItem{}

		for _, migrationItem := range appMigrationCreateItems {
			items = append(items, ParseMigrationItem(appId, migrationItem))
		}

		request := types.AppMigrationRequest{
			MigrationType:      "automatic",
			DtPlanned:          appMigrationCreatePlanned,
			MigrationItemArray: items,
		}

		Level27Client.AppMigrationsCreate(appId, request)
	},
}

func ParseMigrationItem(appID int, values string) types.AppMigrationItem{
	valueSplitted := strings.Split(values, ",")

	item := types.AppMigrationItem{
		Ord: 1,
		Source: "cp4",
	}

	haveAnyDst := false
	haveAnySrc := false
	for _, keyValuePair := range valueSplitted {
		// Go over key value pairs and fill out the migration item as we go.

		key, value := ParseMigrationItemKeyValuePair(keyValuePair)

		switch key {
		case "ord":
			val, err := strconv.Atoi(value)
			cobra.CheckErr(err)
			item.Ord = val

		case "destSystem":
			item.DestinationEntityId = resolveSystem(value)
			item.DestinationEntity = "system"
			haveAnyDst = true

		case "destGroup":
			item.DestinationEntityId = resolveSystemgroup(value)
			item.DestinationEntity = "systemgroup"
			haveAnyDst = true

		case "source":
			appComponent := resolveAppComponent(appID, value)
			appComponentType := Level27Client.AppComponentGetSingle(appID, appComponent).Appcomponenttype
			haveAnySrc = true

			item.SourceInfo = appComponent
			item.Type = appComponentType

		default:
			log.Fatalf("Unknown property in migration item: %s", key)
		}
	}

	if !haveAnyDst {
		cobra.CheckErr(fmt.Sprintf("No destination specified for migration item!"))
	}

	if !haveAnySrc {
		cobra.CheckErr(fmt.Sprintf("No source specified for migration item!"))
	}

	return item
}


func ParseMigrationItemKeyValuePair(keyValuePair string) (string, string) {
	split := strings.SplitN(keyValuePair, "=", 2)
	if len(split) == 1 {
		log.Fatalf("MigrationItem property not defined correctly: '%v'. Use '=' to define properties.", keyValuePair)
	}

	key := strings.TrimSpace(split[0])
	value := strings.TrimSpace(split[1])

	return key, value
}

// ---- UPDATE MIGRATION
var appMigrationsUpdateType, appMigrationsUpdateDtPlanned string
var appMigrationsUpdateCmd = &cobra.Command{
	Use:     "update [appID] [migrationID]",
	Short:   "Update an app migration.",
	Example: "lvl app migrations update MyAppName 3414",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		//search for appID based on name
		appId := resolveApp(args[0])
		// check for valid migrationId type
		migrationId := checkSingleIntID(args[1], "appMigration")

		request := types.AppMigrationUpdate{
			MigrationType: appMigrationsUpdateType,
			DtPlanned:     appMigrationsUpdateDtPlanned,
		}

		Level27Client.AppMigrationsUpdate(appId, migrationId, request)
	},
}

// ---- DESCRIBE MIGRATION
var appMigrationDescribeCmd = &cobra.Command{
	Use:     "describe [appID] [migrationID]",
	Short:   "Get detailed info about a specific migration.",
	Example: "lvl app migrations describe MyAppName 1243",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// search for appId based on name
		appId := resolveApp(args[0])
		// check for valid migrationId type
		migrationId := checkSingleIntID(args[1], "appMigration")

		migration := Level27Client.AppMigrationDescribe(appId, migrationId)

		outputFormatTemplate(migration, "templates/appMigration.tmpl")
	},
}

//-------------------------------------------------  APP MIGRATIONS ACTIONS (CONFIRM / DENY / RESTART) -------------------------------------------------
// ---- MIGRATIONS ACTION COMMAND
var appMigrationsActionCmd = &cobra.Command{
	Use:     "action",
	Short:   "Execute an action for a migration",
	Example: "lvl app migrations action deny MyAppName 241\nlvl app migrations action restart MyAppName 234",
}

// ---- CONFIRM MIGRATION
var appMigrationsActionConfirmCmd = &cobra.Command{
	Use:     "confirm",
	Short:   "Execute confirm action on a migration",
	Example: "lvl app migrations action confirm MyAppName 332",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// search for appId based on name
		appId := resolveApp(args[0])
		// check for valid migrationId type
		migrationId := checkSingleIntID(args[1], "appMigration")

		Level27Client.AppMigrationsAction(appId, migrationId, "confirm")
	},
}

// ---- DENY MIGRATION
var appMigrationsActionDenyCmd = &cobra.Command{
	Use:     "deny",
	Short:   "Execute confirm action on a migration",
	Example: "lvl app migrations action deny MyAppName 332",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// search for appId based on name
		appId := resolveApp(args[0])
		// check for valid migrationId type
		migrationId := checkSingleIntID(args[1], "appMigration")

		Level27Client.AppMigrationsAction(appId, migrationId, "deny")
	},
}

// ---- RETRY MIGRATION
var appMigrationsActionRetryCmd = &cobra.Command{
	Use:     "retry",
	Short:   "Execute confirm action on a migration",
	Example: "lvl app migrations action retry MyAppName 332",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// search for appId based on name
		appId := resolveApp(args[0])
		// check for valid migrationId type
		migrationId := checkSingleIntID(args[1], "appMigration")

		Level27Client.AppMigrationsAction(appId, migrationId, "retry")
	},
}

// ------- APP COMPONENT URLs

// APP COMPONENT URL
var appComponentUrlCmd = &cobra.Command{
	Use: "url",
	Aliases: []string{"urls"},
}

// APP COMPONENT URL GET
var appComponentUrlGetCmd = &cobra.Command{
	Use: "get",

	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		appID := resolveApp(args[0])
		componentID := resolveAppComponent(appID, args[1])

		results := resolveGets(
			args[2:],
			func (name string) []types.AppComponentUrlShort {
				return Level27Client.AppComponentUrlLookup(appID, componentID, name)
			},
			func (i int) types.AppComponentUrlShort {
				return Level27Client.AppComponentUrlGetSingle(appID, componentID, i).ToShort()
			},
			func(cgp types.CommonGetParams) []types.AppComponentUrlShort {
				return Level27Client.AppComponentUrlGetList(appID, componentID, cgp)
			},
		)

		outputFormatTable(
			results,
			[]string{"ID", "CONTENT", "STATUS", "TYPE", "SSL CERT", "FORCE SSL", "HANDLE DNS", "AUTHENTICATE"},
			[]string{"ID", "Content", "Status", "Type", "SslCertificate.Name", "SslForce", "HandleDNS", "Authentication"})
	},
}

// APP COMPONENT URL CREATE
var appComponentUrlCreateAuthentication bool
var appComponentUrlCreateContent string
var appComponentUrlCreateSslForce bool
var appComponentUrlCreateSslCertificate int
var appComponentUrlCreateHandleDns bool
var appComponentUrlCreateAutoSslCertificate bool
var appComponentUrlCreateCmd = &cobra.Command{
	Use: "create",
	Short: "Create an url for an appcomponent.",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		appID := resolveApp(args[0])
		componentID := resolveAppComponent(appID, args[1])

		var cert *int
		if appComponentUrlCreateSslCertificate == 0 {
			cert = nil
		} else {
			cert = &appComponentUrlCreateSslCertificate
		}

		create := types.AppComponentUrlCreate {
			Authentication: appComponentUrlCreateAuthentication,
			Content: appComponentUrlCreateContent,
			SslForce: appComponentUrlCreateSslForce,
			SslCertificate: cert,
			HandleDns: appComponentUrlCreateHandleDns,
			AutoSslCertificate: appComponentUrlCreateAutoSslCertificate,
		}

		Level27Client.AppComponentUrlCreate(appID, componentID, create)
	},
}

// APP COMPONENT URL DELETE
var appComponentUrlDeleteForce bool
var appComponentUrlDeleteCmd = &cobra.Command{
	Use: "delete",
	Short: "Delete an url from an appcomponent.",
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		appID := resolveApp(args[0])
		componentID := resolveAppComponent(appID, args[1])
		urlID := resolveAppComponentUrl(appID, componentID, args[2])

		if !appComponentUrlDeleteForce {
			url := Level27Client.AppComponentUrlGetSingle(appID, componentID, urlID)

			msg := fmt.Sprintf(
				"Delete URL %s (%d) on app comp %s (%d)?",
				url.Content, url.ID,
				url.Appcomponent.Name, url.Appcomponent.ID)

			if !confirmPrompt(msg) {
				return
			}
		}

		Level27Client.AppComponentUrlDelete(appID, componentID, urlID)
	},
}