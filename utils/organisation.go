package utils

import (
	"fmt"

	"bitbucket.org/level27/lvl/types"
)

// Get a single organisation from the API.
func (c *Client) Organisation(organisationId int) types.Organisation {
	var orgs struct {
		Organisation types.Organisation `json:"organisation"`
	}

	endpoint := fmt.Sprintf("organisations/%d", organisationId)
	err := c.invokeAPI("GET", endpoint, nil, &orgs)
	AssertApiError(err, "organisation")

	return orgs.Organisation
}

//Organisation gets a organisation from the API
func (c *Client) Organisations(getParams types.CommonGetParams) []types.Organisation {
	var orgs struct {
		Organisation []types.Organisation `json:"organisations"`
	}

	endpoint := fmt.Sprintf("organisations?%s", formatCommonGetParams(getParams))
	err := c.invokeAPI("GET", endpoint, nil, &orgs)
	AssertApiError(err, "organisation")

	return orgs.Organisation
}

func (c *Client) LookupOrganisation(name string) []types.Organisation {
	results := []types.Organisation{}
	orgs := c.Organisations(types.CommonGetParams{ Filter: name })
	for _, org := range orgs {
		if org.Name == name {
			results = append(results, org)
		}
	}

	return results
}