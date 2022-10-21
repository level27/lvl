package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/level27/l27-go"
	"github.com/level27/lvl/utils"
	"github.com/spf13/cobra"
)

func init() {
	// MAIL
	RootCmd.AddCommand(mailCmd)

	// MAIL GET
	mailCmd.AddCommand(mailGetCmd)
	addCommonGetFlags(mailGetCmd)

	// MAIL CREATE
	mailCmd.AddCommand(mailCreateCmd)
	addWaitFlag(mailCreateCmd)
	mailCreateCmd.Flags().StringVar(&mailCreateName, "name", "", "Name of the new mailgroup")
	mailCreateCmd.Flags().StringVar(&mailCreateOrganisation, "organisation", "", "Organisation owning the new mailgroup")
	mailCreateCmd.Flags().StringVar(&mailCreateExternalInfo, "externalInfo", "", "")

	// MAIL DELETE
	mailCmd.AddCommand(mailDeleteCmd)
	addWaitFlag(mailDeleteCmd)
	mailDeleteCmd.Flags().BoolVar(&mailDeleteForce, "force", false, "Do not ask for confirmation to delete the mail group")

	// MAIL UPDATE
	mailCmd.AddCommand(mailUpdateCmd)
	settingsFileFlag(mailUpdateCmd)
	settingString(mailUpdateCmd, updateSettings, "organisation", "New organisation for the mailgroup")
	settingString(mailUpdateCmd, updateSettings, "name", "New name for the mailgroup")

	// MAIL ACTIONS
	mailCmd.AddCommand(mailActionsCmd)

	// MAIL ACTIONS ACTIVATE
	mailActionsCmd.AddCommand(mailActionsActivateCmd)

	// MAIL ACTIONS DEACTIVATE
	mailActionsCmd.AddCommand(mailActionsDeactivateCmd)

	// MAIL DOMAIN
	mailCmd.AddCommand(mailDomainCmd)

	// MAIL DOMAIN LINK
	mailDomainCmd.AddCommand(mailDomainLinkCmd)
	mailDomainLinkCmd.Flags().BoolVar(&mailDomainLinkNoHandleDns, "no-handle-dns", false, "Disable automatic creation of domain records")

	// MAIL DOMAIN UNLINK
	mailDomainCmd.AddCommand(mailDomainUnlinkCmd)

	// MAIL DOMAIN SETPRIMARY
	mailDomainCmd.AddCommand(mailDomainSetPrimaryCmd)

	// MAIL DOMAIN UPDATE
	mailDomainCmd.AddCommand(mailDomainUpdateCmd)
	settingsFileFlag(mailDomainUpdateCmd)
	settingString(mailDomainUpdateCmd, updateSettings, "handleMailDns", "")

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

	// MAIL BOX ADDRESS
	mailBoxCmd.AddCommand(mailBoxAddressCmd)

	// MAIL BOX ADDRESS ADD
	mailBoxAddressCmd.AddCommand(mailBoxAddressAddCmd)

	// MAIL BOX ADDRESS REMOVE
	mailBoxAddressCmd.AddCommand(mailBoxAddressRemoveCmd)

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

	// MAIL ACCESS
	addAccessCmds(mailCmd, "mailgroups", resolveMailgroup)
	// MAIL BILLING
	addBillingCmds(mailCmd, "mailgroups", resolveMailgroup)
	// MAIL JOBS
	addJobCmds(mailCmd, "mailgroup", resolveMailgroup)
	// MAIL INTEGRITY
	addIntegrityCheckCmds(mailCmd, "mailgroups", resolveMailgroup)
}

