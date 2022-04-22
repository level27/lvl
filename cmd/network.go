package cmd

import (
	"fmt"
	"strconv"

	"bitbucket.org/level27/lvl/types"
	"github.com/spf13/cobra"
)

func resolveNetwork(arg string) int {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id
	}


	return resolveShared(
		Level27Client.LookupNetwork(arg),
		arg,
		"network",
		func (app types.Network) string { return fmt.Sprintf("%s (%d)", app.Name, app.ID) }).ID
}

func init() {
	RootCmd.AddCommand(networkCmd)

	networkCmd.AddCommand(networkGetCmd)
	addCommonGetFlags(networkCmd)
}

var networkCmd = &cobra.Command{
	Use: "network",
}

var networkGetCmd = &cobra.Command{
	Use: "get",

	Run: func(cmd *cobra.Command, args []string) {
		networks := Level27Client.GetNetworks(optGetParameters)
		outputFormatTableFuncs(networks, []string{"ID", "Type", "Name", "VLAN", "Organisation", "Zone"}, []interface{}{"ID", func(net types.Network) string {
			if net.Public { return "public" }
			if net.Customer { return "customer" }
			if net.Internal { return "internal" }
			return ""
		}, "Name", "Vlan", "Organisation.Name", "Zone.Name"})
	},
}