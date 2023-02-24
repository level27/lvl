package cmd

import (
	"fmt"

	"github.com/level27/l27-go"
	"github.com/level27/lvl/utils"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(&sshKeyCmd)

	sshKeyCmd.AddCommand(&sshKeyFavoriteCmd)
}

var sshKeyCmd = cobra.Command{
	Use:   "sshkey",
	Short: "Commands for managing SSH keys",
}

var sshKeyFavoriteCmd = cobra.Command{
	Use:   "favorite <SSH key name>",
	Short: "Favorite an SSH key for use in other lvl commands",
	Long: `Favorite an SSH key for use in other lvl commands.
This is used by commands like 'lvl system ssh' to automatically add your SSH key to a system.
The SSH key must first be uploaded to your account on CP4:
https://app.level27.eu/account/profile/ssh-keys`,

	Example: "lvl sshkey favorite pieter-jan",

	Args: cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		user, err := Level27Client.LoginInfo()
		if err != nil {
			return err
		}

		orgID := user.User.Organisation.ID
		userID := user.User.ID

		argKey := args[0]
		sshKeyID, err := resolveSshKey(orgID, userID, argKey)
		if err != nil {
			return err
		}

		sshKey, err := Level27Client.OrganisationUserSshKeysGetSingle(orgID, userID, sshKeyID)
		if err != nil {
			return err
		}

		utils.SaveConfig("ssh_favoriteKey", sshKey.ID)
		fmt.Printf("Key %s (%d) has been set as favorite.", sshKey.Description, sshKey.ID)

		return nil
	},
}

func resolveSshKey(organisationID l27.IntID, userID l27.IntID, arg string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.OrganisationUserSshKeysLookup(organisationID, userID, arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"ssh key",
		func(app l27.SshKey) string { return fmt.Sprintf("%s (%d)", app.Description, app.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, err
}
