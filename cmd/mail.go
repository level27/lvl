package cmd

import (
	"fmt"
	"strconv"

	"bitbucket.org/level27/lvl/types"
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
	mailCreateCmd.Flags().StringVar(&mailCreateName, "name", "", "Name of the new mailgroup")
	mailCreateCmd.Flags().StringVar(&mailCreateOrganisation, "organisation", "", "Organisation owning the new mailgroup")
	mailCreateCmd.Flags().StringVar(&mailCreateExternalInfo, "externalInfo", "", "")

	// MAIL DELETE
	mailCmd.AddCommand(mailDeleteCmd)
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
	mailBoxCreateCmd.Flags().StringVar(&mailBoxCreateName, "name", "", "Name of the new mail box")
	mailBoxCreateCmd.Flags().StringVar(&mailBoxCreatePassword, "password", "", "Password of the new mail box")
	mailBoxCreateCmd.Flags().BoolVar(&mailBoxCreateOooEnabled, "oooEnabled", false, "Whether the account is marked as out-of-office.")
	mailBoxCreateCmd.Flags().StringVar(&mailBoxCreateOooSubject, "oooSubject", "", "Subject line for out-of-office status. Required if --oooEnabled is true.")
	mailBoxCreateCmd.Flags().StringVar(&mailBoxCreateOooText, "oooText", "", "Body for out-of-office status. Required if --oooEnabled is true.")

	// MAIL BOX DELETE
	mailBoxCmd.AddCommand(mailBoxDeleteCmd)
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

}

// Resolve the integer ID of a mail group, from a commandline-passed argument.
// Returns ID if it's a numeric ID, otherwise resolves by name.
func resolveMailgroup(arg string) int {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id
	}

	mailgroup := Level27Client.MailgroupsLookup(arg)
	if mailgroup == nil {
		cobra.CheckErr(fmt.Sprintf("Unable to find mailgroup: %s", arg))
		return 0
	}
	return mailgroup.ID
}

func resolveMailbox(mailgroupID int, arg string) int {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id
	}

	mailbox := Level27Client.MailgroupsMailboxesLookup(mailgroupID, arg)
	if mailbox == nil {
		cobra.CheckErr(fmt.Sprintf("Unable to find mailbox: %s", arg))
		return 0
	}
	return mailbox.ID
}

func resolveMailboxAdress(mailgroupID int, mailboxID int, arg string) int {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id
	}

	address := Level27Client.MailgroupsMailboxesAddressesLookup(mailgroupID, mailboxID, arg)
	if address == nil {
		cobra.CheckErr(fmt.Sprintf("Unable to find mailbox address: %s", arg))
		return 0
	}
	return address.ID
}

// MAIL
var mailCmd = &cobra.Command{
	Use: "mail",
	Short: "Commands to manage mailgroups and mailboxes",
}

// MAIL GET
var mailGetCmd = &cobra.Command{
	Use: "get",
	Short: "Get mailgroups",

	Run: func(cmd *cobra.Command, args []string) {
		mails := resolveGets(
			args,
			Level27Client.MailgroupsLookup,
			Level27Client.MailgroupsGetSingle,
			Level27Client.MailgroupsGetList)

		outputFormatTableFuncs(
			mails,
			[]string{"ID", "PRIMARY/NAME", "STATUS", "DOMAINS", "BOXES", "FORWARDERS"},
			[]interface{}{
				"ID",
				mailgroupDisplayName,
				"Status",
				func(m types.Mailgroup) int { return len(m.Domains) },
				"MailboxCount",
				"MailforwarderCount",
			})
	},
}

// MAIL CREATE
var mailCreateName string
var mailCreateOrganisation string
var mailCreateAutoTeams string
var mailCreateExternalInfo string

var mailCreateCmd = &cobra.Command{
	Use: "create",
	Short: "Create a new mail group",
	Long: "Does not automatically link any domains to the mail group. Use separate commands after the mail group has been created.",

	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		create := types.MailgroupCreate{
			Name: mailCreateName,
			Organisation: resolveOrganisation(mailCreateOrganisation),
			AutoTeams: mailCreateAutoTeams,
			ExternalInfo: mailCreateExternalInfo,
			Type: "level27",
		}

		Level27Client.MailgroupsCreate(create)
	},
}

