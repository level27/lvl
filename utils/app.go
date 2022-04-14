package utils

import (
	"fmt"
	"log"
	"strings"

	"bitbucket.org/level27/lvl/types"
)

//------------------------------------------------- APP MAIN SUBCOMMANDS (GET / CREATE  / UPDATE / DELETE / DESCRIBE)-------------------------------------------------
// #region APP MAIN SUBCOMMANDS (GET / CREATE  / UPDATE / DELETE / DESCRIBE)

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

// ---- UPDATE APP
func (c *Client) AppUpdate(appId int, req types.AppPutRequest) {
	endpoint := fmt.Sprintf("apps/%v", appId)
	err := c.invokeAPI("PUT", endpoint, req, nil)
	AssertApiError(err, "Apps")
}

// #endregion

//------------------------------------------------- APP ACTIONS (ACTIVATE / DEACTIVATE)-------------------------------------------------
// ---- ACTION (ACTIVATE OR DEACTIVATE) ON AN APP
func (c *Client) AppAction(appId int, action string) {
	request := types.AppActionRequest{
		Type: action,
	}
	endpoint := fmt.Sprintf("apps/%v/actions", appId)
	err := c.invokeAPI("POST", endpoint, request, nil)
	AssertApiError(err, "app")
}

//------------------------------------------------- APP COMPONENTS (GET / DESCRIBE / CREATE)-------------------------------------------------

// ---- GET LIST OF COMPONENTS
func (c *Client) AppComponentsGet(appid int, getParams types.CommonGetParams) []types.AppComponent2 {
	var components struct {
		Data []types.AppComponent2 `json:"components"`
	}

	endpoint := fmt.Sprintf("apps/%v/components?%v", appid, formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &components)
	AssertApiError(err, "app")

	log.Print("hallo")
	return components.Data
}

// ---- DESCRIBE COMPONENT (GET SINGLE COMPONENT)
func (c *Client) AppComponentGetSingle(appId int, id int) types.AppComponent2 {
	var component struct {
		Data types.AppComponent2 `json:"component"`
	}

	endpoint := fmt.Sprintf("apps/%d/components/%v", appId, id)
	err := c.invokeAPI("GET", endpoint, nil, &component)
	AssertApiError(err, "app")
	return component.Data
}

//------------------------------------------------- APP COMPONENTS HELPERS (CATEGORY )-------------------------------------------------
// ---- GET LIST OFF APPCOMPONENTTYPES
func (c *Client) AppComponenttypesGet() types.Appcomponenttype {
	var componenttypes struct {
		Data types.Appcomponenttype `json:"appcomponenttypes"`
	}

	endpoint := "appcomponenttypes"
	err := c.invokeAPI("GET", endpoint, nil, &componenttypes)
	AssertApiError(err, "appcomponent")
	return componenttypes.Data
}


func (c *Client) AppCertificateGet(appId int) []types.SslCertificate{
	var certificates struct{
		Data []types.SslCertificate `json:"sslCertificates"`
	}

	endpoint := fmt.Sprintf("apps/%v/sslcertificates", appId)
	err := c.invokeAPI("GET", endpoint, nil, &certificates)
	AssertApiError(err, "appCertificate")
	return certificates.Data
}


// ---- ADD SSL CERTIFICATE
func (c *Client)AppCertificateAdd(appId int, req interface{}) types.SslCertificate{
	var certificate struct{
		Data types.SslCertificate `json:"sslCertificate"`
	}
	endpoint := fmt.Sprintf("apps/%v/sslcertificates", appId)
	err := c.invokeAPI("POST", endpoint, req, &certificate)
	AssertApiError(err, "appCertificate")
	return certificate.Data
}