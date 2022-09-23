package cmd

import (
	"fmt"
	"strconv"

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
	Use:   "favorite",
	Short: "Favorite an SSH key for use in other lvl commands",

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

		utils.SaveConfig("ssh_favoriteKey", sshKey.Id)
		fmt.Printf("Key %s (%d) has been set as favorite.", sshKey.Description, sshKey.Id)

		return nil
	},
}

func resolveSshKey(organisationID int, userID int, arg string) (int, error) {
	id, err := strconv.Atoi(arg)
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
		func(app l27.SshKey) string { return fmt.Sprintf("%s (%d)", app.Description, app.Id) })

	if err != nil {
		return 0, err
	}

	return res.Id, err
}
