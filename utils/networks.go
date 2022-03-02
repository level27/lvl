package utils

import (
	"fmt"

	"bitbucket.org/level27/lvl/types"
)

func (c *Client) GetNetworks(get types.CommonGetParams) []types.Network {
	var networks struct {
		Networks []types.Network `json:"network"`
	}

	endpoint := fmt.Sprintf("networks?%s", formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &networks)

	AssertApiError(err, "Networks")
	return networks.Networks
}

func Ipv4IntToString(ipv4 int) string {
	a := (ipv4 >> 24) & 0xFF
	b := (ipv4 >> 16) & 0xFF
	c := (ipv4 >> 8) & 0xFF
	d := (ipv4 >> 0) & 0xFF

	return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
}