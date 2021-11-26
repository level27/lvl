package get

import (
	"bitbucket.org/level27/lvl/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var optNumber = "20"
var optFilter = ""

var GetCmd = &cobra.Command{
	Use: "get",
}

func init() {
	GetCmd.PersistentFlags().StringVarP(&optNumber, "number", "n", optNumber, "How many things should we retrieve from the API?")
	GetCmd.PersistentFlags().StringVarP(&optFilter, "filter", "f", optFilter, "How to filter API results?")
	viper.BindPFlag("number", GetCmd.PersistentFlags().Lookup("number"))
	viper.BindPFlag("filter", GetCmd.PersistentFlags().Lookup("filter"))

	cmd.RootCmd.AddCommand(GetCmd)
}