// Resolve the integer ID of a mail group, from a commandline-passed argument.
// Returns ID if it's a numeric ID, otherwise resolves by name.
func resolveMailgroup(arg string) (int, error) {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.MailgroupsLookup(arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"mailgroup",
		func(group l27.Mailgroup) string { return fmt.Sprintf("%s (%d)", group.Name, group.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, nil
}

func resolveMailbox(mailgroupID int, arg string) (int, error) {
	id, err := strconv.Atoi(arg)
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

func resolveMailboxAdress(mailgroupID int, mailboxID int, arg string) (int, error) {
	id, err := strconv.Atoi(arg)
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

func resolveMailforwarder(mailgroupID int, arg string) (int, error) {
	id, err := strconv.Atoi(arg)
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

// MAIL
var mailCmd = &cobra.Command{
	Use:   "mail",
	Short: "Commands to manage mailgroups and mailboxes",
}

// MAIL GET
var mailGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get mailgroups",

	RunE: func(cmd *cobra.Command, args []string) error {
		mails, err := resolveGets(
			args,
			Level27Client.MailgroupsLookup,
			Level27Client.MailgroupsGetSingle,
			Level27Client.MailgroupsGetList)

		if err != nil {
			return err
		}

		outputFormatTableFuncs(
			mails,
			[]string{"ID", "PRIMARY/NAME", "STATUS", "DOMAINS", "BOXES", "FORWARDERS"},
			[]interface{}{
				"ID",
				mailgroupDisplayName,
				"Status",
				func(m l27.Mailgroup) int { return len(m.Domains) },
				"MailboxCount",
				"MailforwarderCount",
			})

		return nil
	},
}

// MAIL CREATE
var mailCreateName string
var mailCreateOrganisation string
var mailCreateAutoTeams string
var mailCreateExternalInfo string

var mailCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new mail group",
	Long:  "Does not automatically link any domains to the mail group. Use separate commands after the mail group has been created.",

	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		org, err := resolveOrganisation(mailCreateOrganisation)
		if err != nil {
			return err
		}

		create := l27.MailgroupCreate{
			Name:         mailCreateName,
			Organisation: org,
			AutoTeams:    mailCreateAutoTeams,
			ExternalInfo: mailCreateExternalInfo,
			Type:         "level27",
		}

		group, err := Level27Client.MailgroupsCreate(create)
		if err != nil {
			return err
		}

		if optWait {
			group, err = waitForStatus(
				func() (l27.Mailgroup, error) { return Level27Client.MailgroupsGetSingle(group.ID) },
				func(s l27.Mailgroup) string { return s.Status },
				"ok",
				[]string{"to_create", "creating"},
			)

			if err != nil {
				return fmt.Errorf("waiting on mailgroup status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(group, "templates/entities/mail/create.tmpl")
		return nil
	},
}

// MAIL DELETE
var mailDeleteForce bool
var mailDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a mail group",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mailgroupID, err := resolveMailgroup(args[0])
		if err != nil {
			return err
		}

		if !mailDeleteForce {
			mailgroup, err := Level27Client.MailgroupsGetSingle(mailgroupID)
			if err != nil {
				return err
			}

			displayName := mailgroupDisplayName(mailgroup)
			if !confirmPrompt(fmt.Sprintf("Delete mailgroup %s (%d)?", displayName, mailgroupID)) {
				return nil
			}
		}

		err = Level27Client.MailgroupsDelete(mailgroupID)
		if err != nil {
			return err
		}

		if optWait {
			err = waitForDelete(
				func() (l27.Mailgroup, error) { return Level27Client.MailgroupsGetSingle(mailgroupID) },
				func(a l27.Mailgroup) string { return a.Status },
				[]string{"deleting", "to_delete"},
			)

			if err != nil {
				return fmt.Errorf("waiting on mailgroup status failed: %s", err.Error())
			}
		}

		return nil
	},
}

// MAIL UPDATE
var mailUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update settings on a mail group",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		settings, err := loadMergeSettings(updateSettingsFile, updateSettings)
		if err != nil {
			return err
		}

		mailgroupID, err := resolveMailgroup(args[0])
		if err != nil {
			return err
		}

		mailgroup, err := Level27Client.MailgroupsGetSingle(mailgroupID)
		if err != nil {
			return err
		}

		mailgroupPut := l27.MailgroupPut{
			Name:         mailgroup.Name,
			Type:         mailgroup.Type,
			Organisation: mailgroup.Organisation.ID,
			Systemgroup:  mailgroup.Systemgroup.ID,
		}

		data := utils.RoundTripJson(mailgroupPut).(map[string]interface{})
		data = mergeMaps(data, settings)

		data["organisation"], err = resolveOrganisation(fmt.Sprint(data["organisation"]))
		if err != nil {
			return err
		}

		err = Level27Client.MailgroupsUpdate(mailgroupID, data)
		return err
	},
}

// Get the display name for a mailgroup
// Usually just the primary domain name, falling back to .Name if necessary.
func mailgroupDisplayName(m l27.Mailgroup) string {
	if len(m.Domains) == 0 {
		return m.Name
	}

	// Normally we try to display the primary domain, but in some cases there is none.
	// This logic falls back to [0] if there is no primary domain.
	domain := m.Domains[0]
	for _, iterDomain := range m.Domains {
		if iterDomain.MailPrimary {
			domain = iterDomain
		}
	}

	return fmt.Sprintf("%s.%s", domain.Name, domain.Domaintype.Extension)
}

// MAIL ACTIONS
var mailActionsCmd = &cobra.Command{
	Use: "actions",
}

// MAIL ACTIONS ACTIVATE
var mailActionsActivateCmd = &cobra.Command{
	Use: "activate",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mailgroupID, err := resolveMailgroup(args[0])
		if err != nil {
			return err
		}

		_, err = Level27Client.MailgroupsAction(mailgroupID, "activate")
		return err
	},
}

// MAIL ACTIONS ACTIVATE
var mailActionsDeactivateCmd = &cobra.Command{
	Use: "deactivate",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mailgroupID, err := resolveMailgroup(args[0])
		if err != nil {
			return err
		}

		_, err = Level27Client.MailgroupsAction(mailgroupID, "deactivate")
		return err
	},
}

