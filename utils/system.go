package utils

import (
	"fmt"

	"bitbucket.org/level27/lvl/types"
)

// returning a list of all current systems [lvl system get]
func (c *Client) SystemGetList() []types.System {

	//creating an array of systems.
	var systems struct {
		Data []types.System `json:"systems"`
	}

	endpoint := "systems"
	err := c.invokeAPI("GET", endpoint, nil, &systems)
	AssertApiError(err, "Systems")
	return systems.Data

}

// Returning a single system by its ID
func (c *Client) SystemGetSingle(id int) types.System{
	var system types.System
	endpoint := fmt.Sprintf("system?%v", id)
	err := c.invokeAPI("GET", endpoint, nil, &system)

	AssertApiError(err, "System")
	return system

}