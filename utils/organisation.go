package utils

import (
	"fmt"
	"net/url"

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
func (c *Client) Organisations(filter string, number int) types.Organisations {
	var organisations types.Organisations

	endpoint := fmt.Sprintf("organisations?limit=%d&filter=%s", number, url.QueryEscape(filter))
	err := c.invokeAPI("GET", endpoint, nil, &organisations)
	AssertApiError(err)

	return organisations
}