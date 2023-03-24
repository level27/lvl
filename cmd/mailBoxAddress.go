package cmd

import (
	"fmt"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

func init() {
	// MAIL BOX ADDRESS
	mailBoxCmd.AddCommand(mailBoxAddressCmd)

	// MAIL BOX ADDRESS ADD
	mailBoxAddressCmd.AddCommand(mailBoxAddressAddCmd)

	// MAIL BOX ADDRESS REMOVE
	mailBoxAddressCmd.AddCommand(mailBoxAddressRemoveCmd)
}

func resolveMailboxAdress(mailgroupID l27.IntID, mailboxID l27.IntID, arg string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.MailgroupsMailboxesAddressesLookup(mailgroupID, mailboxID, arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"mailbox address",
		func(address l27.MailboxAddress) string { return fmt.Sprintf("%s (%d)", address.Address, address.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, nil
}

// MAIL BOX ADDRESS
var mailBoxAddressCmd = &cobra.Command{
	Use: "address",
}

// MAIL BOX ADDRESS ADD
var mailBoxAddressAddCmd = &cobra.Command{
	Use: "add",

	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		mailgroupID, err := resolveMailgroup(args[0])
		if err != nil {
			return err
		}

		mailboxID, err := resolveMailbox(mailgroupID, args[1])
		if err != nil {
			return err
		}

		address, err := Level27Client.MailgroupsMailboxesAddressesCreate(
			mailgroupID,
			mailboxID,
			l27.MailboxAddressCreate{
				Address: args[2],
			},
		)
		if err != nil {
			return err
		}

		outputFormatTemplate(address, "templates/entities/mailBoxAddress/add.tmpl")
		return nil
	},
}

// MAIL BOX ADDRESS REMOVE
var mailBoxAddressRemoveCmd = &cobra.Command{
	Use: "remove",

	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		mailgroupID, err := resolveMailgroup(args[0])
		if err != nil {
			return err
		}

		mailboxID, err := resolveMailbox(mailgroupID, args[1])
		if err != nil {
			return err
		}

		addressID, err := resolveMailboxAdress(mailgroupID, mailboxID, args[2])
		if err != nil {
			return err
		}

		err = Level27Client.MailgroupsMailboxesAddressesDelete(
			mailgroupID,
			mailboxID,
			addressID,
		)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/mailBoxAddress/remove.tmpl")
		return err
	},
}
