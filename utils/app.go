package utils

import (
	"fmt"

	"bitbucket.org/level27/lvl/types"
)

// Gets an app from the API
func (c *Client) App(method string, id interface{}, data interface{}) types.App {
	var app types.App

	switch method {
	case "GET":
		endpoint := fmt.Sprintf("apps/%s", id)
		c.invokeAPI("GET", endpoint, nil, &app)
	case "CREATE":
		endpoint := "apps"
		c.invokeAPI("POST", endpoint, data, &app)
	case "UPDATE":
		endpoint := fmt.Sprintf("apps/%s", id)
		c.invokeAPI("PUT", endpoint, data, &app)
	case "DELETE":
		endpoint := fmt.Sprintf("apps/%s", id)
		c.invokeAPI("DELETE", endpoint, nil, nil)
	}

	return app
}

func (c *Client) Apps(filter string, number string) types.Apps {
	var apps types.Apps

	endpoint := "apps"
	err := c.invokeAPI("GET", endpoint, nil, &apps)
	AssertApiError(err)

	return apps
}

