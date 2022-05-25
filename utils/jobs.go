package utils

import (
	"fmt"

	"bitbucket.org/level27/lvl/types"
)

func (c *Client) JobHistoryRootGet(rootJobId int) types.Job {
	var job types.Job
	endpoint := fmt.Sprintf("jobs/history/root/%v", rootJobId)
	err := c.invokeAPI("GET", endpoint, nil, &job)
	AssertApiError(err, "root job history")

	return job
}


func (c *Client) EntityJobHistoryGet(entityType string, domainId int) []types.Job {
	var historyResult []types.Job

	endpoint := fmt.Sprintf("jobs/history/%s/%v", entityType, domainId)
	err := c.invokeAPI("GET", endpoint, nil, &historyResult)

	AssertApiError(err, "job history")

	return historyResult
}
