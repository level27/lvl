package cmd

import (
	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

func init() {
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

	// MAIL DOMAIN DKIM
	mailDomainCmd.AddCommand(mailDomainDkimCmd)

	// MAIL DOMAIN DKIM ENABLE
	mailDomainDkimCmd.AddCommand(mailDomainDkimEnableCmd)

	// MAIL DOMAIN DKIM DISABLE
	mailDomainDkimCmd.AddCommand(mailDomainDkimDisableCmd)
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
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/mailDomain/link.tmpl")
		return nil
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
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/mailDomain/unlink.tmpl")
		return nil
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
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/mailDomain/setPrimary.tmpl")
		return nil
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
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/mailDomain/update.tmpl")
		return err
	},
}

// MAIL DOMAIN DKIM
var mailDomainDkimCmd = &cobra.Command{
	Use:   "dkim",
	Short: "Commands to manage DKIM on mail domains",
}

// MAIL DOMAIN DKIM ENABLE
var mailDomainDkimEnableCmd = &cobra.Command{
	Use:   "enable <mail group> <mail domain>",
	Short: "Enable DKIM on a mail domain",

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

		result, err := Level27Client.MailgroupsDomainAction(mailgroupID, domainID, "createDkim")
		if err != nil {
			return err
		}

		outputFormatTemplate(result, "templates/entities/mailDomain/enableDkim.tmpl")
		return nil
	},
}

// MAIL DOMAIN DKIM DISABLE
var mailDomainDkimDisableCmd = &cobra.Command{
	Use:   "disable <mail group> <mail domain>",
	Short: "Disable DKIM on a mail domain",

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

		result, err := Level27Client.MailgroupsDomainAction(mailgroupID, domainID, "deleteDkim")
		if err != nil {
			return err
		}

		outputFormatTemplate(result, "templates/entities/mailDomain/disableDkim.tmpl")
		return nil
	},
}
