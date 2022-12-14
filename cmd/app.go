package cmd

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/level27/l27-go"
	"github.com/level27/lvl/utils"
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
	addWaitFlag(appCreateCmd)
	// flags used for creating app
	flags := appCreateCmd.Flags()
	flags.StringVarP(&appCreateName, "name", "n", "", "Name of the app.")
	flags.StringVarP(&appCreateOrg, "organisation", "", "", "The name of the organisation/owner of the app.")
	flags.Int32SliceVar(&appCreateTeams, "autoTeams", appCreateTeams, "A csv list of team ID's.")
	flags.StringVar(&appCreateExtInfo, "externalInfo", "", "ExternalInfo (required when billableItemInfo entities for an organisation exist in DB.)")
	appCreateCmd.MarkFlagRequired("name")
	appCreateCmd.MarkFlagRequired("organisation")

	// ---- DELETE APP
	appCmd.AddCommand(appDeleteCmd)
	//flag to skip confirmation when deleting an app
	addDeleteConfirmFlag(appCmd)
	addWaitFlag(appDeleteCmd)

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
	addWaitFlag(appComponentCreateCmd)
	appComponentCreateCmd.Flags().StringVarP(&appComponentCreateParamsFile, "params-file", "f", "", "JSON file to read params from. Pass '-' to read from stdin.")
	appComponentCreateCmd.Flags().StringVar(&appComponentCreateName, "name", "", "")
	appComponentCreateCmd.Flags().StringVar(&appComponentCreateType, "type", "", "")
	appComponentCreateCmd.Flags().StringVar(&appComponentCreateSystem, "system", "", "")
	appComponentCreateCmd.Flags().StringVar(&appComponentCreateSystemgroup, "systemgroup", "", "")
	appComponentCreateCmd.Flags().Int32Var(&appComponentCreateSystemprovider, "systemprovider", 0, "")
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
	addDeleteConfirmFlag(AppComponentDeleteCmd)

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
	addWaitFlag(appSslCreateCmd)
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
	addDeleteConfirmFlag(appRestoreDeleteCmd)

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
	addWaitFlag(appComponentUrlCreateCmd)
	appComponentUrlCreateCmd.Flags().BoolVar(&appComponentUrlCreateAuthentication, "authentication", false, "Require HTTP Basic authentication on the URL")
	appComponentUrlCreateCmd.Flags().StringVarP(&appComponentUrlCreateContent, "content", "c", "", "Content for the new URL")
	appComponentUrlCreateCmd.Flags().BoolVar(&appComponentUrlCreateSslForce, "force-ssl", false, "Force usage of SSL on the URL")
	appComponentUrlCreateCmd.Flags().Int32Var(&appComponentUrlCreateSslCertificate, "ssl-certificate", 0, "SSL certificate to use.")
	appComponentUrlCreateCmd.Flags().BoolVar(&appComponentUrlCreateHandleDns, "handle-dns", false, "Automatically create DNS records")
	appComponentUrlCreateCmd.Flags().BoolVar(&appComponentUrlCreateAutoSslCertificate, "auto-ssl-certificate", false, "Automatically create SSL certificate with Let's Encrypt")

	// APP COMPONENT URL DELETE
	appComponentUrlCmd.AddCommand(appComponentUrlDeleteCmd)
	addWaitFlag(appComponentUrlDeleteCmd)
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
func resolveApp(arg string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.AppLookup(arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"app",
		func(app l27.App) string { return fmt.Sprintf("%s (%d)", app.Name, app.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, nil
}

// GET SSL CERTIFICATE ID BASED ON ID
func resolveAppSslCertificate(appID l27.IntID, arg string) (l27.IntID, error) {
	// if arg already int, this is the ID
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.AppSslCertificatesLookup(appID, arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"app SSL certificate",
		func(app l27.AppSslCertificate) string { return fmt.Sprintf("%s (%d)", app.Name, app.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, nil
}

// GET AN APPCOMPONENT ID BASED ON THE NAME
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

// GET AN APP URL ID BASED ON THE NAME
func resolveAppComponentUrl(appID l27.IntID, appComponentID l27.IntID, arg string) (l27.IntID, error) {
	// if arg already int, this is the ID
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.AppComponentUrlLookup(appID, appComponentID, arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"url",
		func(url l27.AppComponentUrlShort) string { return fmt.Sprintf("%s (%d)", url.Content, url.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, err
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
	RunE: func(cmd *cobra.Command, args []string) error {
		ids, err := convertStringsToIDs(args)
		if err != nil {
			return err
		}

		apps, err := getApps(ids)
		if err != nil {
			return err
		}

		outputFormatTable(
			apps,
			[]string{"ID", "NAME", "STATUS"},
			[]string{"ID", "Name", "Status"})

		return nil
	},
}

func getApps(ids []l27.IntID) ([]l27.App, error) {
	c := Level27Client
	if len(ids) == 0 {
		return c.Apps(optGetParameters)
	} else {
		apps := make([]l27.App, len(ids))
		for idx, id := range ids {
			var err error
			apps[idx], err = c.App(id)
			if err != nil {
				return nil, err
			}
		}

		return apps, nil
	}
}

// ---- CREATE NEW APP
var appCreateName, appCreateOrg, appCreateExtInfo string
var appCreateTeams []l27.IntID
var appCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new app.",
	Example: "lvl app create -n myNewApp --organisation level27",
	RunE: func(cmd *cobra.Command, args []string) error {
		// check if name is valid.
		if appCreateName == "" {
			return errors.New("app name cannot be empty")
		}

		// fill in all the props needed for the post request
		organisation, err := resolveOrganisation(appCreateOrg)
		if err != nil {
			return err
		}

		request := l27.AppPostRequest{
			Name:         appCreateName,
			Organisation: organisation,
			AutoTeams:    appCreateTeams,
			ExternalInfo: appCreateExtInfo,
		}

		// when succesfully creating app. app will be returned
		app, err := Level27Client.AppCreate(request)
		if err != nil {
			return err
		}

		if optWait {
			// I'm fairly certain creating apps always completes instantly,
			// but for consistency's sake I'll add the parameter anyways.
			app, err = waitForStatus(
				func() (l27.App, error) { return Level27Client.App(app.ID) },
				func(s l27.App) string { return s.Status },
				"ok",
				[]string{"to_create", "creating"},
			)

			if err != nil {
				return fmt.Errorf("waiting on app status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(app, "templates/entities/app/create.tmpl")
		return nil
	},
}

// ---- DELETE AN APP
var appDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete an app",
	Example: "lvl app delete NameOfMyApp",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// try to find appID based on name
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		if !optDeleteConfirmed {
			app, err := Level27Client.App(appID)
			if err != nil {
				return err
			}

			if !confirmPrompt(fmt.Sprintf("Delete app %s (%d)?", app.Name, app.ID)) {
				return nil
			}
		}

		err = Level27Client.AppDelete(appID)
		if err != nil {
			return err
		}

		if optWait {
			err = waitForDelete(
				func() (l27.App, error) { return Level27Client.App(appID) },
				func(a l27.App) string { return a.Status },
				[]string{"deleting"},
			)

			if err != nil {
				return fmt.Errorf("waiting on app status failed: %s", err.Error())
			}
		}

		return nil
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
	RunE: func(cmd *cobra.Command, args []string) error {
		//check if appID is valid
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		//get the current data from the app. if not changed its needed for put request
		currentData, err := Level27Client.App(appID)
		if err != nil {
			return err
		}

		var currentTeamIDs []string
		for _, team := range currentData.Teams {
			currentTeamIDs = append(currentTeamIDs, fmt.Sprint(team.ID))
		}
		// fill in request with the current data.
		request := l27.AppPutRequest{
			Name:         currentData.Name,
			Organisation: currentData.Organisation.ID,
			AutoTeams:    currentTeamIDs,
		}

		//when flags have been set. we need the currentdata to be updated.
		if cmd.Flag("name").Changed {
			request.Name = appUpdateName
		}

		if cmd.Flag("organisation").Changed {
			organisationID, err := resolveOrganisation(appUpdateOrg)
			if err != nil {
				return err
			}
			request.Organisation = organisationID
		}

		if cmd.Flag("autoTeams").Changed {
			request.AutoTeams = appUpdateTeams
		}

		Level27Client.AppUpdate(appID, request)
		return nil
	},
}

// ---- DESCRIBE APP
var AppDescribeCmd = &cobra.Command{
	Use:     "describe",
	Short:   "Get detailed info about an app.",
	Example: "lvl app describe 2077",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		//check for valid appID
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		// get all data from app by appID
		app, err := Level27Client.App(appID)
		if err != nil {
			return err
		}

		outputFormatTemplate(app, "templates/app.tmpl")
		return nil
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
	RunE: func(cmd *cobra.Command, args []string) error {
		// check for valid appID
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		Level27Client.AppAction(appID, "activate")
		return nil
	},
}

// ---- DEACTIVATE APP
var AppActionDeactivateCmd = &cobra.Command{
	Use:     "deactivate",
	Short:   "Deactivate an app",
	Example: "lvl app action deactivate 2077",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// check for valid appID
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		Level27Client.AppAction(appID, "deactivate")
		return nil
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
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		certs, err := resolveGets(
			// First arg is app ID.
			args[1:],
			func(name string) ([]l27.AppSslCertificate, error) {
				return Level27Client.AppSslCertificatesLookup(appID, name)
			},
			func(certID l27.IntID) (l27.AppSslCertificate, error) {
				return Level27Client.AppSslCertificatesGetSingle(appID, certID)
			},
			func(get l27.CommonGetParams) ([]l27.AppSslCertificate, error) {
				return Level27Client.AppSslCertificatesGetList(appID, appSslGetType, appSslGetStatus, get)
			},
		)

		if err != nil {
			return err
		}

		outputFormatTableFuncs(
			certs,
			[]string{"ID", "Name", "Type", "Status", "SSL Status", "Expiry Date"},
			[]interface{}{"ID", "Name", "SslType", "Status", "SslStatus", "DtExpires", func(c l27.AppSslCertificate) string { return utils.FormatUnixTime(c.DtExpires) }})

		return nil
	},
}

// APP SSL DESCRIBE

var appSslDescribeCmd = &cobra.Command{
	Use:     "describe [app] [SSL cert]",
	Short:   "Get detailed information of an SSL certificate",
	Example: "lvl app ssl describe forum forum.example.com",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		certID, err := resolveAppSslCertificate(appID, args[1])
		if err != nil {
			return err
		}

		cert, err := Level27Client.AppSslCertificatesGetSingle(appID, certID)
		if err != nil {
			return err
		}

		outputFormatTemplate(cert, "templates/appSslCertificate.tmpl")
		return nil
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
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		create := l27.AppSslCertificateCreate{
			Name:                   appSslCreateName,
			SslType:                appSslCreateSslType,
			AutoSslCertificateUrls: appSslCreateAutoSslCertificateUrls,
			SslForce:               appSslCreateSslForce,
			AutoUrlLink:            appSslCreateAutoUrlLink,
		}

		var certificate l27.AppSslCertificate

		switch appSslCreateSslType {
		case "own":
			sslKey, err := readArgFileSupported(appSslCreateSslKey)
			if err != nil {
				return err
			}

			sslCrt, err := readArgFileSupported(appSslCreateSslCrt)
			if err != nil {
				return err
			}

			sslCabundle, err := readArgFileSupported(appSslCreateSslCabundle)
			if err != nil {
				return err
			}

			createOwn := l27.AppSslCertificateCreateOwn{
				AppSslCertificateCreate: create,
				SslKey:                  sslKey,
				SslCrt:                  sslCrt,
				SslCabundle:             sslCabundle,
			}

			certificate, err = Level27Client.AppSslCertificatesCreateOwn(appID, createOwn)
			if err != nil {
				return err
			}

		case "letsencrypt", "xolphin":
			certificate, err = Level27Client.AppSslCertificatesCreate(appID, create)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("invalid SSL type: %s", appSslCreateSslType)
		}

		if optWait {
			certificate, err = waitForStatus(
				func() (l27.AppSslCertificate, error) {
					return Level27Client.AppSslCertificatesGetSingle(appID, certificate.ID)
				},
				func(s l27.AppSslCertificate) string { return s.Status },
				"ok",
				[]string{"to_create", "creating"},
			)

			if err != nil {
				return fmt.Errorf("waiting on certificate status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(certificate, "templates/entities/appSslCertificate/create.tmpl")
		return nil
	},
}

// APP SSL DELETE

var appSslDeleteForce bool
var appSslDeleteCmd = &cobra.Command{
	Use:     "delete [app] [SSL cert]",
	Short:   "Delete an SSL certificate from an app",
	Example: "lvl app ssl delete forum forum.example.com",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		certID, err := resolveAppSslCertificate(appID, args[1])
		if err != nil {
			return err
		}

		if !appSslDeleteForce {
			app, err := Level27Client.App(appID)
			if err != nil {
				return err
			}

			cert, err := Level27Client.AppSslCertificatesGetSingle(appID, certID)
			if err != nil {
				return err
			}

			if !confirmPrompt(fmt.Sprintf("Delete SSL certificate %s (%d) on app %s (%d)?", cert.Name, certID, app.Name, appID)) {
				return nil
			}
		}

		Level27Client.AppSslCertificatesDelete(appID, certID)
		return nil
	},
}

// APP SSL UPDATE
var appSslUpdateCmd = &cobra.Command{
	Use: "update [app] [SSL cert]",

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

		certID, err := resolveAppSslCertificate(appID, args[1])
		if err != nil {
			return err
		}

		cert, err := Level27Client.AppSslCertificatesGetSingle(appID, certID)
		if err != nil {
			return err
		}

		put := l27.AppSslCertificatePut{
			Name:    cert.Name,
			SslType: cert.SslType,
		}

		data := utils.RoundTripJson(put).(map[string]interface{})
		data = mergeMaps(data, settings)

		Level27Client.AppSslCertificatesUpdate(appID, certID, data)
		return nil
	},
}

// APP SSL FIX
var appSslFixCmd = &cobra.Command{
	Use:     "fix [app] [SSL cert]",
	Short:   "Fix an invalid SSL certificate",
	Example: "lvl app ssl fix forum forum.example.com",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		certID, err := resolveAppSslCertificate(appID, args[1])
		if err != nil {
			return err
		}

		Level27Client.AppSslCertificatesFix(appID, certID)
		return nil
	},
}

// APP SSL ACTION

var appSslActionCmd = &cobra.Command{
	Use: "action",
}

var appSslActionRetryCmd = &cobra.Command{
	Use: "retry [app] [SSL cert]",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		certID, err := resolveAppSslCertificate(appID, args[1])
		if err != nil {
			return err
		}

		Level27Client.AppSslCertificatesActions(appID, certID, "retry")
		return nil
	},
}

var appSslActionValidateChallengeCmd = &cobra.Command{
	Use: "validateChallenge [app] [SSL cert]",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		certID, err := resolveAppSslCertificate(appID, args[1])
		if err != nil {
			return err
		}

		Level27Client.AppSslCertificatesActions(appID, certID, "validateChallenge")
		return nil
	},
}

// APP SSL KEY
var appSslKeyCmd = &cobra.Command{
	Use:     "key",
	Short:   "Return a private key for type 'own' sslCertificate.",
	Example: "lvl app ssl key MyAppName",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		certID, err := resolveAppSslCertificate(appID, args[1])
		if err != nil {
			return err
		}

		key, err := Level27Client.AppSslCertificatesKey(appID, certID)
		if err != nil {
			return err
		}

		fmt.Print(key.SslKey)
		return nil
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
	RunE: func(cmd *cobra.Command, args []string) error {
		//search for appID based on Appname
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		ids, err := convertStringsToIDs(args[1:])
		if err != nil {
			return errors.New("nvalid component ID")
		}

		res, err := getComponents(appID, ids)
		if err != nil {
			return err
		}

		outputFormatTable(
			res,
			[]string{"ID", "NAME", "STATUS"},
			[]string{"ID", "Name", "Status"})

		return nil
	},
}

func getComponents(appID l27.IntID, ids []l27.IntID) ([]l27.AppComponent, error) {
	c := Level27Client
	if len(ids) == 0 {
		return c.AppComponentsGet(appID, optGetParameters)
	} else {
		components := make([]l27.AppComponent, len(ids))
		for idx, id := range ids {
			var err error
			components[idx], err = c.AppComponentGetSingle(appID, id)
			if err != nil {
				return nil, err
			}
		}

		return components, nil
	}
}

// ---- CREATE COMPONENT
var appComponentCreateParamsFile string
var appComponentCreateName string
var appComponentCreateType string
var appComponentCreateSystem string
var appComponentCreateSystemgroup string
var appComponentCreateSystemprovider l27.IntID
var appComponentCreateParams *[]string
var appComponentCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new appcomponent.",
	Example: "lvl app component create --name myComponentName --type docker",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if appComponentCreateSystem == "" && appComponentCreateSystemgroup == "" {
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
		create["category"] = "config"
		create["appcomponenttype"] = appComponentCreateType

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
		data["category"] = appComponent.Category
		data["systemgroup"] = appComponent.Systemgroup
		data["system"] = nil
		if appComponent.Systemgroup == nil {
			data["system"] = appComponent.Systems[0].ID
		}

		for k, v := range appComponent.Appcomponentparameters {
			param, ok := parameterTypes[k]
			if !ok {
				// Maybe should be a panic instead?
				return fmt.Errorf("API returned unknown parameter in component data: '%s'", k)
			}

			paramType := parameterTypes[k].Type

			if param.Readonly || param.DisableEdit {
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

	if appComponentType.Servicetype.SSHKeyPossible {
		appComponentType.Servicetype.Parameters = append(
			appComponentType.Servicetype.Parameters,
			l27.AppComponentTypeParameter{
				Name:           "sshkeys",
				DisplayName:    "SSH Keys",
				Description:    "The SSH keys that can be used to log into the component",
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

// ---- DELETE COMPONENT
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
			appComponent, err := Level27Client.AppComponentGetSingle(appID, appComponentID)
			if err != nil {
				return err
			}

			if !confirmPrompt(fmt.Sprintf("Delete app component %s (%d) on app %s (%d)?", appComponent.Name, appComponent.ID, appComponent.App.Name, appComponent.App.ID)) {
				return nil
			}
		}

		Level27Client.AppComponentsDelete(appID, appComponentID)
		return nil
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

// ---- (TYPES) GET LIST OF APPCOMPONENT TYPES
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

// ---- (PARAMETERS) GET LIST OF PARAMETERS FOR A SPECIFIC APPCOMPONENT TYPE
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
	RunE: func(cmd *cobra.Command, args []string) error {
		// search for appID based on name
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		Restores, err := Level27Client.AppComponentRestoresGet(appID)
		if err != nil {
			return err
		}

		outputFormatTableFuncs(Restores,
			[]string{"ID", "FILENAME", "STATUS", "DATE", "APPCOMPONENT_ID", "APPCOMPONENT_NAME"},
			[]interface{}{"ID", "Filename", "Status", func(r l27.AppComponentRestore) string { return utils.FormatUnixTime(r.AvailableBackup.Date) }, "Appcomponent.ID", "Appcomponent.Name"})

		return nil
	},
}

// ---- CREATE A NEW RESTORE
var appComponentRestoreCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new restore for an app.",
	Example: "lvl app restore create MyAppName MyComponentName 453",
	Args:    cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		//search appID based on appname
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		// search componentID based on name
		componentID, err := resolveAppComponent(appID, args[1])
		if err != nil {
			return err
		}

		backupID, err := checkSingleIntID(args[2], "backup")
		if err != nil {
			return err
		}

		request := l27.AppComponentRestoreRequest{
			Appcomponent:    componentID,
			AvailableBackup: backupID,
		}

		restore, err := Level27Client.AppComponentRestoreCreate(appID, request)
		if err != nil {
			return err
		}

		log.Printf("Restore created. [ID: %v].", restore.ID)
		return nil
	},
}

// ---- DELETE A RESTORE
var appRestoreDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a specific restore from an app.",
	Example: "lvl app component restore delete MyAppName 4532",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// search for appID based on name
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		// check if restoreID is valid type
		restoreID, err := checkSingleIntID(args[1], "restore")
		if err != nil {
			return err
		}

		if !optDeleteConfirmed {
			app, err := Level27Client.App(appID)
			if err != nil {
				return err
			}

			if !confirmPrompt(fmt.Sprintf("Delete restore %d on app %s (%d)?", restoreID, app.Name, app.ID)) {
				return nil
			}
		}

		Level27Client.AppComponentRestoresDelete(appID, restoreID)
		return nil
	},
}

// ---- DOWNLOAD RESTORE FILE
var appComponentRestoreDownloadName string
var appComponentRestoreDownloadCmd = &cobra.Command{
	Use:     "download [appname] [restoreID]",
	Short:   "Download the restore file.",
	Example: "lvl app component restore download MyAppName 4123",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// search appID based on name
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		// check if restoreID is valid type
		restoreID, err := checkSingleIntID(args[1], "Restore")
		if err != nil {
			return err
		}

		Level27Client.AppComponentRestoreDownload(appID, restoreID, appComponentRestoreDownloadName)
		return nil
	},
}

//-------------------------------------------------  APP COMPONENT BACKUPS (GET) -------------------------------------------------
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
	RunE: func(cmd *cobra.Command, args []string) error {
		// search appID based on appname
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		// search componentID based on name
		componentID, err := resolveAppComponent(appID, args[1])
		if err != nil {
			return err
		}

		availableBackups, err := Level27Client.AppComponentbackupsGet(appID, componentID)
		if err != nil {
			return err
		}

		outputFormatTableFuncs(availableBackups,
			[]string{"ID", "SNAPSHOTNAME", "DATE"},
			[]interface{}{"ID", "SnapshotName", func(a l27.AppComponentAvailableBackup) string {
				return utils.FormatUnixTime(a.Date)
			}})

		return nil
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
	RunE: func(cmd *cobra.Command, args []string) error {
		//search for appID based on name
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		migrations, err := Level27Client.AppMigrationsGet(appID)
		if err != nil {
			return err
		}

		outputFormatTableFuncs(migrations,
			[]string{"ID", "MIGRATION_TYPE", "STATUS", "DATE_PLANNED"},
			[]interface{}{"ID", "MigrationType", "Status", func(m l27.AppMigration) string {
				return utils.FormatUnixTime(m.DtPlanned)
			}})

		return nil
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

	RunE: func(cmd *cobra.Command, args []string) error {
		//search for appid based on appName
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		items := []l27.AppMigrationItem{}

		for _, migrationItem := range appMigrationCreateItems {
			res, err := ParseMigrationItem(appID, migrationItem)
			if err != nil {
				return err
			}

			items = append(items, res)
		}

		request := l27.AppMigrationRequest{
			MigrationType:      "automatic",
			DtPlanned:          appMigrationCreatePlanned,
			MigrationItemArray: items,
		}

		migration, err := Level27Client.AppMigrationsCreate(appID, request)
		if err != nil {
			return err
		}

		log.Printf("migration created! [ID: '%v']", migration.ID)
		return nil
	},
}

func ParseMigrationItem(appID l27.IntID, values string) (l27.AppMigrationItem, error) {
	valueSplitted := strings.Split(values, ",")

	item := l27.AppMigrationItem{
		Ord:    1,
		Source: "cp4",
	}

	haveAnyDst := false
	haveAnySrc := false
	for _, keyValuePair := range valueSplitted {
		// Go over key value pairs and fill out the migration item as we go.

		key, value, err := ParseMigrationItemKeyValuePair(keyValuePair)
		if err != nil {
			return l27.AppMigrationItem{}, err
		}

		switch key {
		case "ord":
			val, err := strconv.Atoi(value)
			if err != nil {
				return l27.AppMigrationItem{}, err
			}
			item.Ord = int32(val)

		case "destSystem":
			item.DestinationEntityID, err = resolveSystem(value)
			item.DestinationEntity = "system"
			haveAnyDst = true
			if err != nil {
				return l27.AppMigrationItem{}, err
			}

		case "destGroup":
			item.DestinationEntityID, err = resolveSystemgroup(value)
			item.DestinationEntity = "systemgroup"
			haveAnyDst = true
			if err != nil {
				return l27.AppMigrationItem{}, err
			}

		case "source":
			appComponent, err := resolveAppComponent(appID, value)
			if err != nil {
				return l27.AppMigrationItem{}, err
			}

			appComponentType, err := Level27Client.AppComponentGetSingle(appID, appComponent)
			if err != nil {
				return l27.AppMigrationItem{}, err
			}

			haveAnySrc = true

			item.SourceInfo = appComponent
			item.Type = appComponentType.Appcomponenttype

		default:
			return l27.AppMigrationItem{}, fmt.Errorf("unknown property in migration item: %s", key)
		}
	}

	if !haveAnyDst {
		return l27.AppMigrationItem{}, errors.New("no destination specified for migration item")
	}

	if !haveAnySrc {
		return l27.AppMigrationItem{}, errors.New("no source specified for migration item")
	}

	return item, nil
}

func ParseMigrationItemKeyValuePair(keyValuePair string) (string, string, error) {
	split := strings.SplitN(keyValuePair, "=", 2)
	if len(split) == 1 {
		return "", "", fmt.Errorf("migrationItem property not defined correctly: '%v'. Use '=' to define properties", keyValuePair)
	}

	key := strings.TrimSpace(split[0])
	value := strings.TrimSpace(split[1])

	return key, value, nil
}

// ---- UPDATE MIGRATION
var appMigrationsUpdateType, appMigrationsUpdateDtPlanned string
var appMigrationsUpdateCmd = &cobra.Command{
	Use:     "update [appID] [migrationID]",
	Short:   "Update an app migration.",
	Example: "lvl app migrations update MyAppName 3414",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		//search for appID based on name
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		// check for valid migrationID type
		migrationID, err := checkSingleIntID(args[1], "appMigration")
		if err != nil {
			return err
		}

		request := l27.AppMigrationUpdate{
			MigrationType: appMigrationsUpdateType,
			DtPlanned:     appMigrationsUpdateDtPlanned,
		}

		err = Level27Client.AppMigrationsUpdate(appID, migrationID, request)
		if err != nil {
			return err
		}

		log.Print("migration succesfully updated!")
		return nil
	},
}

// ---- DESCRIBE MIGRATION
var appMigrationDescribeCmd = &cobra.Command{
	Use:     "describe [appID] [migrationID]",
	Short:   "Get detailed info about a specific migration.",
	Example: "lvl app migrations describe MyAppName 1243",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// search for appID based on name
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		// check for valid migrationID type
		migrationID, err := checkSingleIntID(args[1], "appMigration")
		if err != nil {
			return err
		}

		migration, err := Level27Client.AppMigrationDescribe(appID, migrationID)
		if err != nil {
			return err
		}

		outputFormatTemplate(migration, "templates/appMigration.tmpl")
		return nil
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
	RunE: func(cmd *cobra.Command, args []string) error {
		// search for appID based on name
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		// check for valid migrationID type
		migrationID, err := checkSingleIntID(args[1], "appMigration")
		if err != nil {
			return err
		}

		Level27Client.AppMigrationsAction(appID, migrationID, "confirm")
		return nil
	},
}

// ---- DENY MIGRATION
var appMigrationsActionDenyCmd = &cobra.Command{
	Use:     "deny",
	Short:   "Execute confirm action on a migration",
	Example: "lvl app migrations action deny MyAppName 332",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// search for appID based on name
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		// check for valid migrationID type
		migrationID, err := checkSingleIntID(args[1], "appMigration")
		if err != nil {
			return err
		}

		Level27Client.AppMigrationsAction(appID, migrationID, "deny")
		return nil
	},
}

// ---- RETRY MIGRATION
var appMigrationsActionRetryCmd = &cobra.Command{
	Use:     "retry",
	Short:   "Execute confirm action on a migration",
	Example: "lvl app migrations action retry MyAppName 332",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// search for appID based on name
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		// check for valid migrationID type
		migrationID, err := checkSingleIntID(args[1], "appMigration")
		if err != nil {
			return err
		}

		Level27Client.AppMigrationsAction(appID, migrationID, "retry")
		return nil
	},
}

// ------- APP COMPONENT URLs

// APP COMPONENT URL
var appComponentUrlCmd = &cobra.Command{
	Use:     "url",
	Aliases: []string{"urls"},
}

// APP COMPONENT URL GET
var appComponentUrlGetCmd = &cobra.Command{
	Use: "get",

	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		componentID, err := resolveAppComponent(appID, args[1])
		if err != nil {
			return err
		}

		results, err := resolveGets(
			args[2:],
			func(name string) ([]l27.AppComponentUrlShort, error) {
				return Level27Client.AppComponentUrlLookup(appID, componentID, name)
			},
			func(i l27.IntID) (l27.AppComponentUrlShort, error) {
				res, err := Level27Client.AppComponentUrlGetSingle(appID, componentID, i)
				if err != nil {
					return l27.AppComponentUrlShort{}, err
				}
				return res.ToShort(), nil
			},
			func(cgp l27.CommonGetParams) ([]l27.AppComponentUrlShort, error) {
				return Level27Client.AppComponentUrlGetList(appID, componentID, cgp)
			},
		)

		if err != nil {
			return err
		}

		outputFormatTable(
			results,
			[]string{"ID", "CONTENT", "STATUS", "TYPE", "SSL CERT", "FORCE SSL", "HANDLE DNS", "AUTHENTICATE"},
			[]string{"ID", "Content", "Status", "Type", "SslCertificate.Name", "SslForce", "HandleDNS", "Authentication"})

		return nil
	},
}

// APP COMPONENT URL CREATE
var appComponentUrlCreateAuthentication bool
var appComponentUrlCreateContent string
var appComponentUrlCreateSslForce bool
var appComponentUrlCreateSslCertificate l27.IntID
var appComponentUrlCreateHandleDns bool
var appComponentUrlCreateAutoSslCertificate bool
var appComponentUrlCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an url for an appcomponent.",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		componentID, err := resolveAppComponent(appID, args[1])
		if err != nil {
			return err
		}

		var cert *l27.IntID
		if appComponentUrlCreateSslCertificate == 0 {
			cert = nil
		} else {
			cert = &appComponentUrlCreateSslCertificate
		}

		create := l27.AppComponentUrlCreate{
			Authentication:     appComponentUrlCreateAuthentication,
			Content:            appComponentUrlCreateContent,
			SslForce:           appComponentUrlCreateSslForce,
			SslCertificate:     cert,
			HandleDns:          appComponentUrlCreateHandleDns,
			AutoSslCertificate: appComponentUrlCreateAutoSslCertificate,
		}

		url, err := Level27Client.AppComponentUrlCreate(appID, componentID, create)
		if err != nil {
			return err
		}

		if optWait {
			url, err = waitForStatus(
				func() (l27.AppComponentUrl, error) {
					return Level27Client.AppComponentUrlGetSingle(appID, componentID, url.ID)
				},
				func(s l27.AppComponentUrl) string { return s.Status },
				"ok",
				[]string{"creating"},
			)

			if err != nil {
				return fmt.Errorf("waiting on URL status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(url, "templates/entities/appComponentUrl/create.tmpl")
		return nil
	},
}

// APP COMPONENT URL DELETE
var appComponentUrlDeleteForce bool
var appComponentUrlDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an url from an appcomponent.",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		componentID, err := resolveAppComponent(appID, args[1])
		if err != nil {
			return err
		}

		urlID, err := resolveAppComponentUrl(appID, componentID, args[2])
		if err != nil {
			return err
		}

		if !appComponentUrlDeleteForce {
			url, err := Level27Client.AppComponentUrlGetSingle(appID, componentID, urlID)
			if err != nil {
				return err
			}

			msg := fmt.Sprintf(
				"Delete URL %s (%d) on app comp %s (%d)?",
				url.Content, url.ID,
				url.Appcomponent.Name, url.Appcomponent.ID)

			if !confirmPrompt(msg) {
				return nil
			}
		}

		err = Level27Client.AppComponentUrlDelete(appID, componentID, urlID)
		if err != nil {
			return err
		}

		if optWait {
			err = waitForDelete(
				func() (l27.AppComponentUrl, error) {
					return Level27Client.AppComponentUrlGetSingle(appID, componentID, urlID)
				},
				func(a l27.AppComponentUrl) string { return a.Status },
				[]string{"to_delete"},
			)

			if err != nil {
				return fmt.Errorf("waiting on app status failed: %s", err.Error())
			}
		}

		return nil
	},
}
