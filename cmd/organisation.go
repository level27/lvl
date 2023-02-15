package cmd

import (
	"fmt"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(organisationCmd)

	organisationCmd.AddCommand(organisationGetCmd)
	addCommonGetFlags(organisationGetCmd)

	organisationCmd.AddCommand(organisationUserCmd)

	organisationUserCmd.AddCommand(organisationUserSshKeyCmd)

	organisationUserSshKeyCmd.AddCommand(organisationUserSshKeyCreateCmd)
	organisationUserSshKeyCreateCmd.Flags().StringVarP(&organisationUserSshKeyCreateDescription, "description", "d", "", "Description new SSH key")
	organisationUserSshKeyCreateCmd.Flags().StringVarP(&organisationUserSshKeyCreateContent, "content", "c", "", "Content of the new SSH key, e.g. ssh-rsa ...")
	organisationUserSshKeyCreateCmd.MarkFlagRequired("description")
	organisationUserSshKeyCreateCmd.MarkFlagRequired("content")
}

var organisationCmd = &cobra.Command{
	Use:   "organisation",
	Short: "Commands for managing organisations",
}

var organisationGetCmd = &cobra.Command{
	Use: "get",

	Args: cobra.ArbitraryArgs,
	RunE: func(ccmd *cobra.Command, args []string) error {
		ids, err := convertStringsToIDs(args)
		if err != nil {
			return err
		}

		options, err := getOrganisations(ids)
		if err != nil {
			return err
		}

		outputFormatTable(options, []string{"ID", "NAME"}, []string{"ID", "Name"})
		return nil
	},
}

var organisationUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Commands for managing users on an organisation",
}

var organisationUserSshKeyCmd = &cobra.Command{
	Use:   "sshkey",
	Short: "Commands for managing SSH keys on a user",
}

var organisationUserSshKeyCreateDescription string
var organisationUserSshKeyCreateContent string
var organisationUserSshKeyCreateCmd = &cobra.Command{
	Use:   "create <organisation> <user> -d <key description> -c <key content>",
	Short: "Add a new SSH key to a user",

	Args: cobra.ExactArgs(2),
	RunE: func(ccmd *cobra.Command, args []string) error {
		organisation, err := resolveOrganisation(args[0])
		if err != nil {
			return err
		}

		user, err := resolveOrganisationUser(organisation, args[1])
		if err != nil {
			return err
		}

		content, err := readArgFileSupported(organisationUserSshKeyCreateContent)
		if err != nil {
			return err
		}

		create := l27.SshKeyCreate{
			Content:     content,
			Description: organisationUserSshKeyCreateDescription,
		}

		resp, err := Level27Client.OrganisationUserSshKeysCreate(organisation, user, create)
		if err != nil {
			return err
		}

		outputFormatTemplate(resp, "templates/entities/organisationUserSshkey/create.tmpl")

		return nil
	},
}

func resolveOrganisation(arg string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.LookupOrganisation(arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"organisation",
		func(app l27.Organisation) string { return fmt.Sprintf("%s (%d)", app.Name, app.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, nil
}

func resolveOrganisationUser(organisationID l27.IntID, arg string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.LookupOrganisationUser(organisationID, arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"user",
		func(user l27.OrganisationUser) string {
			return fmt.Sprintf("%s %s (%d)", user.FirstName, user.LastName, user.ID)
		})

	if err != nil {
		return 0, err
	}

	return res.ID, nil
}

func getOrganisations(ids []l27.IntID) ([]l27.Organisation, error) {
	c := Level27Client
	if len(ids) == 0 {
		return c.Organisations(optGetParameters)
	} else {
		organisations := make([]l27.Organisation, len(ids))
		for idx, id := range ids {
			var err error
			organisations[idx], err = c.Organisation(id)
			if err != nil {
				return nil, err
			}
		}

		return organisations, nil
	}
}
