package utils

import (
	"fmt"
	"log"
	"strings"

	"bitbucket.org/level27/lvl/types"
)

//------------------------------------------------- APP (GET / CREATE  / UPDATE / DELETE / DESCRIBE)-------------------------------------------------

// Gets an app from the API
func (c *Client) App(id int) types.App {
	var app struct {
		App types.App `json:"app"`
	}

	endpoint := fmt.Sprintf("apps/%d", id)
	err := c.invokeAPI("GET", endpoint, nil, &app)
	AssertApiError(err, "app")

	return app.App
}

// Gets a list of apps from the API
func (c *Client) Apps(getParams types.CommonGetParams) []types.App {
	var apps struct {
		Apps []types.App `json:"apps"`
	}

	endpoint := fmt.Sprintf("apps?%s", formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &apps)
	AssertApiError(err, "app")

	return apps.Apps
}

// ---- CREATE NEW APP
func (c *Client) AppCreate(req types.AppPostRequest) types.App {
	var app struct {
		Data types.App `json:"app"`
	}
	endpoint := "apps"
	err := c.invokeAPI("POST", endpoint, req, &app)

	AssertApiError(err, "apps")

	return app.Data
}

// ---- DELETE APP
func (c *Client) AppDelete(appId int, isConfirmed bool) {
	endpoint := fmt.Sprintf("apps/%v", appId)

	if isConfirmed {
		err := c.invokeAPI("DELETE", endpoint, nil, nil)
		AssertApiError(err, "Apps")
	} else {
		var userResponse string
		// ask user for confirmation on deleting the check
		question := fmt.Sprintf("Are you sure you want to delete the app with ID: %v? Please type [y]es or [n]o: ", appId)
		fmt.Print(question)
		//reading user response
		_, err := fmt.Scan(&userResponse)
		if err != nil {
			log.Fatal(err)
		}
		// check if user confirmed the deletion or not
		switch strings.ToLower(userResponse) {
		case "y", "yes":
			err := c.invokeAPI("DELETE", endpoint, nil, nil)
			AssertApiError(err, "Apps")
		case "n", "no":
			log.Printf("Delete canceled for app: %v", appId)
		default:
			log.Println("Please make sure you type (y)es or (n)o and press enter to confirm:")

			c.AppDelete(appId, false)
		}
	}
}
