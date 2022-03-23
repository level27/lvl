package utils

import (
	"fmt"

	"bitbucket.org/level27/lvl/types"
)

//------------------------------------------------- SYSTEMSGROUPS (GET / CREATE  / UPDATE / DELETE)-------------------------------------------------

// ---------------- GET SINGLE (describe)
func (c *Client) SystemgroupsgetSingle(systemgroupId int) types.Systemgroup{
	// var to store API response
	var systemgroup struct {
		Data types.Systemgroup `json:"systemgroup"`
	}

	endpoint := fmt.Sprintf("systemgroups/%v", systemgroupId)
	err := c.invokeAPI("GET", endpoint, nil, &systemgroup)
	AssertApiError(err, "systemgroups")

	return systemgroup.Data

}

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


// ---------------- UPDATE 
func (c *Client) SystemgroupsUpdate(systemgroupId int, req types.SystemgroupRequest){
	endpoint := fmt.Sprintf("systemgroups/%v", systemgroupId)
	err := c.invokeAPI("PUT", endpoint, req, nil)
	AssertApiError(err, "systemgroup")
}