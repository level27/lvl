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

	networkCmd.AddCommand(networkZoneCmd)

	networkZoneCmd.AddCommand(networkZoneAddCmd)
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

var networkZoneCmd = &cobra.Command{
	Use:   "zone",
	Short: "Manage zones that networks available in",
}

var networkZoneAddCmd = &cobra.Command{
	Use:   "add <network> <zone>",
	Short: "Add a zone to a network",
	Args:  cobra.ExactArgs(2),

	RunE: func(cmd *cobra.Command, args []string) error {
		networkId, err := resolveNetwork(args[0])
		if err != nil {
			return err
		}

		zoneId, _, err := resolveZoneRegion(args[1])
		if err != nil {
			return err
		}

		network, err := Level27Client.GetNetwork(networkId)
		if err != nil {
			return err
		}

		put := l27.NetworkPutRequest{
			Remarks:     network.Remarks,
			Description: network.Description,
		}

		for _, zone := range network.Zones {
			if zone.ID == zoneId {
				return fmt.Errorf("zone already has network")
			}

			put.Zones = append(put.Zones, zone.ID)
		}

		put.Zones = append(put.Zones, zoneId)

		err = Level27Client.NetworkUpdate(networkId, put)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/network/update.tmpl")

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
