package cmd

import (
	"fmt"
	"strings"

	"github.com/level27/l27-go"
	"github.com/level27/lvl/utils"
	"github.com/spf13/cobra"
)

func init() {
	// MAIL FORWARDER
	mailCmd.AddCommand(mailForwarderCmd)

	// MAIL FORWARDER GET
	mailForwarderCmd.AddCommand(mailForwarderGetCmd)
	addCommonGetFlags(mailForwarderGetCmd)

	// MAIL FORWARDER CREATE
	mailForwarderCmd.AddCommand(mailForwarderCreateCmd)
	addWaitFlag(mailForwarderCreateCmd)
	mailForwarderCreateCmd.Flags().StringVar(&mailForwarderCreateAddress, "address", "", "Address of the mail forwarder")
	mailForwarderCreateCmd.Flags().StringVar(&mailForwarderCreateDestination, "destination", "", "Comma-separated list of destination addresses to forward to")

	// MAIL FORWARDER DELETE
	mailForwarderCmd.AddCommand(mailForwarderDeleteCmd)
	addWaitFlag(mailForwarderDeleteCmd)
	mailForwarderDeleteCmd.Flags().BoolVar(&mailForwarderDeleteForce, "force", false, "Do not ask for confirmation to delete the mail forwarder")

	// MAIL FORWARDER UPDATE
	mailForwarderCmd.AddCommand(mailForwarderUpdateCmd)
	settingsFileFlag(mailForwarderUpdateCmd)
	settingString(mailForwarderUpdateCmd, updateSettings, "destination", "Comma-separated list of all destinations for this forwarder")

	// MAIL FORWARDER DESTINATION
	mailForwarderCmd.AddCommand(mailForwarderDestinationCmd)

	// MAIL FORWARDER DESTINATION ADD
	mailForwarderDestinationCmd.AddCommand(mailForwarderDestinationAddCmd)

	// MAIL FORWARDER DESTINATION REMOVE
	mailForwarderDestinationCmd.AddCommand(mailForwarderDestinationRemoveCmd)
}

