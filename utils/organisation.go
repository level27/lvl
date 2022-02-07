package utils

import (
	"fmt"

	"bitbucket.org/level27/lvl/types"
)

//Organisation gets a system from the API
func (c *Client) Organisation(method string, id interface{}, data interface{}) types.Organisation {
	var organisation types.Organisation

	switch method {
	case "GET":
		endpoint := fmt.Sprintf("organisations/%s", id)
		c.invokeAPI("GET", endpoint, nil, &organisation)
	case "CREATE":
		endpoint := "organisations"
		c.invokeAPI("POST", endpoint, data, &organisation)
	case "UPDATE":
		endpoint := fmt.Sprintf("organisations/%s", id)
		c.invokeAPI("PUT", endpoint, data, &organisation)
	case "DELETE":
		endpoint := fmt.Sprintf("organisations/%s", id)
		c.invokeAPI("DELETE", endpoint, nil, nil)
	}

	return organisation
}

//Organisation gets a organisation from the API
func (c *Client) Organisations(getParams types.CommonGetParams) types.Organisations {
	var organisations types.Organisations

	endpoint := fmt.Sprintf("organisations?%s", formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &organisations)
	AssertApiError(err, "organisation")

	return organisations
}