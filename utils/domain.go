package utils

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"bitbucket.org/level27/lvl/types"
)

//gets extensions for domains
func (c *Client) Extension() []types.DomainProvider {
	var extensions struct {
		Data []types.DomainProvider `json:"providers"`
	}

	endpoint := "domains/providers"
	err := c.invokeAPI("GET", endpoint, nil, &extensions)
	fmt.Println("dit zijn de providers ")
	fmt.Println(err)
	AssertApiError(err, "extension")

	return extensions.Data
}

//Domain gets a system from the API
func (c *Client) Domain(method string, id interface{}, data interface{}) types.Domain {
	var domain struct {
		Data types.Domain `json:"domain"`
	}

	var err error
	switch method {
	case "GET":
		endpoint := fmt.Sprintf("domains/%s", id)
		err = c.invokeAPI("GET", endpoint, nil, &domain)
	case "CREATE":
		endpoint := "domains"
		err = c.invokeAPI("POST", endpoint, data, &domain)
	case "UPDATE":
		endpoint := fmt.Sprintf("domains/%s", id)
		err = c.invokeAPI("PUT", endpoint, data, nil)
	case "DELETE":
		endpoint := fmt.Sprintf("domains/%s", id)
		err = c.invokeAPI("DELETE", endpoint, nil, nil)
	case "TRANSFER":
		endpoint := fmt.Sprintf("domains/%s/internaltransfer", id)
		err = c.invokeAPI("POST", endpoint, data, &domain)
	}

	AssertApiError(err, "domain")

	return domain.Data
}

//Domain gets a domain from the API
func (c *Client) Domains(filter string, number int) []types.Domain {
	var domains struct {
		Data []types.Domain `json:"domains"`
	}

	endpoint := fmt.Sprintf("domains?limit=%d&filter=%s", number, url.QueryEscape(filter))
	err := c.invokeAPI("GET", endpoint, nil, &domains)
	AssertApiError(err, "domains")

	return domains.Data
}

// ------------------ /DOMAINS --------------------------

// DELETE DOMAIN
func (c *Client) DomainDelete(id []string) {

	// looping over all given args and checking for valid domainId's
	for _, value := range id {

		domainId, err := strconv.Atoi(value)
		if err == nil {
			var userResponse string

			question := fmt.Sprintf("Are you sure you want to delete domain with ID: %v? Please type [y]es or [n]o: ", domainId)
			fmt.Print(question)
			_, err := fmt.Scan(&userResponse)
			if err != nil {
				log.Fatal(err)
			}

			switch strings.ToLower(userResponse) {
			case "y", "yes":
				c.Domain("DELETE", value, nil)
			case "n", "no":
				log.Printf("Delete canceled for domain: %v", value)
			default:
				log.Println("Please make sure you type (y)es or (n)o and press enter to confirm:")
				domID := []string{value}
				c.DomainDelete(domID)
			}
		} else {
			log.Printf("Wrong or invalid domain ID: %v.\n", value)
		}
	}

}

// CREATE DOMAIN [lvl domain create <parmeters>]
func (c *Client) DomainCreate(args []string, req types.DomainRequest) {

	if req.Action == "" {
		req.Action = "none"
	}

	test := c.Domain("CREATE", nil, req)
	fmt.Printf("handle dns: %v ", test.DNSIsHandled)
	log.Printf("domain created: '%v' - ID: '%v'", test.Fullname, test.ID)

}

// TRANSFER DOMAIN [lvl domain transfer <parameters>]
func (c *Client) DomainTransfer(args []string, req types.DomainRequest) {
	if req.Action == "" {
		req.Action = "transfer"
	}

	c.Domain("CREATE", nil, req)
}

// INTERNAL TRANSFER
func (c *Client) DomainInternalTransfer(args []string, req types.DomainRequest) {

	res := c.Domain("TRANSFER", args[0], req)

	fmt.Println(res)
}