// MAIL DOMAIN
var mailDomainCmd = &cobra.Command{
	Use: "domain",
}

// MAIL DOMAIN LINK
var mailDomainLinkNoHandleDns bool
var mailDomainLinkCmd = &cobra.Command{
	Use:   "link [mailgroup] [domain]",
	Short: "Add a domain to a mail group",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		mailgroupID, err := resolveMailgroup(args[0])
		if err != nil {
			return err
		}

		domainID, err := resolveDomain(args[1])
		if err != nil {
			return err
		}

		_, err = Level27Client.MailgroupsDomainsLink(mailgroupID, l27.MailgroupDomainAdd{
			Domain:        domainID,
			HandleMailDns: !mailDomainLinkNoHandleDns,
		})

		return err
	},
}

// MAIL DOMAIN UNLINK
var mailDomainUnlinkCmd = &cobra.Command{
	Use:   "unlink [mailgroup] [domain]",
	Short: "Remove a domain from a mail group",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		mailgroupID, err := resolveMailgroup(args[0])
		if err != nil {
			return err
		}

		domainID, err := resolveDomain(args[1])
		if err != nil {
			return err
		}

		err = Level27Client.MailgroupsDomainsUnlink(mailgroupID, domainID)
		return err
	},
}

// MAIL DOMAIN SETPRIMARY
var mailDomainSetPrimaryCmd = &cobra.Command{
	Use:   "setprimary [mailgroup] [domain]",
	Short: "Set a domain on a mail group as primary",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		mailgroupID, err := resolveMailgroup(args[0])
		if err != nil {
			return err
		}

		domainID, err := resolveDomain(args[1])
		if err != nil {
			return err
		}

		err = Level27Client.MailgroupsDomainsSetPrimary(mailgroupID, domainID)
		return err
	},
}

// MAIL DOMAIN UPDATE
var mailDomainUpdateCmd = &cobra.Command{
	Use:   "update [mailgroup] [domain]",
	Short: "Update settings on a domain",

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

		domainID, err := resolveDomain(args[1])
		if err != nil {
			return err
		}

		err = Level27Client.MailgroupsDomainsPatch(mailgroupID, domainID, settings)
		return err
	},
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
		return err
	},
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

		_, err = Level27Client.MailgroupsMailboxesAddressesCreate(
			mailgroupID,
			mailboxID,
			l27.MailboxAddressCreate{
				Address: args[2],
			},
		)

		return err
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
		return err
	},
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
		return err
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
		return err
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
		return err
	},
}
