package cmd

import (
	"fmt"

	"github.com/level27/l27-go"
	"github.com/level27/lvl/utils"
	"github.com/spf13/cobra"
)

func init() {
	// MAIL BOX
	mailCmd.AddCommand(mailBoxCmd)

	// MAIL BOX GET
	mailBoxCmd.AddCommand(mailBoxGetCmd)
	addCommonGetFlags(mailBoxGetCmd)

	// MAIL BOX DESCRIBE
	mailBoxCmd.AddCommand(mailBoxDescribeCmd)

	// MAIL BOX CREATE
	mailBoxCmd.AddCommand(mailBoxCreateCmd)
	addWaitFlag(mailBoxCreateCmd)
	mailBoxCreateCmd.Flags().StringVar(&mailBoxCreateName, "name", "", "Name of the new mail box")
	mailBoxCreateCmd.Flags().StringVar(&mailBoxCreatePassword, "password", "", "Password of the new mail box")
	mailBoxCreateCmd.Flags().BoolVar(&mailBoxCreateOooEnabled, "oooEnabled", false, "Whether the account is marked as out-of-office.")
	mailBoxCreateCmd.Flags().StringVar(&mailBoxCreateOooSubject, "oooSubject", "", "Subject line for out-of-office status. Required if --oooEnabled is true.")
	mailBoxCreateCmd.Flags().StringVar(&mailBoxCreateOooText, "oooText", "", "Body for out-of-office status. Required if --oooEnabled is true.")

	// MAIL BOX DELETE
	mailBoxCmd.AddCommand(mailBoxDeleteCmd)
	addWaitFlag(mailBoxDeleteCmd)
	mailBoxDeleteCmd.Flags().BoolVar(&mailBoxDeleteForce, "force", false, "Do not ask for confirmation to delete the mail box")

	// MAIL BOX UPDATE
	mailBoxCmd.AddCommand(mailBoxUpdateCmd)
	settingsFileFlag(mailBoxUpdateCmd)
	settingString(mailBoxUpdateCmd, updateSettings, "name", "New name for the mail box")
	settingString(mailBoxUpdateCmd, updateSettings, "password", "New password for the mail box")
	settingBool(mailBoxUpdateCmd, updateSettings, "oooEnabled", "New status for out-of-office")
	settingString(mailBoxUpdateCmd, updateSettings, "oooSubject", "New subject line for out-of-office status")
	settingString(mailBoxUpdateCmd, updateSettings, "oooText", "New body for out-of-office status")

}

