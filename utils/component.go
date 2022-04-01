package utils

import (
	"fmt"

	"bitbucket.org/level27/lvl/types"
)

// Gets a list of components of a specified category (and optionally type) from the API.
func (c *Client) Components(category string, cType string, getParams types.CommonGetParams) []types.AppComponent2 {
	var components struct {
		Components []types.AppComponent2 `json:"components"`
	}

	endpoint := fmt.Sprintf("appcomponents/%s?%s&type=%s", category, formatCommonGetParams(getParams), cType)
	err := c.invokeAPI("GET", endpoint, nil, &components)
	AssertApiError(err, "component")

	return components.Components
}