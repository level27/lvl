package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"bitbucket.org/level27/lvl/types"
	"github.com/Jeffail/gabs/v2"
)

// --------------------------- TOPLEVEL SYSTEM ACTIONS (GET / POST) ------------------------------------
// #region SYSTEM TOPLEVEL (GET / CREATE)
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

// CREATE SYSTEM [lvl system create <parmeters>]
func (c *Client) SystemCreate(req types.SystemPost) {

	var System struct {
		Data types.System `json:"system"`
	}

	err := c.invokeAPI("POST", "systems", req, &System)
	AssertApiError(err, "SystemCreate")

	log.Printf("System created! [Fullname: '%v' , ID: '%v']", System.Data.Name, System.Data.Id)

}

// #endregion

// --------------------------- @PJ please fill in comments about code ------------------------------------
// #region  @PJ please fill in comments about code
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

func (c *Client) SystemGetSshKeys(id int, get types.CommonGetParams) []types.SystemSshkey {
	var keys struct {
		SshKeys []types.SystemSshkey `json:"sshkeys"`
	}

	endpoint := fmt.Sprintf("systems/%d/sshkeys?%s", id, formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &keys)

	AssertApiError(err, "System SSH Keys")
	return keys.SshKeys
}

func (c *Client) SystemGetNonAddedSshKeys(systemID int, organisationID int, userID int, get types.CommonGetParams) []types.SshKey {
	var keys struct {
		SshKeys []types.SshKey `json:"sshKeys"`
	}

	endpoint := fmt.Sprintf("systems/%d/organisations/%d/users/%d/nonadded-sshkeys?%s", systemID, organisationID, userID, formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &keys)

	AssertApiError(err, "system nonadded SSH Keys")
	return keys.SshKeys
}

func (c *Client) SystemAddSshKey(id int, keyID int) types.SshKey {
	var key struct {
		Sshkey types.SshKey `json:"sshKey"`
	}

	var data struct {
		Sshkey int `json:"sshkey"`
	}

	data.Sshkey = keyID

	endpoint := fmt.Sprintf("systems/%d/sshkeys", id)
	err := c.invokeAPI("POST", endpoint, &data, &key)

	AssertApiError(err, "Add SSH key")
	return key.Sshkey
}

