package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"bitbucket.org/level27/lvl/types"
)

func domainStatusCode(e error) {

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

}

//Domain gets a system from the API
func (c *Client) Domain(method string, id interface{}, data interface{}) types.Domain {
	var domain types.Domain

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

	return domain
}

//Domain gets a domain from the API
func (c *Client) Domains(filter string, number string) types.Domains {
	var domains types.Domains

	endpoint := "domains?limit=" + number + "&filter=" + filter
	err := c.invokeAPI("GET", endpoint, nil, &domains)
	AssertApiError(err)

	return domains
}

// ------------------ /DOMAINS --------------------------

// DESCRIBE DOMAIN (get detailed info from specific domain) - [lvl domain describe <id>]
func (c *Client) DomainDescribe(id []string) {
	if len(id) == 1 {
		domainID := id[0]
		domain := c.Domain("GET", domainID, nil).Data

		tmpl := template.Must(template.ParseFiles("templates/domain.tmpl"))
		err := tmpl.Execute(os.Stdout, domain)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("ERROR!")
	}
}

// DELETE DOMAIN
func (c *Client) DomainDelete(id []string) {
	if len(id) == 1 {
		domainID := id[0]
		// Ask for user confirmation to delete domain
		var userResponse string

		question := fmt.Sprintf("Are you sure you want to delete domain with ID: %v? Please type [y]es or [n]o: ", domainID)
		fmt.Print(question)
		_, err := fmt.Scan(&userResponse)
		if err != nil {
			log.Fatal(err)
		}

		switch strings.ToLower(userResponse) {
		case "y", "yes":
			c.Domain("DELETE", domainID, nil)
		case "n", "no":
			log.Fatal("Delete canceled")
		default:
			log.Println("Please make sure you type (y)es or (n)o and press enter to confirm:")
			domID := []string{domainID}
			c.DomainDelete(domID)
		}

	} else {
		fmt.Println("ERROR: wrong or invalid ID")

	}
}

// CREATE DOMAIN [lvl domain create <parmeters>]
func (c *Client) DomainCreate(args []string, req types.DomainRequest) {
	id := ""
	if req.Action == "" {
		req.Action = "none"
	}

	fmt.Println(req)
	c.Domain("CREATE", id, req)

}

// ------------------ /DOMAIN/RECORDS ----------------------
// GET
func (c *Client) DomainRecords(id string) []types.DomainRecord {
	var records struct {
		Records []types.DomainRecord `json:"records"`
	}

	endpoint := fmt.Sprintf("domains/%s/records", id)
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
