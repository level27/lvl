package utils

import (
	"fmt"
	"log"

	"bitbucket.org/level27/lvl/types"
)

// --------------------------- TOPLEVEL SYSTEM ACTIONS (GET / POST) ------------------------------------
//------------------ GET

// returning a list of all current systems [lvl system get]
func (c *Client) SystemGetList(getParams types.CommonGetParams) []types.System {

	//creating an array of systems.
	var systems struct {
		Data []types.System `json:"systems"`
	}

	//creating endpoint
	endpoint := fmt.Sprintf("systems?%s", formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &systems)
	AssertApiError(err, "Systems")
	//returning result as system type
	return systems.Data

}

// Returning a single system by its ID
// this is not for a describe.
func (c *Client) SystemGetSingle(id int) types.System {
	var system struct {
		Data types.System `json:"system"`
	}
	endpoint := fmt.Sprintf("systems/%v", id)
	err := c.invokeAPI("GET", endpoint, nil, &system)

	AssertApiError(err, "System")
	return system.Data

}

//----------------- POST

// CREATE SYSTEM [lvl system create <parmeters>]
func (c *Client) SystemCreate( req types.SystemPost) {
	

	var System struct {
		Data types.System `json:"system"`
	}

	

	err := c.invokeAPI("POST", "systems", req, &System)
	AssertApiError(err, "SystemCreate")

	log.Printf("System created! [Fullname: '%v' , ID: '%v']", System.Data.Name, System.Data.Id)

}

// --------------------------- SYSTEM/CHECKS TOPLEVEL (GET / POST) ------------------------------------
// ------------- GET CHECKS
func (c *Client) SystemCheckGetList(systemId int, getParams types.CommonGetParams) []types.SystemCheck {

	//creating an array of systems.
	var systemChecks struct {
		Data []types.SystemCheck `json:"checks"`
	}

	//creating endpoint
	endpoint := fmt.Sprintf("systems/%v/checks?%s", systemId, formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &systemChecks)
	AssertApiError(err, "Systems")
	//returning result as system check type
	return systemChecks.Data

}

// ------------- CREATE A CHECK
func (c *Client) SystemCheckCreate(systemId int, req types.SystemCheckRequest ) {
	var SystemCheck struct {
		Data types.SystemCheck `json:"check"`
	}
	endpoint := fmt.Sprintf("systems/%v/checks", systemId)
	err := c.invokeAPI("POST", endpoint, req, &SystemCheck )

	AssertApiError(err, "System checks")
	log.Printf("System check created! [Checktype: '%v' , ID: '%v']", SystemCheck.Data.CheckType, SystemCheck.Data.Id)
}

// --------------------------- SYSTEM/COOKBOOKS TOPLEVEL (GET / POST) ------------------------------------
// ------------- GET COOKBOOK

func (c *Client) SystemCookbookGetList(systemId int) []types.Cookbook{
	// creating array of cookbooks to return
	var systemCookbooks struct{
		Data []types.Cookbook `json:"cookbooks"`
	}

	endpoint := fmt.Sprintf("systems/%v/cookbooks", systemId)
	err := c.invokeAPI("GET", endpoint, nil, &systemCookbooks)

	AssertApiError(err, "cookbooks")

	return systemCookbooks.Data

}