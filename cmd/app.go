package cmd

import (
	"errors"
	"fmt"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

func init() {
	// APP
	RootCmd.AddCommand(appCmd)

	// APP GET
	appCmd.AddCommand(appGetCmd)
	addCommonGetFlags(appGetCmd)

	// APP CREATE
	appCmd.AddCommand(appCreateCmd)
	addWaitFlag(appCreateCmd)
	flags := appCreateCmd.Flags()
	flags.StringVarP(&appCreateName, "name", "n", "", "Name of the app.")
	flags.StringVarP(&appCreateOrg, "organisation", "", "", "The name of the organisation/owner of the app.")
	flags.Int32SliceVar(&appCreateTeams, "autoTeams", appCreateTeams, "A csv list of team ID's.")
	flags.StringVar(&appCreateExtInfo, "externalInfo", "", "ExternalInfo (required when billableItemInfo entities for an organisation exist in DB.)")
	appCreateCmd.MarkFlagRequired("name")
	appCreateCmd.MarkFlagRequired("organisation")

	// APP DELETE
	appCmd.AddCommand(appDeleteCmd)
	addDeleteConfirmFlag(appCmd)
	addWaitFlag(appDeleteCmd)

	// APP UPDATE
	appCmd.AddCommand(appUpdateCmd)
	flags = appUpdateCmd.Flags()
	flags.StringVarP(&appUpdateName, "name", "n", "", "Name of the app.")
	flags.StringVarP(&appUpdateOrg, "organisation", "", "", "The name of the organisation/owner of the app.")
	flags.StringSliceVar(&appUpdateTeams, "autoTeams", appUpdateTeams, "A csv list of team ID's.")

	// APP DESCRIBE
	appCmd.AddCommand(AppDescribeCmd)

	// APP ACTION
	appCmd.AddCommand(AppActionCmd)

	// APP ACTION ACTIVATE
	AppActionCmd.AddCommand(AppActionActivateCmd)

	// APP ACTION DEACTIVATE
	AppActionCmd.AddCommand(AppActionDeactivateCmd)

	// APP INTEGRITY
	addIntegrityCheckCmds(appCmd, "apps", resolveApp)

	// APP ACCESS
	addAccessCmds(appCmd, "apps", resolveApp)

	// APP BILLING
	addBillingCmds(appCmd, "apps", resolveApp)
}

// Resolve the ID of an app based on user-provided name or ID.
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

// APP
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
		apps, err := resolveGets(
			args,
			Level27Client.AppLookup,
			Level27Client.App,
			Level27Client.Apps)

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

		outputFormatTemplate(nil, "templates/entities/app/delete.tmpl")
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

		err = Level27Client.AppUpdate(appID, request)
		if err != nil {
			return err
		}

		outputFormatTemplate(request, "templates/entities/app/update.tmpl")

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

		err = Level27Client.AppAction(appID, "activate")
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/app/activate.tmpl")

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

		err = Level27Client.AppAction(appID, "deactivate")
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/app/deactivate.tmpl")
		return nil
	},
}

// ------------------------------------------------- APP COMPONENTS (CREATE / GET / UPDATE / DELETE / DESCRIBE)-------------------------------------------------
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

		components, err := resolveGets(
			// First arg is app ID.
			args[1:],
			func(name string) ([]l27.AppComponent, error) {
				return Level27Client.AppComponentLookup(appID, name)
			},
			func(certID l27.IntID) (l27.AppComponent, error) {
				return Level27Client.AppComponentGetSingle(appID, certID)
			},
			func(get l27.CommonGetParams) ([]l27.AppComponent, error) {
				return Level27Client.AppComponentsGet(appID, get)
			})

		if err != nil {
			return err
		}

		outputFormatTable(
			components,
			[]string{"ID", "NAME", "STATUS"},
			[]string{"ID", "Name", "Status"})

		return nil
	},
}
