package cmd

import (
	"fmt"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

func init() {
	// APP COMPONENT CRON
	appComponentCmd.AddCommand(appComponentCronCmd)

	// APP COMPONENT CRON GET
	appComponentCronCmd.AddCommand(appComponentCronGetCmd)

	// APP COMPONENT CRON CREATE
	appComponentCronCmd.AddCommand(appComponentCronCreateCmd)
	addWaitFlag(appComponentCronCreateCmd)
	appComponentCronCreateCmd.Flags().StringVarP(&optAppComponentCronCreateName, "name", "n", "", "The name of the new cron")
	appComponentCronCreateCmd.Flags().StringVarP(&optAppComponentCronCreateSchedule, "schedule", "s", "", "The schedule that controls when the cron should fire. Takes the form of a standard POSIX crontab pattern.")
	appComponentCronCreateCmd.Flags().StringVarP(&optAppComponentCronCreateCommand, "command", "c", "", "The shell command that will be executed when the cron fires.")
	appComponentCronCreateCmd.MarkFlagRequired("name")
	appComponentCronCreateCmd.MarkFlagRequired("schedule")
	appComponentCronCreateCmd.MarkFlagRequired("command")

	// APP COMPONENT CRON UPDATE
	appComponentCronCmd.AddCommand(appComponentCronUpdateCmd)
	addWaitFlag(appComponentCronUpdateCmd)
	appComponentCronUpdateCmd.Flags().StringVarP(&optAppComponentCronUpdateName, "name", "n", "", "The name of the new cron")
	appComponentCronUpdateCmd.Flags().StringVarP(&optAppComponentCronUpdateSchedule, "schedule", "s", "", "The schedule that controls when the cron should fire. Takes the form of a standard POSIX crontab pattern.")
	appComponentCronUpdateCmd.Flags().StringVarP(&optAppComponentCronUpdateCommand, "command", "c", "", "The shell command that will be executed when the cron fires.")

	// APP COMPONENT CRON DELETE
	appComponentCronCmd.AddCommand(appComponentCronDeleteCmd)
	addWaitFlag(appComponentCronDeleteCmd)
}

// Resolve the ID of an app component cron based on user-provided name or ID.
func resolveAppComponentCron(appID l27.IntID, appComponentID l27.IntID, arg string) (l27.IntID, error) {
	// if arg already int, this is the ID
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.AppComponentCronLookup(appID, appComponentID, arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"cron",
		func(cron l27.AppComponentCronShort) string { return fmt.Sprintf("%s (%d)", cron.Name, cron.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, err
}

// APP COMPONENT CRON
var appComponentCronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Commands for managing crons",
	Long:  "Crons are scheduled automatic commands that can be ran on app components",
}

// APP COMPONENT CRON CREATE
var optAppComponentCronCreateName string
var optAppComponentCronCreateSchedule string
var optAppComponentCronCreateCommand string
var appComponentCronCreateCmd = &cobra.Command{
	Use:   "create <app> <component> -n <name> -s <scedule> -c <command>",
	Short: "Create a new cron on an app component",
	Example: `Run a script at 01:17 every day:
lvl app component cron create my-app php -n prune_database -s "17 01 * * *" -c ./prune_database.sh`,

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

		create := l27.AppComponentCronCreate{
			Name:     optAppComponentCronCreateName,
			Schedule: optAppComponentCronCreateSchedule,
			Command:  optAppComponentCronCreateCommand,
		}

		cron, err := Level27Client.AppComponentCronCreate(appID, componentID, create)
		if err != nil {
			return err
		}

		if optWait {
			cron, err = waitForStatus(
				func() (l27.AppComponentCron, error) {
					return Level27Client.AppComponentCronGetSingle(appID, componentID, cron.ID)
				},
				func(s l27.AppComponentCron) string { return s.Status },
				"ok",
				[]string{"to_create", "creating"},
			)

			if err != nil {
				return fmt.Errorf("waiting on cron status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(cron, "templates/entities/appComponentCron/create.tmpl")
		return nil
	},
}

// APP COMPONENT CRON GET
var appComponentCronGetCmd = &cobra.Command{
	Use:   "get <app> <component>",
	Short: "Get a list of crons on an app component",

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

		results, err := resolveGets(
			args[2:],
			func(name string) ([]l27.AppComponentCronShort, error) {
				return Level27Client.AppComponentCronLookup(appID, componentID, name)
			},
			func(i l27.IntID) (l27.AppComponentCronShort, error) {
				res, err := Level27Client.AppComponentCronGetSingle(appID, componentID, i)
				if err != nil {
					return l27.AppComponentCronShort{}, err
				}
				return res.ToShort(), nil
			},
			func(cgp l27.CommonGetParams) ([]l27.AppComponentCronShort, error) {
				return Level27Client.AppComponentCronGetList(appID, componentID, cgp)
			},
		)

		if err != nil {
			return err
		}

		outputFormatTable(
			results,
			[]string{"ID", "NAME", "STATUS", "SCHEDULE", "COMMAND"},
			[]string{"ID", "Name", "Status", "Schedule", "Command"})

		return nil
	},
}

// APP COMPONENT CRON UPDATE
var optAppComponentCronUpdateName string
var optAppComponentCronUpdateSchedule string
var optAppComponentCronUpdateCommand string
var appComponentCronUpdateCmd = &cobra.Command{
	Use:   "update <app> <component> <cron> [-n <name>] [-s <scedule>] [-c <command>]",
	Short: "Update an existing cron on an app component",

	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		componentID, err := resolveAppComponent(appID, args[1])
		if err != nil {
			return err
		}

		cronID, err := resolveAppComponentCron(appID, componentID, args[2])
		if err != nil {
			return err
		}

		cron, err := Level27Client.AppComponentCronGetSingle(appID, componentID, cronID)
		if err != nil {
			return err
		}

		update := l27.AppComponentCronUpdate{
			Name:     cron.Name,
			Schedule: cron.Schedule,
			Command:  cron.Command,
		}

		if optAppComponentCronUpdateName != "" {
			update.Name = optAppComponentCronUpdateName
		}

		if optAppComponentCronUpdateSchedule != "" {
			update.Schedule = optAppComponentCronUpdateSchedule
		}

		if optAppComponentCronUpdateCommand != "" {
			update.Command = optAppComponentCronUpdateCommand
		}

		err = Level27Client.AppComponentCronUpdate(appID, componentID, cronID, update)
		if err != nil {
			return err
		}

		if optWait {
			cron, err = waitForStatus(
				func() (l27.AppComponentCron, error) {
					return Level27Client.AppComponentCronGetSingle(appID, componentID, cron.ID)
				},
				func(s l27.AppComponentCron) string { return s.Status },
				"ok",
				[]string{"updating"},
			)

			if err != nil {
				return fmt.Errorf("waiting on cron status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(nil, "templates/entities/appComponentCron/update.tmpl")
		return nil
	},
}

var appComponentCronDeleteCmd = &cobra.Command{
	Use:   "delete <app> <component> <cron>",
	Short: "Delete a cron on an app component",

	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		componentID, err := resolveAppComponent(appID, args[1])
		if err != nil {
			return err
		}

		cronID, err := resolveAppComponentCron(appID, componentID, args[2])
		if err != nil {
			return err
		}

		if !appComponentUrlDeleteForce {
			cron, err := Level27Client.AppComponentCronGetSingle(appID, componentID, cronID)
			if err != nil {
				return err
			}

			msg := fmt.Sprintf(
				"Delete cron %s (%d) on app comp %s (%d)?",
				cron.Name, cron.ID,
				cron.Appcomponent.Name, cron.Appcomponent.ID)

			if !confirmPrompt(msg) {
				return nil
			}
		}

		err = Level27Client.AppComponentCronDelete(appID, componentID, cronID)
		if err != nil {
			return err
		}

		if optWait {
			err = waitForDelete(
				func() (l27.AppComponentCron, error) {
					return Level27Client.AppComponentCronGetSingle(appID, componentID, cronID)
				},
				func(a l27.AppComponentCron) string { return a.Status },
				[]string{"to_delete"},
			)

			if err != nil {
				return fmt.Errorf("waiting on cron status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(nil, "templates/entities/appComponentCron/delete.tmpl")
		return nil
	},
}
