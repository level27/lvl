package utils

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"bitbucket.org/level27/lvl/types"
)

//------------------------------------------------- Resolve functions -------------------------------------------------

// GET appID based on name
func (c *Client) AppLookup(name string) []types.App {
	results := []types.App{}
	apps := c.Apps(types.CommonGetParams{Filter: name})
	for _, app := range apps {
		if app.Name == name {
			results = append(results, app)
		}
	}

	return results
}

// GET componentId based on name
func (c *Client) AppComponentLookup(appId int, name string) []types.AppComponent{
	results := []types.AppComponent{}
	components := c.AppComponentsGet(appId, types.CommonGetParams{Filter: name})
	for _, component := range components {
		if component.Name == name {
			results = append(results, component)
		}
	}

	return results
}

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


// APP SSL CERTIFICATES

// GET /apps/{appID}/sslcertificates
func (c *Client) AppSslCertificatesGetList(appID int, sslType string, status string, get types.CommonGetParams) []types.AppSslCertificate {
	var response struct {
		SslCertificates []types.AppSslCertificate `json:"sslCertificates"`
	}

	endpoint := fmt.Sprintf(
		"apps/%d/sslcertificates?sslType=%s&status=%s&%s",
		appID,
		url.QueryEscape(sslType),
		url.QueryEscape(status),
		formatCommonGetParams(get))

	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "AppSslCertificatesGetList")

	return response.SslCertificates
}

// GET /apps/{appID}/sslcertificates/{sslCertificateID}
func (c *Client) AppSslCertificatesGetSingle(appID int, sslCertificateID int) types.AppSslCertificate {
	var response struct {
		SslCertificate types.AppSslCertificate `json:"sslCertificate"`
	}

	endpoint := fmt.Sprintf("apps/%d/sslcertificates/%d", appID, sslCertificateID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "AppSslCertificatesGetSingle")

	return response.SslCertificate
}

// POST /apps/{appID}/sslcertificates
func (c *Client) AppSslCertificatesCreate(appID int, create types.AppSslCertificateCreate) types.AppSslCertificate {
	var response struct {
		SslCertificate types.AppSslCertificate `json:"sslCertificate"`
	}

	endpoint := fmt.Sprintf("apps/%d/sslcertificates", appID)
	err := c.invokeAPI("POST", endpoint, create, &response)
	AssertApiError(err, "AppSslCertificatesCreate")

	return response.SslCertificate
}

// POST /apps/{appID}/sslcertificates (variant for sslType == "own")
func (c *Client) AppSslCertificatesCreateOwn(appID int, create types.AppSslCertificateCreateOwn) types.AppSslCertificate {
	var response struct {
		SslCertificate types.AppSslCertificate `json:"sslCertificate"`
	}

	endpoint := fmt.Sprintf("apps/%d/sslcertificates", appID)
	err := c.invokeAPI("POST", endpoint, create, &response)
	AssertApiError(err, "AppSslCertificatesCreate")

	return response.SslCertificate
}

// DELETE /apps/{appID}/sslcertificates/{sslCertificateID}
func (c *Client) AppSslCertificatesDelete(appID int, sslCertificateID int) {
	endpoint := fmt.Sprintf("apps/%d/sslcertificates/%d", appID, sslCertificateID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)
	AssertApiError(err, "AppSslCertificatesDelete")
}

// PUT /apps/{appID}/sslcertificates/{sslCertificateID}
func (c *Client) AppSslCertificatesUpdate(appID int, sslCertificateID int, data map[string]interface{}) {
	endpoint := fmt.Sprintf("apps/%d/sslcertificates/%d", appID, sslCertificateID)
	err := c.invokeAPI("PUT", endpoint, data, nil)
	AssertApiError(err, "AppSslCertificatesUpdate")
}

// Try to find an SSL certificate on an app by name.
func (c *Client) AppSslCertificatesLookup(appID int, name string) []types.AppSslCertificate {
	results := []types.AppSslCertificate{}
	apps := c.AppSslCertificatesGetList(appID, "", "", types.CommonGetParams{Filter: name})
	for _, cert := range apps {
		if cert.Name == name {
			results = append(results, cert)
		}
	}

	return results
}

// POST /apps/{appID}/sslcertificates/{sslCertificateID}/actions
func (c *Client) AppSslCertificatesActions(appID int, sslCertificateID int, actionType string) {
	var request struct {
		Type string `json:"type"`
	}

	request.Type = actionType

	endpoint := fmt.Sprintf("apps/%d/sslcertificates/%d/actions", appID, sslCertificateID)
	err := c.invokeAPI("POST", endpoint, request, nil)
	AssertApiError(err, "AppSslCertificatesActions")
}

// POST /apps/{appID}/sslcertificates/{sslCertificateID}/fix
func (c *Client) AppSslCertificatesFix(appID int, sslCertificateID int) types.AppSslCertificate {
	var response struct {
		SslCertificate types.AppSslCertificate `json:"sslCertificate"`
	}

	endpoint := fmt.Sprintf("apps/%d/sslcertificates/%d/fix", appID, sslCertificateID)
	err := c.invokeAPI("POST", endpoint, nil, &response)
	AssertApiError(err, "AppSslCertificatesFix")

	return response.SslCertificate
}

// GET /apps/{appID}/sslcertificates/{sslCertificateID}/key
func (c *Client) AppSslCertificatesKey(appID int, sslCertificateID int) types.AppSslcertificateKey {
	var response types.AppSslcertificateKey

	endpoint := fmt.Sprintf("apps/%d/sslcertificates/%d/key", appID, sslCertificateID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "AppSslCertificatesKey")

	return response
}
//------------------------------------------------- APP COMPONENTS (GET / DESCRIBE / CREATE)-------------------------------------------------

// ---- GET LIST OF COMPONENTS
func (c *Client) AppComponentsGet(appid int, getParams types.CommonGetParams) []types.AppComponent {
	var components struct {
		Data []types.AppComponent `json:"components"`
	}

	endpoint := fmt.Sprintf("apps/%v/components?%v", appid, formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &components)
	AssertApiError(err, "app")

	return components.Data
}

// ---- DESCRIBE COMPONENT (GET SINGLE COMPONENT)
func (c *Client) AppComponentGetSingle(appId int, id int) types.AppComponent {
	var component struct {
		Data types.AppComponent `json:"component"`
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
