package utils

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
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

// ---- DELETE COMPONENT
func (c *Client) AppComponentsDelete(appId int, componentId int, isDeleteConfirmed bool) {
	endpoint := fmt.Sprintf("apps/%v/components/%v", appId, componentId)

	if isDeleteConfirmed {
		err := c.invokeAPI("DELETE", endpoint, nil, nil)
		AssertApiError(err, "appcomponent")
	} else {
		var userResponse string
		// ask user for confirmation on deleting the check
		question := fmt.Sprintf("Are you sure you want to delete the appcomponent with ID: %v? Please type [y]es or [n]o: ", componentId)
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
			AssertApiError(err, "appcomponent")
		case "n", "no":
			log.Printf("Delete canceled for appcomponent: %v", componentId)
		default:
			log.Println("Please make sure you type (y)es or (n)o and press enter to confirm:")

			c.AppComponentsDelete(appId, componentId, false)
		}
	}

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

//------------------------------------------------- APP SSL CERTIFICATES (GET/ ADD/ DELETE )-------------------------------------------------

// ---- GET LIST OF SSL CERTIFICATES
func (c *Client) AppCertificateGet(appId int) []types.SslCertificate {
	var certificates struct {
		Data []types.SslCertificate `json:"sslCertificates"`
	}

	endpoint := fmt.Sprintf("apps/%v/sslcertificates", appId)
	err := c.invokeAPI("GET", endpoint, nil, &certificates)
	AssertApiError(err, "appCertificate")
	return certificates.Data
}

// ---- ADD SSL CERTIFICATE
func (c *Client) AppCertificateAdd(appId int, req interface{}) types.SslCertificate {
	var certificate struct {
		Data types.SslCertificate `json:"sslCertificate"`
	}
	endpoint := fmt.Sprintf("apps/%v/sslcertificates", appId)
	err := c.invokeAPI("POST", endpoint, req, &certificate)
	AssertApiError(err, "appCertificate")
	return certificate.Data
}

// ---- DELETE SSL CERTIFICATE
func (c *Client) AppCertificateDelete(appId int, certificateId int, isDeleteConfirmed bool) {

	// when confirmation flag is set, delete check without confirmation question
	if isDeleteConfirmed {
		endpoint := fmt.Sprintf("apps/%v/sslcertificates/%v", appId, certificateId)
		err := c.invokeAPI("DELETE", endpoint, nil, nil)
		AssertApiError(err, "appCertificate")
	} else {
		var userResponse string
		// ask user for confirmation on deleting the check
		question := fmt.Sprintf("Are you sure you want to delete the ssl certificate with ID: %v? Please type [y]es or [n]o: ", certificateId)
		fmt.Print(question)
		//reading user response
		_, err := fmt.Scan(&userResponse)
		if err != nil {
			log.Fatal(err)
		}
		// check if user confirmed the deletion of the check or not
		switch strings.ToLower(userResponse) {
		case "y", "yes":
			endpoint := fmt.Sprintf("apps/%v/sslcertificates/%v", appId, certificateId)
			err := c.invokeAPI("DELETE", endpoint, nil, nil)
			AssertApiError(err, "appCertificate")
		case "n", "no":
			log.Printf("Delete canceled for ssl certificate: %v", certificateId)
		default:
			log.Println("Please make sure you type (y)es or (n)o and press enter to confirm:")

			c.AppCertificateDelete(appId, certificateId, false)
		}
	}
}

// ---- FIX SSL CERTIFICATE
func (c *Client) AppCertificateFix(appID int, certificateId int) {
	endpoint := fmt.Sprintf("apps/%v/sslcertificates/%v/fix", appID, certificateId)
	err := c.invokeAPI("POST", endpoint, nil, nil)
	AssertApiError(err, "appCertificate")
}

// ---- GET PRIVATE KEY (TYPE 'OWN' CERTIFICATE)
func (c *Client) AppCertificateKey(appId int, certificateId int) {
	var key struct {
		Data string `json:"sslKey"`
	}
	endpoint := fmt.Sprintf("apps/%v/sslcertificates/%v/key", appId, certificateId)
	err := c.invokeAPI("GET", endpoint, nil, &key)
	AssertApiError(err, "appCertificate")

	fmt.Print(key.Data)
}