func (c *Client) SystemRemoveSshKey(id int, keyID int) {

	endpoint := fmt.Sprintf("systems/%d/sshkeys/%d", id, keyID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	AssertApiError(err, "Add SSH key")
}

func (c *Client) LookupSystemSshkey(systemID int, name string) *types.SystemSshkey {
	keys := c.SystemGetSshKeys(systemID, types.CommonGetParams{Filter: name})
	for _, key := range keys {
		if key.Description == name {
			return &key
		}
	}

	return nil
}

func (c *Client) LookupSystemNonAddedSshkey(systemID int, organisationID int, userID int, name string) *types.SshKey {
	keys := c.SystemGetNonAddedSshKeys(systemID, organisationID, userID, types.CommonGetParams{Filter: name})
	for _, key := range keys {
		if key.Description == name {
			return &key
		}
	}

	return nil
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

func (c *Client) SystemUpdate(id int, data map[string]interface{}) {
	endpoint := fmt.Sprintf("systems/%d", id)
	err := c.invokeAPI("PUT", endpoint, data, nil)
	AssertApiError(err, "SystemUpdate")
}

// --------------------------- SYSTEM ACTION ---------------------------

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

// ---------------- Delete
func (c *Client) SystemDelete(id int) {
	endpoint := fmt.Sprintf("systems/%v", id)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	AssertApiError(err, "SystemDelete")
}

func (c *Client) SystemDeleteForce(id int) {
	endpoint := fmt.Sprintf("systems/%v/force", id)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	AssertApiError(err, "SystemDelete")
}

// #endregion

// --------------------------- SYSTEM/CHECKS TOPLEVEL (GET / POST ) ------------------------------------
// #region SYSTEM/CHECKS TOPLEVEL (GET / ADD)
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

// ------------- ADD A CHECK
func (c *Client) SystemCheckCreate(systemId int, req interface{}) {
	var SystemCheck struct {
		Data types.SystemCheck `json:"check"`
	}
	endpoint := fmt.Sprintf("systems/%v/checks", systemId)
	err := c.invokeAPI("POST", endpoint, req, &SystemCheck)

	AssertApiError(err, "System checks")
	log.Printf("System check created! [Checktype: '%v' , ID: '%v']", SystemCheck.Data.CheckType, SystemCheck.Data.Id)
}

// #endregion

// --------------------------- SYSTEM/CHECKS PARAMETERS (GET) ------------------------------------
// #region SYSTEM/CHECKS PARAMETERS (GET)

// ---------------- GET CHECK PARAMETERS (for specific checktype)
func (c *Client) SystemCheckTypeGet(checktype string) types.SystemCheckType {
	var checktypes struct {
		Data types.SystemCheckTypeName `json:"checktypes"`
	}
	endpoint := "checktypes"
	err := c.invokeAPI("GET", endpoint, nil, &checktypes)
	AssertApiError(err, "checktypes")

	// check if the given type by user is one of the possible types we got back from the API
	var isTypeValid = false
	for validType := range checktypes.Data {
		if checktype == validType {
			isTypeValid = true
			log.Print()
		}
	}

	// when given type is not valid -> error
	if !isTypeValid {
		message := fmt.Sprintf("given type: '%v' is no valid checktype.", checktype)
		err := errors.New(message)
		log.Fatal(err)
	}

	// return the chosen valid type and its specific data
	return checktypes.Data[checktype]
}

// #endregion

// --------------------------- SYSTEM/CHECKS SPECIFIC ACTIONS (DESCRIBE / DELETE / UPDATE) ------------------------------------
// #region SYSTEM/CHECKS SPECIFIC (DESCRIBE / DELETE / UPDATE)
// ---------------- DESCRIBE A SPECIFIC CHECK
func (c *Client) SystemCheckDescribe(systemID int, CheckID int) types.SystemCheck {
	var check struct {
		Data types.SystemCheck `json:"check"`
	}
	endpoint := fmt.Sprintf("systems/%v/checks/%v", systemID, CheckID)
	err := c.invokeAPI("GET", endpoint, nil, &check)
	AssertApiError(err, "system check")

	return check.Data
}

// ---------------- DELETE A SPECIFIC CHECK
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

// ---------------- UPDATE A SPECIFIC CHECK
func (c *Client) SystemCheckUpdate(systemId int, checkId int, req interface{}) {

	endpoint := fmt.Sprintf("systems/%v/checks/%v", systemId, checkId)
	err := c.invokeAPI("PUT", endpoint, req, nil)

	AssertApiError(err, "System checks")
}

// #endregion

// --------------------------- SYSTEM/COOKBOOKS TOPLEVEL (GET / POST) ------------------------------------
// #region SYSTEM/COOKBOOKS TOPLEVEL (GET / ADD)

// ---------------- GET COOKBOOK
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

// ---------------- ADD COOKBOOK
func (c *Client) SystemCookbookAdd(systemID int, req interface{}) {

	// var to show result of API after succesfull adding cookbook
	var cookbook struct {
		Data types.Cookbook `json:"cookbook"`
	}

	endpoint := fmt.Sprintf("systems/%v/cookbooks", systemID)
	err := c.invokeAPI("POST", endpoint, req, &cookbook)
	AssertApiError(err, "cookbooktype")

}

// #endregion

// --------------------------- SYSTEM/COOKBOOKS PARAMETERS (GET) ------------------------------------
// #region SYSTEM/COOKBOOKS PARAMETERS (GET)
// ---------------- GET COOKBOOKTYPES parameters
func (c *Client) SystemCookbookTypeGet(cookbooktype string) (types.CookbookType, *gabs.Container) {
	var cookbookTypes struct {
		Data types.CookbookTypeName `json:"cookbooktypes"`
	}
	endpoint := "cookbooktypes"
	err := c.invokeAPI("GET", endpoint, nil, &cookbookTypes)
	AssertApiError(err, "cookbooktypes")

	// check if the given type by user is one of the possible types we got back from the API
	var isTypeValid = false
	for validType := range cookbookTypes.Data {
		if cookbooktype == validType {
			isTypeValid = true

		}
	}

	// when given type is not valid -> error
	if !isTypeValid {
		message := fmt.Sprintf("given type: '%v' is no valid cookbooktype.", cookbooktype)
		err := errors.New(message)
		log.Fatal(err)
	}

	// from the valid type we make a JSON string with selectable parameters.
	// we do this because we dont know beforehand if there will be any and how they will be named
	result, err := json.Marshal(cookbookTypes.Data[cookbooktype].CookbookType.ParameterOptions)
	if err != nil {
		log.Fatal(err.Error())
	}
	// parse the slice of bytes into json, this way we can dynamicaly use unknown incomming data
	jsonParsed, err := gabs.ParseJSON([]byte(result))

	if err != nil {
		log.Fatal(err.Error())
	}

	// return the chosen valid type and its specific data
	return cookbookTypes.Data[cookbooktype], jsonParsed
}

// #endregion

// --------------------------- SYSTEM/COOKBOOKS SPECIFIC (DESCRIBE / DELETE / UPDATE) ------------------------------------
// #region SYSTEM/COOKBOOKS SPECIFIC (DESCRIBE / DELETE / UPDATE)

// ---------------- DESCRIBE
func (c *Client) SystemCookbookDescribe(systemId int, cookbookId int) types.Cookbook {
	var cookbook struct {
		Data types.Cookbook `json:"cookbook"`
	}

	endpoint := fmt.Sprintf("systems/%v/cookbooks/%v", systemId, cookbookId)
	err := c.invokeAPI("GET", endpoint, nil, &cookbook)
	AssertApiError(err, "system check")

	return cookbook.Data
}

// ---------------- DELETE
func (c *Client) SystemCookbookDelete(systemId int, cookbookId int, isDeleteConfirmed bool) {

	// when confirmation flag is set, delete check without confirmation question
	if isDeleteConfirmed {
		endpoint := fmt.Sprintf("systems/%v/cookbooks/%v", systemId, cookbookId)
		err := c.invokeAPI("DELETE", endpoint, nil, nil)
		AssertApiError(err, "system cookbook")
	} else {
		var userResponse string
		// ask user for confirmation on deleting the check
		question := fmt.Sprintf("Are you sure you want to delete the systems cookbook with ID: %v? Please type [y]es or [n]o: ", cookbookId)
		fmt.Print(question)
		//reading user response
		_, err := fmt.Scan(&userResponse)
		if err != nil {
			log.Fatal(err)
		}
		// check if user confirmed the deletion of the check or not
		switch strings.ToLower(userResponse) {
		case "y", "yes":
			endpoint := fmt.Sprintf("systems/%v/cookbooks/%v", systemId, cookbookId)
			err := c.invokeAPI("DELETE", endpoint, nil, nil)
			AssertApiError(err, "system cookbook")
		case "n", "no":
			log.Printf("Delete canceled for system check: %v", cookbookId)
		default:
			log.Println("Please make sure you type (y)es or (n)o and press enter to confirm:")

			c.SystemCookbookDelete(systemId, cookbookId, false)
		}

	}

}

// ------------------ UPDATE
func (c *Client) SystemCookbookUpdate(systemId int, cookbookId int, req interface{}) {

	endpoint := fmt.Sprintf("systems/%v/cookbooks/%v", systemId, cookbookId)
	err := c.invokeAPI("PUT", endpoint, req, nil)
	AssertApiError(err, "system/cookbook")

}

// #endregion

// --------------------------- SYSTEM/INTEGRITYCHECKS TOPLEVEL (GET / CREATE) ------------------------------------

// ------------------ GET
func (c *Client) SystemIntegritychecksGet(systemID int) types.IntegrityCheck {

	var integrity struct{
		Data types.IntegrityCheck `json:"integritychecks"`
	}
	
	endpoint := fmt.Sprintf("systems/%v/integritychecks", systemID)
	err := c.invokeAPI("GET", endpoint, nil, integrity.Data)
	AssertApiError(err, "system/integritycheck")

	return integrity.Data
}

// --------------------------- APPLY COOKBOOKCHANGES ON A SYSTEM
func (c *Client) SystemCookbookChangesApply(systemId int) {
	// create json format for post request
	// this function is specifically for updating cookbook status on a system
	requestData := gabs.New()
	requestData.Set("update_cookbooks", "type")

	endpoint := fmt.Sprintf("systems/%v/actions", systemId)
	err := c.invokeAPI("POST", endpoint, requestData, nil)
	AssertApiError(err, "systems/cookbook")

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

// NETWORKS

func (c *Client) LookupSystemHasNetworks(systemID int, name string) *types.SystemHasNetwork {
	networks := c.SystemGetHasNetworks(systemID)
	for _, network := range networks {
		if network.Network.Name == name {
			return &network
		}
	}

	return nil
}

func (c *Client) GetSystemHasNetwork(systemID int, systemHasNetworkID int) types.SystemHasNetwork {
	var response struct {
		SystemHasNetwork types.SystemHasNetwork `json:"systemHasNetwork"`
	}

	endpoint := fmt.Sprintf("systems/%d/networks/%d", systemID, systemHasNetworkID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "GetSystemHasNetwork")

	return response.SystemHasNetwork
}

func (c *Client) SystemAddHasNetwork(systemID int, networkID int) types.SystemHasNetwork {
	var response struct {
		SystemHasNetwork types.SystemHasNetwork `json:"systemHasNetwork"`
	}

	var request struct {
		Network int `json:"network"`
	}

	request.Network = networkID

	endpoint := fmt.Sprintf("systems/%d/networks", systemID)
	err := c.invokeAPI("POST", endpoint, &request, &response)
	AssertApiError(err, "SystemAddHasNetwork")

	return response.SystemHasNetwork
}

func (c *Client) SystemRemoveHasNetwork(systemID int, hasNetworkID int) {
	endpoint := fmt.Sprintf("systems/%d/networks/%d", systemID, hasNetworkID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)
	AssertApiError(err, "SystemRemoveHasNetwork")
}

func (c *Client) SystemGetHasNetworkIp(systemID int, hasNetworkID int, systemHasNetworkIpID int) types.SystemHasNetworkIp {
	var response struct {
		SystemHasNetworkIp types.SystemHasNetworkIp `json:"systemHasNetworkIp"`
	}

	endpoint := fmt.Sprintf("systems/%d/networks/%d/ips/%d", systemID, hasNetworkID, systemHasNetworkIpID)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	AssertApiError(err, "SystemGetHasNetworkIp")
	return response.SystemHasNetworkIp
}

func (c *Client) SystemGetHasNetworkIps(systemID int, hasNetworkID int) []types.SystemHasNetworkIp {
	var response struct {
		SystemHasNetworkIps []types.SystemHasNetworkIp `json:"systemHasNetworkIps"`
	}

	endpoint := fmt.Sprintf("systems/%d/networks/%d/ips", systemID, hasNetworkID)
	err := c.invokeAPI("GET", endpoint, nil, &response)

	AssertApiError(err, "SystemGetHasNetworkIps")
	return response.SystemHasNetworkIps
}

func (c *Client) SystemAddHasNetworkIps(systemID int, hasNetworkID int, add types.SystemHasNetworkIpAdd) types.SystemHasNetworkIp {
	var response struct {
		HasNetwork types.SystemHasNetworkIp `json:"systemHasNetworkIp"`
	}

	endpoint := fmt.Sprintf("systems/%d/networks/%d/ips", systemID, hasNetworkID)
	err := c.invokeAPI("POST", endpoint, add, &response)

	AssertApiError(err, "SystemAddHasNetworkIps")
	return response.HasNetwork
}

func (c *Client) SystemRemoveHasNetworkIps(systemID int, hasNetworkID int, ipID int) {
	endpoint := fmt.Sprintf("systems/%d/networks/%d/ips/%d", systemID, hasNetworkID, ipID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	AssertApiError(err, "SystemRemoveHasNetworkIps")
}

func (c *Client) LookupSystemHasNetworkIp(systemID int, hasNetworkID int, address string) *types.SystemHasNetworkIp {
	ips := c.SystemGetHasNetworkIps(systemID, hasNetworkID)
	for _, ip := range ips {
		if IpsEqual(Ipv4StringIntToString(ip.Ipv4), address) || IpsEqual(ip.Ipv6, address) || IpsEqual(Ipv4StringIntToString(ip.PublicIpv4), address) || IpsEqual(ip.PublicIpv6, address) {
			return &ip
		}
	}

	return nil
}

func (c *Client) SystemHasNetworkIpUpdate(systemID int, hasNetworkID int, hasNetworkIpID int, data map[string]interface{}) {
	endpoint := fmt.Sprintf("systems/%d/networks/%d/ips/%d", systemID, hasNetworkID, hasNetworkIpID)
	err := c.invokeAPI("PUT", endpoint, data, nil)
	AssertApiError(err, "SystemHasNetworkIpUpdate")
}