func resolveMailforwarder(mailgroupID l27.IntID, arg string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.MailgroupsMailforwardersLookup(mailgroupID, arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"mailforwarder",
		func(app l27.Mailforwarder) string { return fmt.Sprintf("%s (%d)", app.Address, app.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, nil
}

// MAIL FORWARDER
var mailForwarderCmd = &cobra.Command{
	Use:   "forwarder",
	Short: "Commands for managing mail forwarders",

	Aliases: []string{"fwd"},
}

// MAIL FORWARDER GET
var mailForwarderGetCmd = &cobra.Command{
	Use:   "get [mailgroup]",
	Short: "Get a list of mail forwarders in a mail group",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mailgroupID, err := resolveMailgroup(args[0])
		if err != nil {
			return err
		}

		mailboxes, err := Level27Client.MailgroupsMailforwardersGetList(mailgroupID, optGetParameters)
		if err != nil {
			return err
		}

		outputFormatTableFuncs(
			mailboxes,
			[]string{"ID", "Status", "Address", "Destimations"},
			[]interface{}{"ID", "Status", "Address", func(f l27.Mailforwarder) string {
				first := true
				result := ""
				for _, dest := range f.Destination {
					if !first {
						result += "\n\t\t\t"
					}

					first = false
					result += dest
				}

				result += "\n\t\t\t\t"

				return result
			}})

		return nil
	},
}

// MAIL FORWARDER CREATE
var mailForwarderCreateAddress string
var mailForwarderCreateDestination string
var mailForwarderCreateCmd = &cobra.Command{
	Use:   "create [mailgroup]",
	Short: "Create a new mail forwarder",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mailgroupID, err := resolveMailgroup(args[0])
		if err != nil {
			return err
		}

		mailforwarder, err := Level27Client.MailgroupsMailforwardersCreate(mailgroupID, l27.MailforwarderCreate{
			Address:     mailForwarderCreateAddress,
			Destination: mailForwarderCreateDestination,
		})

		if err != nil {
			return err
		}

		if optWait {
			mailforwarder, err = waitForStatus(
				func() (l27.Mailforwarder, error) {
					return Level27Client.MailgroupsMailforwardersGetSingle(mailgroupID, mailforwarder.ID)
				},
				func(s l27.Mailforwarder) string { return s.Status },
				"ok",
				[]string{"to_create", "creating"},
			)

			if err != nil {
				return fmt.Errorf("waiting on mailboxforwarder status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(mailforwarder, "templates/entities/mailForwarder/create.tmpl")
		return nil
	},
}

// MAIL FORWARDER DELETE
var mailForwarderDeleteForce bool
var mailForwarderDeleteCmd = &cobra.Command{
	Use:   "delete [mailgroup] [mail forwarder]",
	Short: "Delete a mail forwarder",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		mailgroupID, err := resolveMailgroup(args[0])
		if err != nil {
			return err
		}

		mailforwarderID, err := resolveMailforwarder(mailgroupID, args[1])
		if err != nil {
			return err
		}

		if !mailForwarderDeleteForce {
			mailbox, err := Level27Client.MailgroupsMailforwardersGetSingle(mailgroupID, mailforwarderID)
			if err != nil {
				return err
			}

			if !confirmPrompt(fmt.Sprintf("Delete mail forwarder %s (%d)?", mailbox.Address, mailforwarderID)) {
				return nil
			}
		}

		err = Level27Client.MailgroupsMailforwardersDelete(mailgroupID, mailforwarderID)
		if err != nil {
			return err
		}

		if optWait {
			err = waitForDelete(
				func() (l27.Mailforwarder, error) {
					return Level27Client.MailgroupsMailforwardersGetSingle(mailgroupID, mailforwarderID)
				},
				func(a l27.Mailforwarder) string { return a.Status },
				[]string{"deleting", "to_delete"},
			)

			if err != nil {
				return fmt.Errorf("waiting on mail forwarder status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(nil, "templates/entities/mailForwarder/delete.tmpl")
		return nil
	},
}

// MAIL FORWARDER UPDATE
var mailForwarderUpdateCmd = &cobra.Command{
	Use:   "update [mailgroup] [mail forwarder]",
	Short: "Update settings on a mail forwarder",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		settings, err := loadMergeSettings(updateSettingsFile, updateSettings)
		if err != nil {
			return err
		}

		mailgroupID, err := resolveMailgroup(args[0])
		if err != nil {
			return err
		}

		mailforwarderID, err := resolveMailforwarder(mailgroupID, args[1])
		if err != nil {
			return err
		}

		mailforwarder, err := Level27Client.MailgroupsMailforwardersGetSingle(mailgroupID, mailforwarderID)
		if err != nil {
			return err
		}

		mailforwarderPut := l27.MailforwarderPut{
			Address:     mailforwarder.Address,
			Destination: strings.Join(mailforwarder.Destination, ","),
		}

		data := utils.RoundTripJson(mailforwarderPut).(map[string]interface{})
		data = mergeMaps(data, settings)

		err = Level27Client.MailgroupsMailforwardersUpdate(mailgroupID, mailforwarderID, data)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/mailForwarder/update.tmpl")
		return nil
	},
}

// MAIL FORWARDER DESTINATION
var mailForwarderDestinationCmd = &cobra.Command{
	Use:   "destination",
	Short: "Commands for managing destination addresses on a mail forwarder easily",

	Aliases: []string{"dest"},
}

// MAIL FORWARDER DESTINATION ADD
var mailForwarderDestinationAddCmd = &cobra.Command{
	Use:   "add [mail group] [mail forwarder] [new destination]",
	Short: "Add a single destination address to a mail forwarder",

	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		mailgroupID, err := resolveMailgroup(args[0])
		if err != nil {
			return err
		}

		mailforwarderID, err := resolveMailforwarder(mailgroupID, args[1])
		if err != nil {
			return err
		}

		mailforwarder, err := Level27Client.MailgroupsMailforwardersGetSingle(mailgroupID, mailforwarderID)
		if err != nil {
			return err
		}

		destination := append(mailforwarder.Destination, args[2])

		mailforwarderPut := l27.MailforwarderPut{
			Address:     mailforwarder.Address,
			Destination: strings.Join(destination, ","),
		}

		data := utils.RoundTripJson(mailforwarderPut).(map[string]interface{})
		err = Level27Client.MailgroupsMailforwardersUpdate(mailgroupID, mailforwarderID, data)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/mailForwarderDestination/add.tmpl")
		return nil
	},
}

// MAIL FORWARDER DESTINATION REMOVE
var mailForwarderDestinationRemoveCmd = &cobra.Command{
	Use:   "remove [mail group] [mail forwarder] [old destination]",
	Short: "Remove a single destination address from a mail forwarder",

	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		mailgroupID, err := resolveMailgroup(args[0])
		if err != nil {
			return err
		}

		mailforwarderID, err := resolveMailforwarder(mailgroupID, args[1])
		if err != nil {
			return err
		}

		destAddress := args[2]

		mailforwarder, err := Level27Client.MailgroupsMailforwardersGetSingle(mailgroupID, mailforwarderID)
		if err != nil {
			return err
		}

		destination := mailforwarder.Destination
		idx := indexOf(destination, destAddress)
		if idx == -1 {
			fmt.Printf("'%s' is not a destination on this mail forwarder.", destAddress)
			return nil
		}

		destination = append(destination[:idx], destination[idx+1:]...)

		mailforwarderPut := l27.MailforwarderPut{
			Address:     mailforwarder.Address,
			Destination: strings.Join(destination, ","),
		}

		data := utils.RoundTripJson(mailforwarderPut).(map[string]interface{})
		err = Level27Client.MailgroupsMailforwardersUpdate(mailgroupID, mailforwarderID, data)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/mailForwarderDestination/remove.tmpl")
		return nil
	},
}
