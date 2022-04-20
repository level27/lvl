package cmd

import (
	"fmt"
	"log"
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

	//------------------------------------------------- APP (GET / CREATE / DELETE / UPDATE / DESCRIBE)-------------------------------------------------
	// ACTION COMMAND
	appCmd.AddCommand(AppActionCmd)

	// ACTIVATE APP
	AppActionCmd.AddCommand(AppActionActivateCmd)

	// DEACTIVATE APP
	AppActionCmd.AddCommand(AppActionDeactivateCmd)

	// APP ACCESS
	addAccessCmds(appCmd, "apps", resolveApp)


	// APP SSLCERT
	appCmd.AddCommand(appSslcertCmd)

	// APP SSLCERT GET
	appSslcertCmd.AddCommand(appSslcertGetCmd)
	addCommonGetFlags(appSslcertCmd)

	// APP SSLCERT DESCRIBE
	appSslcertCmd.AddCommand(appSslcertDescribeCmd)

	// APP SSLCERT ADD
	appSslcertCmd.AddCommand(appSslcertCreateCmd)
	appSslcertCreateCmd.Flags().StringVar(&appSslcertCreateName, "name", "", "Name of this SSL certificate")
	appSslcertCreateCmd.Flags().StringVar(&appSslcertCreateSslType, "type", "", "Type of SSL certificate to use. Options are: letsencrypt, xolphin, own")
	appSslcertCreateCmd.Flags().StringVar(&appSslcertCreateAutoSslCertificateUrls, "auto-urls", "", "URL or CSV list of URLs (required for Let's Encrypt)")
	appSslcertCreateCmd.Flags().StringVar(&appSslcertCreateSslKey, "ssl-key", "", "SSL key for own certificate")
	appSslcertCreateCmd.Flags().StringVar(&appSslcertCreateSslCabundle, "ssl-cabundle", "", "SSL CA bundle for own certificate")
	appSslcertCreateCmd.Flags().StringVar(&appSslcertCreateSslCrt, "ssl-crt", "", "SSL CRT for own certificate. Can be read from a file by specifying @filename.")
	appSslcertCreateCmd.Flags().BoolVar(&appSslcertCreateAutoUrlLink, "auto-link", false, "After creation, automatically link to any URLs without existing certificate")
	appSslcertCreateCmd.Flags().BoolVar(&appSslcertCreateSslForce, "ssl-force", false, "Force SSL")

	// APP SSLCERT DELETE

	appSslcertCmd.AddCommand(appSslcertDeleteCmd)
	appSslcertDeleteCmd.Flags().BoolVar(&appSslCertDeleteForce, "force", false, "Do not ask for confirmation to delete the SSL certificate")

	// APP SSLCERT UPDATE
	appSslcertCmd.AddCommand(appSslCertUpdateCmd)
	settingsFileFlag(appSslCertUpdateCmd)
	settingString(appSslCertUpdateCmd, updateSettings, "name", "New name for the SSL certificate")

	appSslcertCmd.AddCommand(appSslCertFixCmd)

	// APP SSLCERT ACTIONS
	appSslcertCmd.AddCommand(appSslCertActionCmd)

	appSslCertActionCmd.AddCommand(appSslCertActionRetryCmd)
	appSslCertActionCmd.AddCommand(appSslCertActionValidateChallengeCmd)

}

func resolveApp(arg string) int {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id
	}

	app := Level27Client.AppLookup(arg)
	if app == nil {
		cobra.CheckErr(fmt.Sprintf("Unable to find app: %s", arg))
		return 0
	}
	return app.ID
}

func resolveAppSslCertificate(appID int, arg string) int {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id
	}

	cert := Level27Client.AppSslCertificatesLookup(appID, arg)
	if cert == nil {
		cobra.CheckErr(fmt.Sprintf("Unable to find app SSL certificate: %s", arg))
		return 0
	}
	return cert.ID
}


// MAIN COMMAND APPS
var appCmd = &cobra.Command{
	Use:     "app",
	Short:   "Commands to manage apps",
	Example: "lvl app [subcommmand]",
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
	Example: "lvl app delete 2593",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// check for valid appID type
		appId := checkSingleIntID(args[0], "app")

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
		appId := checkSingleIntID(args[0], "app")

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
		appId := checkSingleIntID(args[0], "app")
		// get all data from app by appId
		app := Level27Client.App(appId)
		outputFormatTemplate(app, "templates/app.tmpl")
	},
}

// #endregion

//------------------------------------------------- APP ACTIONS (GET / CREATE  / UPDATE / DELETE / DESCRIBE)-------------------------------------------------

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
		appId := checkSingleIntID(args[0], "app")

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
		appId := checkSingleIntID(args[0], "app")

		Level27Client.AppAction(appId, "deactivate")
	},
}


// APP SSL CERTIFICATES

// APP SSLCERT
var appSslcertCmd = &cobra.Command{
	Use: "sslcert",
	Short: "Commands for managing SSL certificates on apps",
	Example: "lvl app sslcert get forum\nlvl app sslcert describe forum forum.example.com",

	Aliases: []string{"ssl"},
}