// UPDATE DOMAIN [lvl update <parameters>]
func (c *Client) DomainUpdate(id int, data map[string]interface{}) {
	endpoint := fmt.Sprintf("domains/%d", id)
	err := c.invokeAPI("PATCH", endpoint, data, nil)
	AssertApiError(err, "domain update")
}

// ------------------ /DOMAIN/RECORDS ----------------------
// GET
func (c *Client) DomainRecords(id int, recordType string, limit int, filter string) []types.DomainRecord {
	var records struct {
		Records []types.DomainRecord `json:"records"`
	}

	endpoint := fmt.Sprintf("domains/%d/records?limit=%d", id, limit)
	if recordType != "" {
		endpoint += fmt.Sprintf("&type=%s", recordType)
	}
	if filter != "" {
		endpoint += fmt.Sprintf("&filter=%s", url.QueryEscape(filter))
	}
	err := c.invokeAPI("GET", endpoint, nil, &records)
	AssertApiError(err, "domain record")

	return records.Records
}

func (c *Client) DomainRecord(domainId int, recordId int) types.DomainRecord {
	var records struct {
		Record types.DomainRecord `json:"record"`
	}

	endpoint := fmt.Sprintf("domains/%d/records/%d", domainId, recordId)
	err := c.invokeAPI("GET", endpoint, nil, &records)
	AssertApiError(err, "domain record")

	return records.Record
}

// CREATE
func (c *Client) DomainRecordCreate(id int, req types.DomainRecordRequest) types.DomainRecord {
	record := types.DomainRecord{}

	endpoint := fmt.Sprintf("domains/%d/records", id)
	err := c.invokeAPI("POST", endpoint, &req, &record)

	AssertApiError(err, "domain record")

	return record
}

// DELETE
func (c *Client) DomainRecordDelete(domainId int, recordId int) {
	endpoint := fmt.Sprintf("domains/%d/records/%d", domainId, recordId)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	AssertApiError(err, "domain record")
}

// UPDATE
func (c *Client) DomainRecordUpdate(domainId int, recordId int, req types.DomainRecordRequest) {
	endpoint := fmt.Sprintf("domains/%d/records/%d", domainId, recordId)
	err := c.invokeAPI("PUT", endpoint, &req, nil)

	AssertApiError(err, "domain record")
}

// --------------------------------------------------- ACCESS --------------------------------------------------------
//add access to a domain

func (c *Client) DomainAccesAdd(domainId int, req types.DomainAccessRequest) {
	endpoint := fmt.Sprintf("domains/%v/acls", domainId)

	err := c.invokeAPI("POST", endpoint, &req, nil)

	AssertApiError(err, "Access")

}

//remove acces from a domain

func (c *Client) DomainAccesRemove(domainId int, organisationId int) {
	endpoint := fmt.Sprintf("domains/%v/acls/%v", domainId, organisationId)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	AssertApiError(err, "Access")
}

// --------------------------------------------------- NOTIFICATIONS --------------------------------------------------------
// GET LIST OF ALL NOTIFICATIONS FOR DOMAIN
func (c *Client) DomainNotificationGet(domainId int) []types.Notification {
	var notifications struct {
		Notifications []types.Notification `json:"notifications"`
	}
	endpoint := fmt.Sprintf("domains/%v/notifications", domainId)
	err := c.invokeAPI("GET", endpoint, nil, &notifications)
	AssertApiError(err, "notifications")
	return notifications.Notifications
}

// CREATE A NOTIFICATION
func (c *Client) DomainNotificationAdd(domainId int, req types.DomainNotificationPostRequest) {
	enpoint := fmt.Sprintf("domains/%v/notifications", domainId)
	err := c.invokeAPI("POST", enpoint, req, nil)

	AssertApiError(err, "notifications")
}

// --------------------------------------------------- BILLABLEITEM --------------------------------------------------------

func (c *Client) DomainBillableItemsGet(domainId int) types.BillableItemGet {
	var billableItem types.BillableItemGet
	endpoint := fmt.Sprintf("domains/%v/billableitem", domainId)
	err := c.invokeAPI("GET", endpoint, nil, &billableItem)
	AssertApiError(err, "BillableItem")

	return billableItem

}
