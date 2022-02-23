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
func (c *Client) SystemCreate(args []string, req types.SystemPost) {
	

	var System struct {
		Data types.System `json:"system"`
	}

	

	err := c.invokeAPI("POST", "systems", req, &System)
	AssertApiError(err, "SystemCreate")

	log.Printf("System created! [Fullname: '%v' , ID: '%v']", System.Data.Name, System.Data.Id)

}
