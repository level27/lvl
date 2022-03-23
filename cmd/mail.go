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