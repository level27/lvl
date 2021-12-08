package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"bitbucket.org/level27/lvl/types"
)

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
		err = c.invokeAPI("POST", endpoint, data, &domain)
	case "UPDATE":
		endpoint := fmt.Sprintf("domains/%s", id)
		err = c.invokeAPI("PUT", endpoint, data, &domain)
	case "DELETE":
		endpoint := fmt.Sprintf("domains/%s", id)
		err = c.invokeAPI("DELETE", endpoint, nil, nil)
	}

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

// CREATE DOMAIN [lvl domain create <id>]
func (c *Client) DomainCreate(name string) {

}
