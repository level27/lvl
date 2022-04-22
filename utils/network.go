package utils

import (
	"fmt"
	"net"
	"strconv"

	"bitbucket.org/level27/lvl/types"
)

func (c *Client) GetNetworks(get types.CommonGetParams) []types.Network {
	var networks struct {
		Networks []types.Network `json:"networks"`
	}

	endpoint := fmt.Sprintf("networks?%s", formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &networks)

	AssertApiError(err, "Networks")
	return networks.Networks
}

func (c *Client) GetNetwork(id int) types.Network {
	var network struct {
		Network types.Network `json:"network"`
	}

	endpoint := fmt.Sprintf("network/%d", id)
	err := c.invokeAPI("GET", endpoint, nil, &network)

	AssertApiError(err, "Network")
	return network.Network
}

func (c *Client) LookupNetwork(name string) []types.Network {
	results := []types.Network{}
	networks := c.GetNetworks(types.CommonGetParams{Filter: name})
	for _, net := range networks {
		if net.Name == name {
			results = append(results, net)
		}
	}

	return results
}

func (c *Client) NetworkLocate(networkID int) types.NetworkLocate {
	var response types.NetworkLocate

	endpoint := fmt.Sprintf("networks/%d/locate", networkID)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	AssertApiError(err, "NetworkLocate")
	return response
}

func Ipv4IntToString(ipv4 int) string {
	a := (ipv4 >> 24) & 0xFF
	b := (ipv4 >> 16) & 0xFF
	c := (ipv4 >> 8) & 0xFF
	d := (ipv4 >> 0) & 0xFF

	return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
}

func Ipv4StringIntToString(ipv4 string) string {
	i, err := strconv.Atoi(ipv4)
	if err != nil {
		return ""
	}

	return Ipv4IntToString(i)
}

func IpsEqual(a string, b string) bool {
	ipA := net.ParseIP(a)
	ipB := net.ParseIP(b)

	return ipA.Equal(ipB)
}