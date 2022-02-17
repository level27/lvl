package cmd

import "github.com/spf13/cobra"

// MAIN COMMAND
var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Commands for managing systems",
}

func init() {
	//adding main command to root
	RootCmd.AddCommand(systemCmd)

	//Toplevel subcommands (get/post)
	systemCmd.AddCommand(systemGetCmd)
}


var systemGetCmd = &cobra.Command{
	Use: "get",
	Short: "get a list of all curent systems",
	Run: func(cmd *cobra.Command, args []string) {

	},
}