func resolveMailbox(mailgroupID l27.IntID, arg string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.MailgroupsMailboxesLookup(mailgroupID, arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"mailbox",
		func(box l27.MailboxShort) string { return fmt.Sprintf("%s (%s, %d)", box.Name, box.Username, box.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, nil
}

// MAIL BOX
var mailBoxCmd = &cobra.Command{
	Use: "box",
}

// MAIL BOX GET
var mailBoxGetCmd = &cobra.Command{
	Use: "get",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mailgroupID, err := resolveMailgroup(args[0])
		if err != nil {
			return err
		}

		mailboxes, err := Level27Client.MailgroupsMailboxesGetList(mailgroupID, optGetParameters)
		if err != nil {
			return err
		}

		outputFormatTable(
			mailboxes,
			[]string{"ID", "Name", "Username", "Status"},
			[]string{"ID", "Name", "Username", "Status"})

		return nil
	},
}

// MAIL BOX DESCRIBE
var mailBoxDescribeCmd = &cobra.Command{
	Use: "describe",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		mailgroupID, err := resolveMailgroup(args[0])
		if err != nil {
			return err
		}

		mailboxID, err := resolveMailbox(mailgroupID, args[1])
		if err != nil {
			return err
		}

		mailbox, err := Level27Client.MailgroupsMailboxesGetSingle(mailgroupID, mailboxID)
		if err != nil {
			return err
		}

		addresses, err := Level27Client.MailgroupsMailboxesAddressesGetList(mailgroupID, mailboxID, l27.CommonGetParams{})
		if err != nil {
			return err
		}

		describe := l27.MailboxDescribe{
			Mailbox:   mailbox,
			Addresses: addresses,
		}

		outputFormatTemplate(describe, "templates/mailbox.tmpl")
		return nil
	},
}

// MAIL BOX CREATE
var mailBoxCreateName string
var mailBoxCreatePassword string
var mailBoxCreateOooEnabled bool
var mailBoxCreateOooSubject string
var mailBoxCreateOooText string
var mailBoxCreateCmd = &cobra.Command{
	Use: "create",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mailgroupID, err := resolveMailgroup(args[0])
		if err != nil {
			return err
		}

		mailbox, err := Level27Client.MailgroupsMailboxesCreate(mailgroupID, l27.MailboxCreate{
			Name:       mailBoxCreateName,
			Password:   mailBoxCreatePassword,
			OooEnabled: mailBoxCreateOooEnabled,
			OooSubject: mailBoxCreateOooSubject,
			OooText:    mailBoxCreateOooText,
		})

		if err != nil {
			return err
		}

		if optWait {
			mailbox, err = waitForStatus(
				func() (l27.Mailbox, error) {
					return Level27Client.MailgroupsMailboxesGetSingle(mailgroupID, mailbox.ID)
				},
				func(s l27.Mailbox) string { return s.Status },
				"ok",
				[]string{"to_create", "creating"},
			)

			if err != nil {
				return fmt.Errorf("waiting on mailbox status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(mailbox, "templates/entities/mailBox/create.tmpl")
		return nil
	},
}

// MAIL BOX DELETE
var mailBoxDeleteForce bool
var mailBoxDeleteCmd = &cobra.Command{
	Use: "delete",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		mailgroupID, err := resolveMailgroup(args[0])
		if err != nil {
			return err
		}

		mailboxID, err := resolveMailbox(mailgroupID, args[1])
		if err != nil {
			return err
		}

		if !mailBoxDeleteForce {
			mailbox, err := Level27Client.MailgroupsMailboxesGetSingle(mailgroupID, mailboxID)
			if err != nil {
				return err
			}

			if !confirmPrompt(fmt.Sprintf("Delete mailbox %s (%d)?", mailbox.Username, mailboxID)) {
				return nil
			}
		}

		err = Level27Client.MailgroupsMailboxesDelete(mailgroupID, mailboxID)
		if err != nil {
			return err
		}

		if optWait {
			err = waitForDelete(
				func() (l27.Mailbox, error) { return Level27Client.MailgroupsMailboxesGetSingle(mailgroupID, mailboxID) },
				func(a l27.Mailbox) string { return a.Status },
				[]string{"deleting", "to_delete"},
			)

			if err != nil {
				return fmt.Errorf("waiting on mailbox status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(nil, "templates/entities/mailBox/delete.tmpl")
		return nil
	},
}

// MAIL BOX UPDATE
var mailBoxUpdateCmd = &cobra.Command{
	Use: "update",

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

		mailboxID, err := resolveMailbox(mailgroupID, args[1])
		if err != nil {
			return err
		}

		mailbox, err := Level27Client.MailgroupsMailboxesGetSingle(mailgroupID, mailboxID)
		if err != nil {
			return err
		}

		mailboxPut := l27.MailboxPut{
			Name:       mailbox.Name,
			Password:   "",
			OooEnabled: mailbox.OooEnabled,
			OooText:    mailbox.OooText,
			OooSubject: mailbox.OooSubject,
		}

		data := utils.RoundTripJson(mailboxPut).(map[string]interface{})
		data = mergeMaps(data, settings)

		err = Level27Client.MailgroupsMailboxesUpdate(mailgroupID, mailboxID, data)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/mailBox/update.tmpl")
		return nil
	},
}
