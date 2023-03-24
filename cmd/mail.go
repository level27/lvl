package cmd

import (
	"fmt"

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
func resolveMailgroup(arg string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
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

		outputFormatTemplate(nil, "templates/entities/mail/delete.tmpl")

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
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/mail/update.tmpl")
		return nil
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
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/mail/activate.tmpl")
		return nil
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
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/mail/deactivate.tmpl")
		return nil
	},
}
