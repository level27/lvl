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
