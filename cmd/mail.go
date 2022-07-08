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

	// MAIL FORWARDER
	mailCmd.AddCommand(mailForwarderCmd)

	// MAIL FORWARDER GET
	mailForwarderCmd.AddCommand(mailForwarderGetCmd)
	addCommonGetFlags(mailForwarderGetCmd)

	// MAIL FORWARDER CREATE
	mailForwarderCmd.AddCommand(mailForwarderCreateCmd)
	mailForwarderCreateCmd.Flags().StringVar(&mailForwarderCreateAddress, "address", "", "Address of the mail forwarder")
	mailForwarderCreateCmd.Flags().StringVar(&mailForwarderCreateDestination, "destination", "", "Comma-separated list of destination addresses to forward to")

	// MAIL FORWARDER DELETE
	mailForwarderCmd.AddCommand(mailForwarderDeleteCmd)
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
func resolveMailgroup(arg string) int {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id
	}

	return resolveShared(
		Level27Client.MailgroupsLookup(arg),
		arg,
		"mailgroup",
		func(group l27.Mailgroup) string { return fmt.Sprintf("%s (%d)", group.Name, group.ID) }).ID
}

func resolveMailbox(mailgroupID int, arg string) int {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id
	}

	return resolveShared(
		Level27Client.MailgroupsMailboxesLookup(mailgroupID, arg),
		arg,
		"mailbox",
		func(box l27.MailboxShort) string { return fmt.Sprintf("%s (%s, %d)", box.Name, box.Username, box.ID) }).ID
}

func resolveMailboxAdress(mailgroupID int, mailboxID int, arg string) int {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id
	}

	return resolveShared(
		Level27Client.MailgroupsMailboxesAddressesLookup(mailgroupID, mailboxID, arg),
		arg,
		"mailbox address",
		func(address l27.MailboxAddress) string { return fmt.Sprintf("%s (%d)", address.Address, address.ID) }).ID
}

func resolveMailforwarder(mailgroupID int, arg string) int {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id
	}

	return resolveShared(
		Level27Client.MailgroupsMailforwardersLookup(mailgroupID, arg),
		arg,
		"mailforwarder",
		func(app l27.Mailforwarder) string { return fmt.Sprintf("%s (%d)", app.Address, app.ID) }).ID
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
				func(m l27.Mailgroup) int { return len(m.Domains) },
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
	Use:   "create",
	Short: "Create a new mail group",
	Long:  "Does not automatically link any domains to the mail group. Use separate commands after the mail group has been created.",

	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		create := l27.MailgroupCreate{
			Name:         mailCreateName,
			Organisation: resolveOrganisation(mailCreateOrganisation),
			AutoTeams:    mailCreateAutoTeams,
			ExternalInfo: mailCreateExternalInfo,
			Type:         "level27",
		}

		Level27Client.MailgroupsCreate(create)
	},
}

// MAIL DELETE
var mailDeleteForce bool
var mailDeleteCmd = &cobra.Command{
	Use:   "delete",
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
	Use:   "update",
	Short: "Update settings on a mail group",

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		settings := loadMergeSettings(updateSettingsFile, updateSettings)

		mailgroupID := resolveMailgroup(args[0])
		mailgroup := Level27Client.MailgroupsGetSingle(mailgroupID)

		mailgroupPut := l27.MailgroupPut{
			Name:         mailgroup.Name,
			Type:         mailgroup.Type,
			Organisation: mailgroup.Organisation.ID,
			Systemgroup:  mailgroup.Systemgroup.ID,
		}

		data := utils.RoundTripJson(mailgroupPut).(map[string]interface{})
		data = mergeMaps(data, settings)

		data["organisation"] = resolveOrganisation(fmt.Sprint(data["organisation"]))

		Level27Client.MailgroupsUpdate(mailgroupID, data)
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
var mailDomainLinkNoHandleDns bool
var mailDomainLinkCmd = &cobra.Command{
	Use:   "link [mailgroup] [domain]",
	Short: "Add a domain to a mail group",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		mailgroupID := resolveMailgroup(args[0])
		domainID := resolveDomain(args[1])

		Level27Client.MailgroupsDomainsLink(mailgroupID, l27.MailgroupDomainAdd{
			Domain:        domainID,
			HandleMailDns: !mailDomainLinkNoHandleDns,
		})
	},
}

