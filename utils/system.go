package utils

import (
	"fmt"

	"bitbucket.org/level27/lvl/types"
)

func (c *Client) SystemGetList() []types.SystemGet {

	var systems struct {
		Data []types.SystemGet `json:"systems"`
	}
	endpoint := fmt.Sprint("systems")
	err := c.invokeAPI("GET", endpoint, nil, &systems)
	AssertApiError(err, "Systems")
	return systems.Data

}


