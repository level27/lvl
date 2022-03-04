package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"bitbucket.org/level27/lvl/types"
	"github.com/Jeffail/gabs/v2"
)

// --------------------------- TOPLEVEL SYSTEM ACTIONS (GET / POST) ------------------------------------
//------------------ GET

// returning a list of all current systems [lvl system get]
func (c *Client) SystemGetList(getParams types.CommonGetParams) []types.System {

	//creating an array of systems.
	var systems struct {
		Data []types.System `json:"systems"`
	}

	//creating endpoint
	endpoint := fmt.Sprintf("systems?%s", formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &systems)
	AssertApiError(err, "Systems")
	//returning result as system type
	return systems.Data

}

func (c *Client) LookupSystem(name string) *types.System {
	systems := c.SystemGetList(types.CommonGetParams{Filter: name})
	for _, system := range systems {
		if system.Name == name {
			return &system
		}
	}

	return nil
}

// Returning a single system by its ID
// this is not for a describe.
func (c *Client) SystemGetSingle(id int) types.System {
	var system struct {
		Data types.System `json:"system"`
	}
	endpoint := fmt.Sprintf("systems/%v", id)
	err := c.invokeAPI("GET", endpoint, nil, &system)

	AssertApiError(err, "System")
	return system.Data

}

func (c *Client) SystemGetSshKeys(id int, get types.CommonGetParams) []types.SshKey {
	var keys struct {
		SshKeys []types.SshKey `json:"sshKeys"`
	}

	endpoint := fmt.Sprintf("systems/%d/sshkeys?%s", id, formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &keys)

	AssertApiError(err, "System SSH Keys")
	return keys.SshKeys
}

func (c *Client) SystemGetHasNetworks(id int) []types.SystemHasNetwork {
	var keys struct {
		SystemHasNetworks []types.SystemHasNetwork `json:"systemHasNetworks"`
	}

	endpoint := fmt.Sprintf("systems/%d/networks", id)
	err := c.invokeAPI("GET", endpoint, nil, &keys)

	AssertApiError(err, "System has networks")
	return keys.SystemHasNetworks
}

func (c *Client) SystemGetVolumes(id int, get types.CommonGetParams) []types.SystemVolume {
	var keys struct {
		Volumes []types.SystemVolume `json:"volumes"`
	}

	endpoint := fmt.Sprintf("systems/%d/volumes?%s", id, formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &keys)

	AssertApiError(err, "Volumes")
	return keys.Volumes
}

func (c *Client) SecurityUpdateDates() []string {
	var updates struct {
		SecurityUpdateDates []string `json:"securityUpdateDates"`
	}

	endpoint := "systems/securityupdatedates"
	err := c.invokeAPI("GET", endpoint, nil, &updates)

	AssertApiError(err, "Security updates")
	return updates.SecurityUpdateDates
}

//----------------- POST
//Get request to see all curent checktypes (valid checktype needed to create new check)
func (c *Client) SystemCheckTypeGet() []string {
	var checks struct {
		Data types.SystemCheckTypeName `json:"checktypes"`
	}

	endpoint := "checktypes"
	err := c.invokeAPI("GET", endpoint, nil, &checks)
	AssertApiError(err, "checktypes")

	//creating an array from the maps keys. the keys of the map are the possible checktypes
	validTypes := make([]string, 0, len(checks.Data))
	values := make([]types.SystemCheckType, 0, len(checks.Data))

	for K, V := range checks.Data {
		validTypes = append(validTypes, K)
		values = append(values, V)
	}

	return validTypes

}

// CREATE SYSTEM [lvl system create <parmeters>]
func (c *Client) SystemCreate(req types.SystemPost) {

	var System struct {
		Data types.System `json:"system"`
	}

	err := c.invokeAPI("POST", "systems", req, &System)
	AssertApiError(err, "SystemCreate")

	log.Printf("System created! [Fullname: '%v' , ID: '%v']", System.Data.Name, System.Data.Id)

}

// SYSTEM ACTION

func (c *Client) SystemAction(id int, action string) types.System {
	var request struct {
		Type string `json:"type"`
	}

	var response struct {
		System types.System `json:"system"`
	}

	request.Type = action
	endpoint := fmt.Sprintf("systems/%d/actions", id)
	err := c.invokeAPI("POST", endpoint, request, &response)
	AssertApiError(err, "SystemAction")

	return response.System
}

// --------------------------- SYSTEM/CHECKS TOPLEVEL (GET / POST) ------------------------------------
// ------------- GET CHECKS
func (c *Client) SystemCheckGetList(systemId int, getParams types.CommonGetParams) []types.SystemCheck {

	//creating an array of systems.
	var systemChecks struct {
		Data []types.SystemCheck `json:"checks"`
	}

	//creating endpoint
	endpoint := fmt.Sprintf("systems/%v/checks?%s", systemId, formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &systemChecks)
	AssertApiError(err, "Systems")
	//returning result as system check type
	return systemChecks.Data

}

// ------------- CREATE A CHECK
func (c *Client) SystemCheckCreate(systemId int, req interface{}) {
	var SystemCheck struct {
		Data types.SystemCheck `json:"check"`
	}
	endpoint := fmt.Sprintf("systems/%v/checks", systemId)
	err := c.invokeAPI("POST", endpoint, req, &SystemCheck)

	AssertApiError(err, "System checks")
	log.Printf("System check created! [Checktype: '%v' , ID: '%v']", SystemCheck.Data.CheckType, SystemCheck.Data.Id)
}