// APP SSLCERT GET
var appSslcertGetType string
var appSslcertGetStatus string
var appSslcertGetCmd = &cobra.Command{
	Use: "get [app]",
	Short: "Get a list of SSL certificates for an app",
	Example: "lvl app sslcert get forum\nlvl app sslcert get forum -f admin",

	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appID := resolveApp(args[0])

		certs := resolveGets(
			// First arg is app ID.
			args[1:],
			func(name string) *types.AppSslCertificate { return Level27Client.AppSslCertificatesLookup(appID, name) },
			func(certID int) types.AppSslCertificate { return Level27Client.AppSslCertificatesGetSingle(appID, certID)},
			func(get types.CommonGetParams) []types.AppSslCertificate { return Level27Client.AppSslCertificatesGetList(appID, appSslcertGetType, appSslcertGetStatus, get) },
			)

		outputFormatTable(certs, []string{"ID", "Name", "Status", "SSL Status"}, []string{"ID", "Name", "Status", "SslStatus"})
	},
}

// APP SSLCERT DESCRIBE

var appSslcertDescribeCmd = &cobra.Command{
	Use: "describe [app] [SSL cert]",
	Short: "Get detailed information of an SSL certificate",
	Example: "lvl app sslcert describe forum forum.example.com",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		appID := resolveApp(args[0])
		certID := resolveAppSslCertificate(appID, args[1])

		cert := Level27Client.AppSslCertificatesGetSingle(appID, certID)

		outputFormatTemplate(cert, "templates/appSslCertificate.tmpl")
	},
}

// APP SSLCERT CREATE
var appSslcertCreateName string
var appSslcertCreateSslType string
var appSslcertCreateAutoSslCertificateUrls string
var appSslcertCreateSslKey string
var appSslcertCreateSslCrt string
var appSslcertCreateSslCabundle string
var appSslcertCreateAutoUrlLink bool
var appSslcertCreateSslForce bool

var appSslcertCreateCmd = &cobra.Command {
	Use: "create [app]",
	Short: "Create a new SSL certificate on an app",
	Example: "lvl app sslcert create forum --name forum.example.com --auto-urls forum.example.com --auto-link --type letsencrypt\nlvl app sslcert create forum --name forum.example.com --type own --ssl-cabundle '@cert.ca-bundle' --ssl-key '@key.pem' --ssl-crt '@cert.crt'",

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appID := resolveApp(args[0])

		create := types.AppSslCertificateCreate {
			Name: appSslcertCreateName,
			SslType: appSslcertCreateSslType,
			AutoSslCertificateUrls: appSslcertCreateAutoSslCertificateUrls,
			SslKey: readArgFileSupported(appSslcertCreateSslKey),
			SslCrt: readArgFileSupported(appSslcertCreateSslCrt),
			SslCabundle: readArgFileSupported(appSslcertCreateSslCabundle),
			SslForce: appSslcertCreateSslForce,
			AutoUrlLink: appSslcertCreateAutoUrlLink,
		}

		Level27Client.AppSslCertificatesCreate(appID, create)
	},
}

// APP SSLCERT DELETE

var appSslCertDeleteForce bool
var appSslcertDeleteCmd = &cobra.Command{
	Use: "delete [app] [SSL cert]",
	Short: "Delete an SSL certificate from an app",
	Example: "lvl app sslcert delete forum forum.example.com",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		appID := resolveApp(args[0])
		certID := resolveAppSslCertificate(appID, args[1])

		if !appSslCertDeleteForce {
			app := Level27Client.App(appID)
			cert := Level27Client.AppSslCertificatesGetSingle(appID, certID)
			if !confirmPrompt(fmt.Sprintf("Delete SSL certificate %s (%d) on app %s (%d)?", cert.Name, certID, app.Name, appID)) {
				return
			}
		}

		Level27Client.AppSslCertificatesDelete(appID, certID)
	},
}

var appSslCertUpdateCmd = &cobra.Command{
	Use: "update [app] [SSL cert]",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		settings := loadMergeSettings(updateSettingsFile, updateSettings)

		appID := resolveApp(args[0])
		certID := resolveAppSslCertificate(appID, args[1])

		cert := Level27Client.AppSslCertificatesGetSingle(appID, certID)
		put := types.AppSslCertificatePut {
			Name: cert.Name,
			SslType: cert.SslType,
		}

		data := utils.RoundTripJson(put).(map[string]interface{})
		data = mergeMaps(data, settings)

		Level27Client.AppSslCertificatesUpdate(appID, certID, data)
	},
}



var appSslCertFixCmd = &cobra.Command{
	Use: "fix [app] [SSL cert]",
	Short: "Fix an invalid SSL certificate",
	Example: "lvl app sslcert fix forum forum.example.com",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		appID := resolveApp(args[0])
		certID := resolveAppSslCertificate(appID, args[1])

		Level27Client.AppSslCertificatesFix(appID, certID)
	},
}

// APP SSLCERT ACTION

var appSslCertActionCmd = &cobra.Command{
	Use: "action",
}

var appSslCertActionRetryCmd = &cobra.Command{
	Use: "retry [app] [SSL cert]",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		appID := resolveApp(args[0])
		certID := resolveAppSslCertificate(appID, args[1])

		Level27Client.AppSslCertificatesActions(appID, certID, "retry")
	},
}

var appSslCertActionValidateChallengeCmd = &cobra.Command{
	Use: "validateChallenge [app] [SSL cert]",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		appID := resolveApp(args[0])
		certID := resolveAppSslCertificate(appID, args[1])

		Level27Client.AppSslCertificatesActions(appID, certID, "validateChallenge")
	},
}
