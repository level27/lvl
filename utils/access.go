package utils

import (
	"fmt"

	"bitbucket.org/level27/lvl/types"
)

// GET /{entityType}/{entityID}/organisations
func (c *Client) EntityGetOrganisations(entityType string, entityID int) []types.OrganisationAccess {
	var response struct {
		Organisations []types.OrganisationAccess `json:"organisations"`
	}

	endpoint := fmt.Sprintf("%s/%d/organisations", entityType, entityID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "EntityGetOrganisations")

	return response.Organisations
}

// POST /{entityType}/{entityID}/acls
func (c *Client) EntityAddAcl(entityType string, entityID int, add types.AclAdd) types.Acl {
	var response struct {
		Acl types.Acl `json:"acl"`
	}

	endpoint := fmt.Sprintf("%s/%d/acls", entityType, entityID)
	err := c.invokeAPI("POST", endpoint, add, &response)
	AssertApiError(err, "EntityAddAcl")
	return response.Acl
}

// DELETE /{entityType}/{entityID}/acls/{organisationID}
func (c *Client) EntityRemoveAcl(entityType string, entityID int, organisationID int) {
	endpoint := fmt.Sprintf("%s/%d/acls/%d", entityType, entityID, organisationID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)
	AssertApiError(err, "EntityRemoveAcl")
}

