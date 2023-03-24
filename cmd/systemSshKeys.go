package cmd

import (
	"errors"
	"fmt"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	//-------------------------------------  SYSTEMS/SSH KEYS (get/ add / delete) --------------------------------------
	// #region SYSTEMS/SSH KEYS (get/ add / describe / delete)

	// SSH KEYS
	systemCmd.AddCommand(systemSshKeysCmd)

	systemSshKeysCmd.AddCommand(systemSshKeysGetCmd)
	addCommonGetFlags(systemSshKeysGetCmd)

	systemSshKeysCmd.AddCommand(systemSshKeysAddCmd)
	systemSshKeysCmd.AddCommand(systemSshKeysRemoveCmd)
}

//------------------------------------------------- SYSTEMS / SSH KEYS (GET / ADD / DELETE)

var systemSshKeysCmd = &cobra.Command{
	Use: "sshkeys",
}

// #region SYSTEMS/SHH KEYS (GET / ADD / DELETE)

// --- GET
var systemSshKeysGetCmd = &cobra.Command{
	Use: "get",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		keys, err := Level27Client.SystemGetSshKeys(id, optGetParameters)
		if err != nil {
			return err
		}

		outputFormatTable(keys, []string{"ID", "DESCRIPTION", "STATUS", "FINGERPRINT"}, []string{"ID", "Description", "ShsStatus", "Fingerprint"})
		return nil
	},
}

// --- ADD
var systemSshKeysAddCmd = &cobra.Command{
	Use: "add",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		keyName := args[1]
		keyID, err := l27.ParseID(keyName)
		if err != nil {
			user := viper.GetInt32("user_id")
			org := viper.GetInt32("org_id")
			system, err := Level27Client.LookupSystemNonAddedSshkey(systemID, org, user, keyName)
			if err != nil {
				return err
			}

			if system == nil {
				existing, err := Level27Client.LookupSystemSshkey(systemID, keyName)
				if err != nil {
					return err
				}

				if existing != nil {
					return errors.New("SSH key already exists on system")
				}

				return fmt.Errorf("unable to find SSH key to add: '%s'", keyName)
			}

			keyID = system.ID
		}

		key, err := Level27Client.SystemAddSshKey(systemID, keyID)
		if err != nil {
			return err
		}

		outputFormatTemplate(key, "templates/entities/systemSshkey/add.tmpl")
		return nil
	},
}

// --- DELETE
var systemSshKeysRemoveCmd = &cobra.Command{
	Use: "remove",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		keyName := args[1]
		keyID, err := l27.ParseID(keyName)
		if err != nil {
			existing, err := Level27Client.LookupSystemSshkey(systemID, keyName)
			if err != nil {
				return err
			}

			if existing == nil {
				return fmt.Errorf("unable to find SSH key to remove: %s", keyName)
			}

			keyID = existing.ID
		}

		err = Level27Client.SystemRemoveSshKey(systemID, keyID)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/systemSshkey/remove.tmpl")
		return nil
	},
}

// #endregion
