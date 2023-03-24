package cmd

import (
	"fmt"
	"strconv"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

func init() {
	systemCmd.AddCommand(systemNetworkCmd)

	systemNetworkCmd.AddCommand(systemNetworkGetCmd)

	systemNetworkCmd.AddCommand(systemNetworkDescribeCmd)

	systemNetworkCmd.AddCommand(systemNetworkAddCmd)

	systemNetworkCmd.AddCommand(systemNetworkRemoveCmd)
}

func resolveSystemHasNetwork(systemID l27.IntID, arg string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.LookupSystemHasNetworks(systemID, arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"system network",
		func(app l27.SystemHasNetwork) string { return fmt.Sprintf("%s (%d)", app.Network.Description, app.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, err
}

// NETWORKS

var systemNetworkCmd = &cobra.Command{
	Use: "network",
}

var systemNetworkGetCmd = &cobra.Command{
	Use:   "get [system]",
	Short: "Get list of networks on a system",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		system, err := Level27Client.SystemGetSingle(systemID)
		if err != nil {
			return err
		}

		outputFormatTableFuncs(system.Networks, []string{"ID", "Network ID", "Type", "Name", "MAC", "IPs"}, []interface{}{"ID", "NetworkID", func(net l27.SystemNetwork) string {
			if net.NetPublic {
				return "public"
			}
			if net.NetCustomer {
				return "customer"
			}
			if net.NetInternal {
				return "internal"
			}
			return ""
		}, "Name", "Mac", func(net l27.SystemNetwork) string {
			return strconv.Itoa(len(net.Ips))
		}})

		return nil
	},
}

var systemNetworkDescribeCmd = &cobra.Command{
	Use:   "describe [system]",
	Short: "Display detailed information about all networks on a system",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		system, err := Level27Client.SystemGetSingle(systemID)
		if err != nil {
			return err
		}

		networks, err := Level27Client.SystemGetHasNetworks(systemID)
		if err != nil {
			return err
		}

		outputFormatTemplate(DescribeSystemNetworks{
			Networks:    system.Networks,
			HasNetworks: networks,
		}, "templates/systemNetworks.tmpl")

		return nil
	},
}

var systemNetworkAddCmd = &cobra.Command{
	Use:   "add [system] [network]",
	Short: "Add a network to a system",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		networkID, err := resolveNetwork(args[1])
		if err != nil {
			return err
		}

		network, err := Level27Client.SystemAddHasNetwork(systemID, networkID)
		if err != nil {
			return err
		}

		outputFormatTemplate(network, "templates/entities/systemNetwork/add.tmpl")
		return nil
	},
}

var systemNetworkRemoveCmd = &cobra.Command{
	Use:   "remove [system] [network]",
	Short: "Remove a network from a system",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		networkID, err := resolveSystemHasNetwork(systemID, args[1])
		if err != nil {
			return err
		}

		err = Level27Client.SystemRemoveHasNetwork(systemID, networkID)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/systemNetwork/remove.tmpl")
		return nil
	},
}

type DescribeSystemNetworks struct {
	Networks    []l27.SystemNetwork    `json:"networks"`
	HasNetworks []l27.SystemHasNetwork `json:"hasNetworks"`
}
