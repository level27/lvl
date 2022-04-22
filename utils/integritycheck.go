package utils

import (
	"fmt"
	"io"
	"log"
	"os"

	"bitbucket.org/level27/lvl/types"
)

// GET /{entityType}/{entityID}/integritychecks/{checkID}
func (c *Client) EntityIntegrityCheck(entityType string, entityID int, checkId int) types.IntegrityCheck {
	var result struct {
		IntegrityCheck types.IntegrityCheck `json:"integritycheck"`
	}

	endpoint := fmt.Sprintf("%s/%d/integritychecks/%d", entityType, entityID, checkId)
	err := c.invokeAPI("GET", endpoint, nil, &result)
	AssertApiError(err, "EntityIntegrityCheck")

	return result.IntegrityCheck
}

// GET /{entityType}/{entityID}/integritychecks
func (c *Client) EntityIntegrityChecks(entityType string, entityID int, getParams types.CommonGetParams) []types.IntegrityCheck {
	var result struct {
		IntegrityChecks []types.IntegrityCheck `json:"integritychecks"`
	}

	endpoint := fmt.Sprintf("%s/%d/integritychecks?%s", entityType, entityID, formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &result)
	AssertApiError(err, "EntityIntegrityChecks")

	return result.IntegrityChecks
}

// POST /{entityType}/{entityID}/integritychecks
func (c *Client) EntityIntegrityCreate(entityType string, entityID int, runJobs bool, forceRunJobs bool) types.IntegrityCheck {
	var result struct {
		IntegrityCheck types.IntegrityCheck `json:"integritycheck"`
	}

	endpoint := fmt.Sprintf("%s/%d/integritychecks", entityType, entityID)
	data := &types.IntegrityCreateRequest{Dojobs: runJobs, Forcejobs: forceRunJobs}
	err := c.invokeAPI("POST", endpoint, data, &result)
	AssertApiError(err, "EntityIntegrityCreate")

	return result.IntegrityCheck
}

// Download entity integrity check report to file.
func (c *Client) EntityIntegrityCheckDownload(entityType string, entityID int, checkId int, fileName string) {
	endpoint := fmt.Sprintf("%s/%d/integritychecks/%d/report", entityType, entityID, checkId)
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

	AssertApiError(err, "EntityIntegrityCheckDownload")

	if fileName == "" {
		fileName = parseContentDispositionFilename(res, fmt.Sprintf("integritycheck_%d_%s_%d.pdf", checkId, entityType, entityID))
	}

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to create file! %s", err.Error())
	}

	fmt.Printf("Saving report to %s\n", fileName)

	defer file.Close()

	io.Copy(file, res.Body)
}
