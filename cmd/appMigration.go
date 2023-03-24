package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/level27/l27-go"
	"github.com/level27/lvl/utils"
	"github.com/spf13/cobra"
)

func init() {
	//-------------------------------------------------  APP MIGRATIONS (GET / DESCRIBE / CREATE / UPDATE) -------------------------------------------------
	// ---- MIGRATIONS COMMAND
	appCmd.AddCommand(appMigrationsCmd)

	// ---- GET LIST OF MIGRATIONS
	appMigrationsCmd.AddCommand(appMigrationsGetCmd)

	// ---- CREATE NEW APP MIGRATION
	appMigrationsCmd.AddCommand(appMigrationsCreateCmd)
	// flags needed to create new migration
	flags := appMigrationsCreateCmd.Flags()
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

// -------------------------------------------------  APP MIGRATIONS (GET / DESCRIBE / CREATE / UPDATE) -------------------------------------------------
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

		outputFormatTemplate(migration, "templates/entities/appMigration/create.tmpl")
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

		outputFormatTemplate(nil, "templates/entities/appMigration/update.tmpl")
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

// -------------------------------------------------  APP MIGRATIONS ACTIONS (CONFIRM / DENY / RESTART) -------------------------------------------------
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

		err = Level27Client.AppMigrationsAction(appID, migrationID, "confirm")
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/appMigration/confirm.tmpl")
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

		err = Level27Client.AppMigrationsAction(appID, migrationID, "deny")
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/appMigration/deny.tmpl")
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

		err = Level27Client.AppMigrationsAction(appID, migrationID, "retry")
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/appMigration/retry.tmpl")
		return nil
	},
}
