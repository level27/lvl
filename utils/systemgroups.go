package utils

import (
	"fmt"

	"bitbucket.org/level27/lvl/types"
)

func (c *Client) SystemgroupsGet(optParameters types.CommonGetParams) []types.Systemgroup {
	var systemgroups struct {
		Data []types.Systemgroup `json:"systemgroups"`
	}
	endpoint := fmt.Sprintf("systemgroups?%v",formatCommonGetParams(optParameters))
	err := c.invokeAPI("GET", endpoint, nil, &systemgroups)
	AssertApiError(err, "systemgroups")

	return systemgroups.Data
}
