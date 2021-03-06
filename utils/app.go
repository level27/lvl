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
func (c *Client) AppComponentLookup(appId int, name string) []types.AppComponent {
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

	log.Print("App succesfully updated!")
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


func (c *Client) AppComponentCreate(appId int, req interface{}) types.AppComponent {
	var app struct {
		Data types.AppComponent `json:"app"`
	}
	endpoint := fmt.Sprintf("apps/%d/components", appId)
	err := c.invokeAPI("POST", endpoint, req, &app)

	AssertApiError(err, "apps")

	return app.Data
}

func (c *Client) AppComponentUpdate(appId int, appComponentID int, req interface{}) {
	endpoint := fmt.Sprintf("apps/%d/components/%d", appId, appComponentID)
	err := c.invokeAPI("PUT", endpoint, req, nil)

	AssertApiError(err, "apps")
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

			c.AppComponentRestoresDelete(appId, restoreId, false)
		}
	}
}

// ---- DOWNLOAD RESTORE FILE
func (c *Client) AppComponentRestoreDownload(appId int, restoreId int, filename string) {
	endpoint := fmt.Sprintf("apps/%v/restores/%v/download", appId, restoreId)
	res, err := c.sendRequestRaw("GET", endpoint, nil, map[string]string{"Accept": "application/gzip"})

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
// ---- GET LIST OF COMPONENT AVAILABLEBACKUPS
func (c *Client) AppComponentbackupsGet(appId int, componentId int) []types.AppComponentAvailableBackup {
	var backups struct {
		Data []types.AppComponentAvailableBackup `json:"availableBackups"`
	}
	endpoint := fmt.Sprintf("apps/%v/components/%v/availablebackups", appId, componentId)
	err := c.invokeAPI("GET", endpoint, nil, &backups)
	AssertApiError(err, "availablebackup")

	return backups.Data
}

//-------------------------------------------------  APP MIGRATIONS (GET / DESCRIBE / CREATE / UPDATE) -------------------------------------------------
// ---- GET LIST OF MIGRATIONS
func (c *Client) AppMigrationsGet(appId int) []types.AppMigration {
	var migrations struct {
		Data []types.AppMigration `json:"migrations"`
	}

	endpoint := fmt.Sprintf("apps/%v/migrations", appId)
	err := c.invokeAPI("GET", endpoint, nil, &migrations)
	AssertApiError(err, "appMigration")

	return migrations.Data
}

// ---- CREATE APP MIGRATION
func (c *Client) AppMigrationsCreate(appId int, req types.AppMigrationRequest) {
	var migration struct {
		Data types.AppMigration `json:"migration"`
	}
	endpoint := fmt.Sprintf("apps/%v/migrations", appId)
	err := c.invokeAPI("POST", endpoint, req, &migration)
	AssertApiError(err, "appMigration")

	log.Printf("migration created! [ID: '%v']", migration.Data.ID)
}

// ---- UPDATE APP MIGRATION
func (c *Client) AppMigrationsUpdate(appId int, migrationId int, req interface{}) {
	endpoint := fmt.Sprintf("apps/%v/migrations/%v", appId, migrationId)
	err := c.invokeAPI("PUT", endpoint, req, nil)
	AssertApiError(err, "appMigration")

	log.Print("migration succesfully updated!")
}

// ---- DESCRIBE APP MIGRATION
func (c *Client) AppMigrationDescribe(appId int, migrationId int) types.AppMigration {
	var migration struct {
		Data types.AppMigration `json:"migration"`
	}

	endpoint := fmt.Sprintf("apps/%v/migrations/%v", appId, migrationId)
	err := c.invokeAPI("GET", endpoint, nil, &migration)
	AssertApiError(err, "appMigration")

	return migration.Data
}

//-------------------------------------------------  APP MIGRATIONS ACTIONS (CONFIRM / DENY / RESTART) -------------------------------------------------
// ---- MIGRATIONS ACTION COMMAND
func (c *Client) AppMigrationsAction(appId int, migrationId int, ChosenAction string) {
	var action struct {
		Type string `json:"type"`
	}

	action.Type = ChosenAction
	endpoint := fmt.Sprintf("apps/%v/migrations/%v/actions", appId, migrationId)
	err := c.invokeAPI("POST", endpoint, action, nil)

	AssertApiError(err, "appMigrationAction")
}


// ------------ COMPONENT URL MANAGEMENT

// GET /apps/{appId}/components/{componentId}/urls
func (c *Client) AppComponentUrlGetList(appID int, componentID int, get types.CommonGetParams) []types.AppComponentUrlShort {
	var resp struct {
		Urls []types.AppComponentUrlShort `json:"urls"`
	}

	endpoint := fmt.Sprintf("apps/%d/components/%d/urls?%s", appID, componentID, formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &resp)
	AssertApiError(err, "AppComponentUrlGetList")

	return resp.Urls
}

// GET /apps/{appId}/components/{componentId}/urls/{urlId}
func (c *Client) AppComponentUrlGetSingle(appID int, componentID int, urlID int) types.AppComponentUrl {
	var resp struct {
		Url types.AppComponentUrl `json:"url"`
	}

	endpoint := fmt.Sprintf("apps/%d/components/%d/urls/%d", appID, componentID, urlID)
	err := c.invokeAPI("GET", endpoint, nil, &resp)
	AssertApiError(err, "AppComponentUrlGetSingle")

	return resp.Url
}

// POST /apps/{appId}/components/{componentId}/urls
func (c *Client) AppComponentUrlCreate(appID int, componentID int, create types.AppComponentUrlCreate) types.AppComponentUrl {
	var resp struct {
		Url types.AppComponentUrl `json:"url"`
	}

	endpoint := fmt.Sprintf("apps/%d/components/%d/urls", appID, componentID)
	err := c.invokeAPI("POST", endpoint, create, &resp)
	AssertApiError(err, "AppComponentUrlCreate")

	return resp.Url
}

// PUT /apps/{appId}/components/{componentId}/urls/{urlId}
func (c *Client) AppComponentUrlUpdate(appID int, componentID int, urlID int, data interface{}) {
	endpoint := fmt.Sprintf("apps/%d/components/%d/urls/%d", appID, componentID, urlID)
	err := c.invokeAPI("PUT", endpoint, data, nil)
	AssertApiError(err, "AppComponentUrlUpdate")
}

// DELETE /apps/{appId}/components/{componentId}/urls/{urlId}
func (c *Client) AppComponentUrlDelete(appID int, componentID int, urlID int) {
	endpoint := fmt.Sprintf("apps/%d/components/%d/urls/%d", appID, componentID, urlID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)
	AssertApiError(err, "AppComponentUrlDelete")
}

func (c *Client) AppComponentUrlLookup(appID int, componentID int, name string) []types.AppComponentUrlShort {
	results := []types.AppComponentUrlShort{}
	urls := c.AppComponentUrlGetList(appID, componentID, types.CommonGetParams{Filter: name})
	for _, url := range urls {
		if url.Content == name {
			results = append(results, url)
		}
	}

	return results
}