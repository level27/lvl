package utils

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
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
func (c *Client) Domains(getParams types.CommonGetParams) []types.Domain {
	var domains struct {
		Data []types.Domain `json:"domains"`
	}

	endpoint := fmt.Sprintf("domains?%s", formatCommonGetParams(getParams))
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

	result := c.Domain("CREATE", nil, req)
	log.Printf("Domain created! [Fullname: '%v' , ID: '%v']", result.Fullname, result.ID)

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
func (c *Client) DomainRecords(id int, recordType string, getParams types.CommonGetParams) []types.DomainRecord {
	var records struct {
		Records []types.DomainRecord `json:"records"`
	}

	endpoint := fmt.Sprintf("domains/%d/records?%s", id, formatCommonGetParams(getParams))
	if recordType != "" {
		endpoint += fmt.Sprintf("&type=%s", recordType)
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
// func (c *Client) DomainNotificationGet(domainId int) []types.Notification {
// 	var notifications struct {
// 		Notifications []types.Notification `json:"notifications"`
// 	}
// 	endpoint := fmt.Sprintf("domains/%v/notifications", domainId)
// 	err := c.invokeAPI("GET", endpoint, nil, &notifications)
// 	AssertApiError(err, "notifications")
// 	return notifications.Notifications
// }

// // CREATE A NOTIFICATION
// func (c *Client) DomainNotificationAdd(domainId int, req types.DomainNotificationPostRequest) {
// 	endpoint := fmt.Sprintf("domains/%v/notifications", domainId)
// 	err := c.invokeAPI("POST", endpoint, req, nil)

// 	AssertApiError(err, "notifications")
// }

// --------------------------------------------------- BILLABLE ITEM --------------------------------------------------------

//--------------------------- CREATE (Turn invoicing on)
//CREATE BILLABLEITEM
func (c *Client) DomainBillableItemCreate(domainid int, req types.DomainBillPostRequest) {

	endpoint := fmt.Sprintf("domains/%v/bill", domainid)

	err := c.invokeAPI("POST", endpoint, req, nil)
	AssertApiError(err, "billable item")

}

// ---------------------------- DELETE (turn invoicing off)
//DELETE
func (c *Client) DomainBillableItemDelete(domainId int) {
	endpoint := fmt.Sprintf("domains/%v/billableitem", domainId)

	err := c.invokeAPI("DELETE", endpoint, nil, nil)
	AssertApiError(err, "Billable item")

}

// -------------------------------------------------------CHECK AVAILABILITY---------------------------------------------------------------------
// Check domain availability
func (c *Client) DomainCheck(name string, extension string) types.DomainCheckResult {
	var checkResult types.DomainCheckResult

	endpoint := fmt.Sprintf("domains/check?name=%s&extension=%s", url.QueryEscape(name), url.QueryEscape(extension))
	err := c.invokeAPI("GET", endpoint, nil, &checkResult)
	AssertApiError(err, "domainCheck")

	return checkResult
}

// ------------------------------------------------------- JOB HISTORY ---------------------------------------------------------------------
func (c *Client) DomainJobHistoryGet(domainId int) []types.Job {
	var historyResult []types.Job

	endpoint := fmt.Sprintf("jobs/history/domain/%v", domainId)
	err := c.invokeAPI("GET", endpoint, nil, &historyResult)

	AssertApiError(err, "job history")

	return historyResult
}

func (c *Client) DomainJobHistoryRootGet(rooJobId int) types.Job {
	var job types.Job
	endpoint := fmt.Sprintf("jobs/history/root/%v", rooJobId)
	err := c.invokeAPI("GET", endpoint, nil, &job)
	AssertApiError(err, "root job history")

	return job
}


// ---------------------------------------------- INTEGRITY CHECKS DOMAINS ------------------------------------------------

func (c *Client) DomainIntegrityCheck(domainId int, checkId int) types.DomainIntegrityCheck {
	var result struct {
		IntegrityCheck types.DomainIntegrityCheck `json:"integritycheck"`
	}

	endpoint := fmt.Sprintf("domains/%d/integritychecks/%d", domainId, checkId)
	err := c.invokeAPI("GET", endpoint, nil, &result)
	AssertApiError(err, "domainIntegrityCheck")

	return result.IntegrityCheck
}

func (c *Client) DomainIntegrityChecks(domainId int, getParams types.CommonGetParams) []types.IntegrityCheckSummary {
	var result struct {
		IntegrityChecks []types.IntegrityCheckSummary `json:"integritychecks"`
	}

	endpoint := fmt.Sprintf("domains/%d/integritychecks?%s", domainId, formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &result)
	AssertApiError(err, "domainIntegrityCheck")

	return result.IntegrityChecks
}

// Create domain integrity check
func (c *Client) DomainIntegrityCreate(domainId int, runJobs bool, forceRunJobs bool) types.DomainIntegrityCheck {
	var result struct {
		IntegrityCheck types.DomainIntegrityCheck `json:"integritycheck"`
	}

	endpoint := fmt.Sprintf("domains/%d/integritychecks", domainId)
	data := &types.IntegrityCreateRequest{Dojobs: runJobs, Forcejobs: forceRunJobs}
	err := c.invokeAPI("POST", endpoint, data, &result)
	AssertApiError(err, "domainIntegrityCheck")

	return result.IntegrityCheck
}

// Download domain integrity check report to file.
func (c *Client) DomainIntegrityCheckDownload(domainId int, checkId int, fileName string) {
	endpoint := fmt.Sprintf("domains/%d/integritychecks/%d/report", domainId, checkId)
	res, err := c.sendRequestRaw("GET", endpoint, nil, map[string]string{"Accept": "application/pdf"})

	if err == nil {
		defer res.Body.Close()

		if isErrorCode(res.StatusCode) {
			var body []byte
			body, err = io.ReadAll(res.Body)
			if err == nil {
				err = formatRequestError(res.StatusCode, body)
			}
		}
	}

	AssertApiError(err, "domainIntegrityCheckDownload")

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to create file! %s", err.Error())
	}

	fmt.Printf("Saving report to %s\n", fileName)

	defer file.Close()

	io.Copy(file, res.Body)
}
