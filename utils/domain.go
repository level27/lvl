package utils

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"bitbucket.org/level27/lvl/types"
)

func domainStatusCode(e error) {
	if e != nil {
		splittedError := strings.Split(e.Error(), " ")
		var result string
		switch splittedError[len(splittedError)-1] {
		case "204":
			result = "Request succesfully processed"
		case "400":
			result = "Bad request"
		case "403":
			result = "You do not have acces to this domain"
		case "404":
			result = "Domain not found"
		case "500":
			result = "You have no proper rights to acces the controller"
		default:
			result = "No Status code received"
		}

		log.Println(result)
	} else {
		log.Println("Request succesfully processed")
	}

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
		fmt.Println(data)
		err = c.invokeAPI("POST", endpoint, data, &domain)
	case "UPDATE":
		endpoint := fmt.Sprintf("domains/%s", id)
		err = c.invokeAPI("PUT", endpoint, data, &domain)
	case "DELETE":
		endpoint := fmt.Sprintf("domains/%s", id)

		err = c.invokeAPI("DELETE", endpoint, nil, nil)
	}

	domainStatusCode(err)
	AssertApiError(err)

	return domain.Data
}

//Domain gets a domain from the API
func (c *Client) Domains(filter string, number int) []types.Domain {
	var domains struct {
		Data []types.Domain `json:"domains"`
	}

	endpoint := fmt.Sprintf("domains?limit=%d&filter=%s", number, url.QueryEscape(filter))
	err := c.invokeAPI("GET", endpoint, nil, &domains)
	AssertApiError(err)

	return domains.Data
}

// ------------------ /DOMAINS --------------------------

// DELETE DOMAIN
func (c *Client) DomainDelete(id []string) {

	// looping over all given args and checking for valid domainId's
	for _, value := range id{

		domainId, err := strconv.Atoi(value)
		if err == nil  {
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
		}else{
			log.Printf("Wrong or invalid domain ID: %v.\n", value)
		}
	}

}

// CREATE DOMAIN [lvl domain create <parmeters>]
func (c *Client) DomainCreate(args []string, req types.DomainRequest) {

	if req.Action == "" {
		req.Action = "none"
	}
	if *req.DomainContactOnSite == 0 {
		req.DomainContactOnSite = nil
	}

	fmt.Println(req)

	c.Domain("CREATE", nil, req)

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
	AssertApiError(err)

	return records.Records
}

func (c *Client) DomainRecord(domainId int, recordId int) types.DomainRecord {
	var records struct {
		Record types.DomainRecord `json:"record"`
	}

	endpoint := fmt.Sprintf("domains/%d/records/%d", domainId, recordId)
	err := c.invokeAPI("GET", endpoint, nil, &records)
	AssertApiError(err)

	return records.Record
}

// CREATE
func (c *Client) DomainRecordCreate(id int, req types.DomainRecordRequest) types.DomainRecord {
	record := types.DomainRecord{}

	endpoint := fmt.Sprintf("domains/%d/records", id)
	err := c.invokeAPI("POST", endpoint, &req, &record)

	AssertApiError(err)

	return record
}

// DELETE
func (c *Client) DomainRecordDelete(domainId int, recordId int) {
	endpoint := fmt.Sprintf("domains/%d/records/%d", domainId, recordId)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)

	AssertApiError(err)
}

func (c *Client) DomainRecordUpdate(domainId int, recordId int, req types.DomainRecordRequest) {
	endpoint := fmt.Sprintf("domains/%d/records/%d", domainId, recordId)
	err := c.invokeAPI("PUT", endpoint, &req, nil)

	AssertApiError(err)
}