//------------------------------------------------- APP SSL CERTIFICATES (ACTIONS)-------------------------------------------------
// ACTION RETRY / VALIDATECHALLENGE (SSL)
func (c *Client) AppCertificateAction(appId int, certificateId int, actionType string) {
	// create request data
	request := types.AppSslCertificateActionRequest{
		Type: actionType,
	}

	endpoint := fmt.Sprintf("apps/%v/sslcertificates/%v/actions", appId, certificateId)
	err := c.invokeAPI("POST", endpoint, request, nil)
	AssertApiError(err, "appCertificate")

	log.Println("Action retry sent to certificate.")
}

//-------------------------------------------------  APP RESTORE (GET / DESCRIBE / CREATE / UPDATE / DELETE / DOWNLOAD) -------------------------------------------------

// ---- GET LIST OF APP RESTORES
func (c *Client) AppComponentRestoresGet(appId int) []types.AppComponentRestore {
	var restores struct {
		Data []types.AppComponentRestore `json:"restores"`
	}

	endpoint := fmt.Sprintf("apps/%v/restores", appId)
	err := c.invokeAPI("GET", endpoint, nil, &restores)
	AssertApiError(err, "appRestore")
	return restores.Data
}

// ---- CREATE NEW RESTORE
func (c *Client) AppComponentRestoreCreate(appId int, req types.AppComponentRestoreRequest) types.AppComponentRestore {
	var restore struct {
		Data types.AppComponentRestore `json:"restore"`
	}
	endpoint := fmt.Sprintf("apps/%v/restores", appId)
	err := c.invokeAPI("POST", endpoint, req, &restore)
	AssertApiError(err, "appRestores")

	return restore.Data
}

// ---- DELETE RESTORE
func (c *Client) AppComponentRestoresDelete(appId int, restoreId int, isDeleteConfirmed bool) {
	endpoint := fmt.Sprintf("apps/%v/restores/%v", appId, restoreId)

	// when confirmation flag is set, delete check without confirmation question
	if isDeleteConfirmed {
		err := c.invokeAPI("DELETE", endpoint, nil, nil)
		AssertApiError(err, "appRestore")
		log.Print("Restore succesfully deleted.")
	} else {
		var userResponse string
		// ask user for confirmation on deleting the check
		question := fmt.Sprintf("Are you sure you want to delete the restore with ID: %v? Please type [y]es or [n]o: ", restoreId)
		fmt.Print(question)
		//reading user response
		_, err := fmt.Scan(&userResponse)
		if err != nil {
			log.Fatal(err)
		}
		// check if user confirmed the deletion of the check or not
		switch strings.ToLower(userResponse) {
		case "y", "yes":
			err := c.invokeAPI("DELETE", endpoint, nil, nil)
			AssertApiError(err, "appRestore")
			log.Print("Restore succesfully deleted.")
		case "n", "no":
			log.Printf("Delete canceled for app restore: %v", restoreId)
		default:
			log.Println("Please make sure you type (y)es or (n)o and press enter to confirm:")

			c.AppCertificateDelete(appId, restoreId, false)
		}
	}
}


// ---- DOWNLOAD RESTORE FILE
func (c *Client) AppComponentRestoreDownload(appId int, restoreId int, filename string){
	endpoint := fmt.Sprintf("apps/%v/restores/%v/download", appId, restoreId)
	res, err := c.sendRequestRaw("GET", endpoint , nil , map[string]string{"Accept": "application/gzip"})

	if filename == "" {
		filename = parseContentDispositionFilename(res, "restore.tar.gz")
	}

	defer res.Body.Close()

	if err == nil {
		if isErrorCode(res.StatusCode) {
			var body []byte
			body, err = io.ReadAll(res.Body)
			if err == nil {
				err = formatRequestError(res.StatusCode, body)
			}
		}
	}
	AssertApiError(err, "appRestore")

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create file! %s", err.Error())
	}

	fmt.Printf("Saving report to %s\n", filename)

	defer file.Close()

	io.Copy(file, res.Body)
}
//-------------------------------------------------  APP COMPONENT BACKUP (GET) -------------------------------------------------
func (c *Client) AppComponentbackupsGet(appId int, componentId int) []types.AppComponentAvailableBackup {
	var backups struct {
		Data []types.AppComponentAvailableBackup `json:"availableBackups"`
	}
	endpoint := fmt.Sprintf("apps/%v/components/%v/availablebackups", appId, componentId)
	err := c.invokeAPI("GET", endpoint, nil, &backups)
	AssertApiError(err, "availablebackup")

	return backups.Data
}