// MAIL DELETE
var mailDeleteForce bool
var mailDeleteCmd = &cobra.Command{
	Use: "delete",
	Short: "Delete a mail group",

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mailgroupID := resolveMailgroup(args[0])

		if !mailDeleteForce {
			mailgroup := Level27Client.MailgroupsGetSingle(mailgroupID)
			displayName := mailgroupDisplayName(mailgroup)
			if !confirmPrompt(fmt.Sprintf("Delete mailgroup %s (%d)?", displayName, mailgroupID)) {
				return
			}
		}

		Level27Client.MailgroupsDelete(mailgroupID)
	},
}

// MAIL UPDATE
var mailUpdateCmd = &cobra.Command{
	Use: "update",
	Short: "Update settings on a mail group",

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		settings := loadMergeSettings(updateSettingsFile, updateSettings)

		mailgroupID := resolveMailgroup(args[0])
		mailgroup := Level27Client.MailgroupsGetSingle(mailgroupID)

		mailgroupPut := types.MailgroupPut {
			Name: mailgroup.Name,
			Type: mailgroup.Type,
			Organisation: mailgroup.Organisation.ID,
			Systemgroup: mailgroup.Systemgroup.ID,
		}

		data := roundTripJson(mailgroupPut).(map[string]interface{})
		data = mergeMaps(data, settings)

		data["organisation"] = resolveOrganisation(fmt.Sprint(data["organisation"]))

		Level27Client.MailgroupsUpdate(mailgroupID, data)
	},
}

// Get the display name for a mailgroup
// Usually just the primary domain name, falling back to .Name if necessary.
func mailgroupDisplayName(m types.Mailgroup) string {
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
	Run: func(cmd *cobra.Command, args []string) {
		mailgroupID := resolveMailgroup(args[0])

		Level27Client.MailgroupsAction(mailgroupID, "activate")
	},
}

// MAIL ACTIONS ACTIVATE
var mailActionsDeactivateCmd = &cobra.Command{
	Use: "deactivate",

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mailgroupID := resolveMailgroup(args[0])

		Level27Client.MailgroupsAction(mailgroupID, "deactivate")
	},
}


// MAIL DOMAIN
var mailDomainCmd = &cobra.Command{
	Use: "domain",
}

// MAIL DOMAIN LINK
var mailDomainLinkNoHandleDns bool;
var mailDomainLinkCmd = &cobra.Command{
	Use: "link [mailgroup] [domain]",
	Short: "Add a domain to a mail group",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		mailgroupID := resolveMailgroup(args[0])
		domainID := resolveDomain(args[1])

		Level27Client.MailgroupsDomainsLink(mailgroupID, types.MailgroupDomainAdd{
			Domain: domainID,
			HandleMailDns: !mailDomainLinkNoHandleDns,
		})
	},
}

// MAIL DOMAIN UNLINK
var mailDomainUnlinkCmd = &cobra.Command{
	Use: "unlink [mailgroup] [domain]",
	Short: "Remove a domain from a mail group",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		mailgroupID := resolveMailgroup(args[0])
		domainID := resolveDomain(args[1])

		Level27Client.MailgroupsDomainsUnlink(mailgroupID, domainID)
	},
}

// MAIL DOMAIN SETPRIMARY
var mailDomainSetPrimaryCmd = &cobra.Command{
	Use: "setprimary [mailgroup] [domain]",
	Short: "Set a domain on a mail group as primary",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		mailgroupID := resolveMailgroup(args[0])
		domainID := resolveDomain(args[1])

		Level27Client.MailgroupsDomainsSetPrimary(mailgroupID, domainID)
	},
}


