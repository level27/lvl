package cmd

import (
	"fmt"

	"github.com/level27/l27-go"
	"github.com/level27/lvl/utils"
	"github.com/spf13/cobra"
)

// Contains backups AND restores since they're so closely related.

func init() {
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

}

// -------------------------------------------------  APP COMPONENT RESTORE (GET / CREATE / UPDATE / DELETE / DOWNLOAD) -------------------------------------------------
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

		outputFormatTemplate(restore, "templates/entities/appComponentRestore/create.tmpl")
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

		err = Level27Client.AppComponentRestoresDelete(appID, restoreID)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/appComponentRestore/delete.tmpl")
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

		err = Level27Client.AppComponentRestoreDownload(appID, restoreID, appComponentRestoreDownloadName)
		if err != nil {
			return err
		}

		return nil
	},
}

// -------------------------------------------------  APP COMPONENT BACKUPS (GET) -------------------------------------------------
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
