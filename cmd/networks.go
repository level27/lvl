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

	network := Level27Client.LookupNetwork(arg)
	if network == nil {
		cobra.CheckErr(fmt.Sprintf("Unable to find network: %s", arg))
		return 0
	}

	return network.ID
}

func init() {
	RootCmd.AddCommand(networksCmd)

	networksCmd.AddCommand(networksGetCmd)
	addCommonGetFlags(networksCmd)
}

var networksCmd = &cobra.Command{
	Use: "networks",
}

var networksGetCmd = &cobra.Command{
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