// MAIL DOMAIN UPDATE
var mailDomainUpdateCmd = &cobra.Command{
	Use: "update [mailgroup] [domain]",
	Short: "Update settings on a domain",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		settings := loadMergeSettings(updateSettingsFile, updateSettings)

		mailgroupID := resolveMailgroup(args[0])
		domainID := resolveDomain(args[1])

		Level27Client.MailgroupsDomainsPatch(mailgroupID, domainID, settings)
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
	Run: func(cmd *cobra.Command, args []string) {
		mailgroupID := resolveMailgroup(args[0])

		mailboxes := Level27Client.MailgroupsMailboxesGetList(mailgroupID, optGetParameters)

		outputFormatTable(
			mailboxes,
			[]string{"ID", "Name", "Username", "Status"},
			[]string{"ID", "Name", "Username", "Status"})
	},
}

// MAIL BOX DESCRIBE
var mailBoxDescribeCmd = &cobra.Command{
	Use: "describe",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		mailgroupID := resolveMailgroup(args[0])
		mailboxID := resolveMailbox(mailgroupID, args[1])

		mailbox := Level27Client.MailgroupsMailboxesGetSingle(mailgroupID, mailboxID)
		addresses := Level27Client.MailgroupsMailboxesAddressesGetList(mailgroupID, mailboxID, types.CommonGetParams{})
		describe := types.MailboxDescribe{
			Mailbox: mailbox,
			Addresses: addresses,
		}

		outputFormatTemplate(describe, "templates/mailbox.tmpl")
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
	Run: func(cmd *cobra.Command, args []string) {
		mailgroupID := resolveMailgroup(args[0])

		Level27Client.MailgroupsMailboxesCreate(mailgroupID, types.MailboxCreate{
			Name: mailBoxCreateName,
			Password: mailBoxCreatePassword,
			OooEnabled: mailBoxCreateOooEnabled,
			OooSubject: mailBoxCreateOooSubject,
			OooText: mailBoxCreateOooText,
		})
	},
}

// MAIL BOX DELETE
var mailBoxDeleteForce bool
var mailBoxDeleteCmd = &cobra.Command{
	Use: "delete",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		mailgroupID := resolveMailgroup(args[0])
		mailboxID := resolveMailbox(mailgroupID, args[1])

		if !mailBoxDeleteForce {
			mailbox := Level27Client.MailgroupsMailboxesGetSingle(mailgroupID, mailboxID)
			if !confirmPrompt(fmt.Sprintf("Delete mailbox %s (%d)?", mailbox.Username, mailgroupID)) {
				return
			}
		}

		Level27Client.MailgroupsMailboxesDelete(mailgroupID, mailboxID)
	},
}

// MAIL BOX UPDATE
var mailBoxUpdateCmd = &cobra.Command{
	Use: "update",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		settings := loadMergeSettings(updateSettingsFile, updateSettings)

		mailgroupID := resolveMailgroup(args[0])
		mailboxID := resolveMailbox(mailgroupID, args[1])
		mailbox := Level27Client.MailgroupsMailboxesGetSingle(mailgroupID, mailboxID)

		mailboxPut := types.MailboxPut {
			Name: mailbox.Name,
			Password: "",
			OooEnabled: mailbox.OooEnabled,
			OooText: mailbox.OooText,
			OooSubject: mailbox.OooSubject,
		}

		data := roundTripJson(mailboxPut).(map[string]interface{})
		data = mergeMaps(data, settings)

		Level27Client.MailgroupsMailboxesUpdate(mailgroupID, mailboxID, data)
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
	Run: func(cmd *cobra.Command, args []string) {
		mailgroupID := resolveMailgroup(args[0])
		mailboxID := resolveMailbox(mailgroupID, args[1])

		Level27Client.MailgroupsMailboxesAddressesCreate(
			mailgroupID,
			mailboxID,
			types.MailboxAddressCreate{
				Address: args[2],
			},
		)
	},
}

// MAIL BOX ADDRESS REMOVE
var mailBoxAddressRemoveCmd = &cobra.Command{
	Use: "remove",

	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		mailgroupID := resolveMailgroup(args[0])
		mailboxID := resolveMailbox(mailgroupID, args[1])
		addressID := resolveMailboxAdress(mailgroupID, mailboxID, args[2])

		Level27Client.MailgroupsMailboxesAddressesDelete(
			mailgroupID,
			mailboxID,
			addressID,
		)
	},
}