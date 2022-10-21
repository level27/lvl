package cmd

import (
	"fmt"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

func resolveNetwork(arg string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.LookupNetwork(arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"network",
		func(app l27.Network) string { return fmt.Sprintf("%s (%d)", app.Name, app.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, err
}

func init() {
	RootCmd.AddCommand(networkCmd)

	networkCmd.AddCommand(networkGetCmd)
	addCommonGetFlags(networkCmd)
}

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Commands for managing networks",
}

var networkGetCmd = &cobra.Command{
	Use: "get",

	RunE: func(cmd *cobra.Command, args []string) error {
		networks, err := Level27Client.GetNetworks(optGetParameters)
		if err != nil {
			return err
		}

		outputFormatTableFuncs(networks, []string{"ID", "Type", "Name", "VLAN", "Organisation", "Zone"}, []interface{}{"ID", func(net l27.Network) string {
			if net.Public {
				return "public"
			}
			if net.Customer {
				return "customer"
			}
			if net.Internal {
				return "internal"
			}
			return ""
		}, "Name", "Vlan", "Organisation.Name", "Zone.Name"})
		return nil
	},
}

func ipv4IntToString(ipv4 int) string {
	a := (ipv4 >> 24) & 0xFF
	b := (ipv4 >> 16) & 0xFF
	c := (ipv4 >> 8) & 0xFF
	d := (ipv4 >> 0) & 0xFF

	return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
}
