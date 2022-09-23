package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(tokenCmd)
}

var tokenCmd = &cobra.Command{
	Use:    "token",
	Short:  "Show authentication token.",
	Long:   "Shows the authentication token. This is intended for manually doing API-requests in scripts and such.",
	Hidden: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		key := viper.GetString("apiKey")

		if key == "" {
			return fmt.Errorf("API token not set, log in first")
		}

		fmt.Print(key)
		return nil
	},
}