// --------------------------- SYSTEM/CHECKS ACTIONS (GET / DELETE / UPDATE) ------------------------------------
// ------------- DESCRIBE A SPECIFIC CHECK
func (c *Client) SystemCheckDescribe(systemID int, CheckID int) types.SystemCheck {
	var check struct {
		Data types.SystemCheck `json:"check"`
	}
	endpoint := fmt.Sprintf("systems/%v/checks/%v", systemID, CheckID)
	err := c.invokeAPI("GET", endpoint, nil, &check)
	AssertApiError(err, "system check")

	return check.Data
}

// ------------- DELETE A SPECIFIC CHECK
func (c *Client) SystemCheckDelete(systemId int, checkId int, isDeleteConfirmed bool) {

	// when confirmation flag is set, delete check without confirmation question
	if isDeleteConfirmed {
		endpoint := fmt.Sprintf("systems/%v/checks/%v", systemId, checkId)
		err := c.invokeAPI("DELETE", endpoint, nil, nil)
		AssertApiError(err, "system check")
	} else {
		var userResponse string
		// ask user for confirmation on deleting the check
		question := fmt.Sprintf("Are you sure you want to delete the systems check with ID: %v? Please type [y]es or [n]o: ", checkId)
		fmt.Print(question)
		//reading user response
		_, err := fmt.Scan(&userResponse)
		if err != nil {
			log.Fatal(err)
		}
		// check if user confirmed the deletion of the check or not
		switch strings.ToLower(userResponse) {
		case "y", "yes":
			endpoint := fmt.Sprintf("systems/%v/checks/%v", systemId, checkId)
			err := c.invokeAPI("DELETE", endpoint, nil, nil)
			AssertApiError(err, "system check")
		case "n", "no":
			log.Printf("Delete canceled for system check: %v", checkId)
		default:
			log.Println("Please make sure you type (y)es or (n)o and press enter to confirm:")

			c.SystemCheckDelete(systemId, checkId, false)
		}

	}

}

// ------------- UPDATE A SPECIFIC CHECK
func (c *Client) SystemCheckUpdate(systemId int, checkId int, req types.SystemCheckRequestHttp) {
	var SystemCheck struct {
		Data types.SystemCheck `json:"check"`
	}
	endpoint := fmt.Sprintf("systems/%v/checks", systemId)
	err := c.invokeAPI("PUT", endpoint, req, &SystemCheck)

	AssertApiError(err, "System checks")
}

// --------------------------- SYSTEM/COOKBOOKS TOPLEVEL (GET / POST) ------------------------------------
// ------------- GET COOKBOOK

func (c *Client) SystemCookbookGetList(systemId int) []types.Cookbook {
	// creating array of cookbooks to return
	var systemCookbooks struct {
		Data []types.Cookbook `json:"cookbooks"`
	}

	endpoint := fmt.Sprintf("systems/%v/cookbooks", systemId)
	err := c.invokeAPI("GET", endpoint, nil, &systemCookbooks)

	AssertApiError(err, "cookbooks")

	return systemCookbooks.Data

}

// ------------- ADD COOKBOOK
func (c *Client) SystemCookbookAdd(systemID int, req interface{}) {

	// var to show result of API after succesfull adding cookbook
	var cookbook struct{
		Data types.Cookbook `json:"cookbook"`
	}

	endpoint := fmt.Sprintf("systems/%v/cookbooks", systemID)
	err := c.invokeAPI("POST", endpoint, req, &cookbook)
	AssertApiError(err, "cookbooktype")


}

// --------------------------- SPECIFIC COOKBOOKTYPES (GET) ------------------------------------
// ------- GET ALL CURRENT COOKBOOKTYPES
func (c *Client) SystemCookbookTypesGet() ([]string, *gabs.Container) {
	var cookbookTypes struct {
		Data types.CookbookTypeName `json:"cookbooktypes"`
	}
	endpoint := "cookbooktypes"
	err := c.invokeAPI("GET", endpoint, nil, &cookbookTypes)
	AssertApiError(err, "cookbooktypes")

	//creating an array from the maps keys. the keys of the map are the possible cookbooktypes
	validTypes := make([]string, 0, len(cookbookTypes.Data))

	for K, _ := range cookbookTypes.Data {
		validTypes = append(validTypes, K)

	}

	// marshal all current data into a byteslice to parse it into json.
	result, err := json.Marshal(cookbookTypes)
	if err != nil {
		log.Fatal(err.Error())
	}

	// parse the slice of bytes into json, this way we can dynamicaly user unknown incomming data
	jsonParsed, err := gabs.ParseJSON([]byte(result))

	if err != nil {
		log.Fatal(err.Error())
	}
	
	return validTypes, jsonParsed
}

// ------------------ GET PROVIDERS

func (c *Client) GetSystemProviderConfigurations() []types.SystemProviderConfiguration {
	var response struct {
		ProviderConfigurations []types.SystemProviderConfiguration `json:"providerConfigurations"`
	}

	err := c.invokeAPI("GET", "systems/provider/configurations", nil, &response)
	AssertApiError(err, "GetSystemProviderConfigurations")

	return response.ProviderConfigurations
}