// MAIL DOMAIN UNLINK
var mailDomainUnlinkCmd = &cobra.Command{
	Use:   "unlink [mailgroup] [domain]",
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
	Use:   "setprimary [mailgroup] [domain]",
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
	Use:   "update [mailgroup] [domain]",
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
		addresses := Level27Client.MailgroupsMailboxesAddressesGetList(mailgroupID, mailboxID, l27.CommonGetParams{})
		describe := l27.MailboxDescribe{
			Mailbox:   mailbox,
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

		Level27Client.MailgroupsMailboxesCreate(mailgroupID, l27.MailboxCreate{
			Name:       mailBoxCreateName,
			Password:   mailBoxCreatePassword,
			OooEnabled: mailBoxCreateOooEnabled,
			OooSubject: mailBoxCreateOooSubject,
			OooText:    mailBoxCreateOooText,
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
			if !confirmPrompt(fmt.Sprintf("Delete mailbox %s (%d)?", mailbox.Username, mailboxID)) {
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

		mailboxPut := l27.MailboxPut{
			Name:       mailbox.Name,
			Password:   "",
			OooEnabled: mailbox.OooEnabled,
			OooText:    mailbox.OooText,
			OooSubject: mailbox.OooSubject,
		}

		data := utils.RoundTripJson(mailboxPut).(map[string]interface{})
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
			l27.MailboxAddressCreate{
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
	Run: func(cmd *cobra.Command, args []string) {
		mailgroupID := resolveMailgroup(args[0])

		mailboxes := Level27Client.MailgroupsMailforwardersGetList(mailgroupID, optGetParameters)

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
	},
}

// MAIL FORWARDER CREATE
var mailForwarderCreateAddress string
var mailForwarderCreateDestination string
var mailForwarderCreateCmd = &cobra.Command{
	Use:   "create [mailgroup]",
	Short: "Create a new mail forwarder",

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mailgroupID := resolveMailgroup(args[0])

		Level27Client.MailgroupsMailforwardersCreate(mailgroupID, l27.MailforwarderCreate{
			Address:     mailForwarderCreateAddress,
			Destination: mailForwarderCreateDestination,
		})
	},
}

// MAIL FORWARDER DELETE
var mailForwarderDeleteForce bool
var mailForwarderDeleteCmd = &cobra.Command{
	Use:   "delete [mailgroup] [mail forwarder]",
	Short: "Delete a mail forwarder",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		mailgroupID := resolveMailgroup(args[0])
		mailforwarderID := resolveMailforwarder(mailgroupID, args[1])

		if !mailForwarderDeleteForce {
			mailbox := Level27Client.MailgroupsMailforwardersGetSingle(mailgroupID, mailforwarderID)
			if !confirmPrompt(fmt.Sprintf("Delete mail forwarder %s (%d)?", mailbox.Address, mailforwarderID)) {
				return
			}
		}

		Level27Client.MailgroupsMailforwardersDelete(mailgroupID, mailforwarderID)
	},
}

// MAIL FORWARDER UPDATE
var mailForwarderUpdateCmd = &cobra.Command{
	Use:   "update [mailgroup] [mail forwarder]",
	Short: "Update settings on a mail forwarder",

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		settings := loadMergeSettings(updateSettingsFile, updateSettings)

		mailgroupID := resolveMailgroup(args[0])
		mailforwarderID := resolveMailforwarder(mailgroupID, args[1])
		mailforwarder := Level27Client.MailgroupsMailforwardersGetSingle(mailgroupID, mailforwarderID)

		mailforwarderPut := l27.MailforwarderPut{
			Address:     mailforwarder.Address,
			Destination: strings.Join(mailforwarder.Destination, ","),
		}

		data := utils.RoundTripJson(mailforwarderPut).(map[string]interface{})
		data = mergeMaps(data, settings)

		Level27Client.MailgroupsMailforwardersUpdate(mailgroupID, mailforwarderID, data)
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
	Run: func(cmd *cobra.Command, args []string) {
		mailgroupID := resolveMailgroup(args[0])
		mailforwarderID := resolveMailforwarder(mailgroupID, args[1])

		mailforwarder := Level27Client.MailgroupsMailforwardersGetSingle(mailgroupID, mailforwarderID)
		destination := append(mailforwarder.Destination, args[2])

		mailforwarderPut := l27.MailforwarderPut{
			Address:     mailforwarder.Address,
			Destination: strings.Join(destination, ","),
		}

		data := utils.RoundTripJson(mailforwarderPut).(map[string]interface{})
		Level27Client.MailgroupsMailforwardersUpdate(mailgroupID, mailforwarderID, data)
	},
}

// MAIL FORWARDER DESTINATION REMOVE
var mailForwarderDestinationRemoveCmd = &cobra.Command{
	Use:   "remove [mail group] [mail forwarder] [old destination]",
	Short: "Remove a single destination address from a mail forwarder",

	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		mailgroupID := resolveMailgroup(args[0])
		mailforwarderID := resolveMailforwarder(mailgroupID, args[1])
		destAddress := args[2]

		mailforwarder := Level27Client.MailgroupsMailforwardersGetSingle(mailgroupID, mailforwarderID)
		destination := mailforwarder.Destination
		idx := indexOf(destination, destAddress)
		if idx == -1 {
			fmt.Printf("'%s' is not a destination on this mail forwarder.", destAddress)
			return
		}

		destination = append(destination[:idx], destination[idx+1:]...)

		mailforwarderPut := l27.MailforwarderPut{
			Address:     mailforwarder.Address,
			Destination: strings.Join(destination, ","),
		}

		data := utils.RoundTripJson(mailforwarderPut).(map[string]interface{})
		Level27Client.MailgroupsMailforwardersUpdate(mailgroupID, mailforwarderID, data)
	},
}
