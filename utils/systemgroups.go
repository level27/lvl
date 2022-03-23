package utils

import (
	"fmt"

	"bitbucket.org/level27/lvl/types"
)

//------------------------------------------------- SYSTEMSGROUPS (GET / CREATE  / UPDATE / DELETE)-------------------------------------------------


// ---------------- GET 
func (c *Client) SystemgroupsGet(optParameters types.CommonGetParams) []types.Systemgroup {
	// var to store API response
	var systemgroups struct {
		Data []types.Systemgroup `json:"systemgroups"`
	}
	endpoint := fmt.Sprintf("systemgroups?%v",formatCommonGetParams(optParameters))
	err := c.invokeAPI("GET", endpoint, nil, &systemgroups)
	AssertApiError(err, "systemgroups")

	return systemgroups.Data
}


// ---------------- CREATE
func (c *Client) SystemgroupsCreate(req types.SystemgroupRequest) types.Systemgroup{
	// var to store API response
	var systemgroup struct {
		Data types.Systemgroup `json:"systemgroup"`
	}

	endpoint := "systemgroups"
	err := c.invokeAPI("POST", endpoint, req, &systemgroup)
	AssertApiError(err, "systemgroup")

	return systemgroup.Data
}