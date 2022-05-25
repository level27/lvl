package utils

import (
	"fmt"

	"bitbucket.org/level27/lvl/types"
)

// POST /{entityType}/{systemID}/bill
func (c *Client) EntityBillableItemCreate(entityType string, entityID int, req types.BillPostRequest) {

	endpoint := fmt.Sprintf("%s/%v/bill", entityType, entityID)

	err := c.invokeAPI("POST", endpoint, req, nil)
	AssertApiError(err, "EntityBillableItemCreate")
}

// DELETE /{entityType}/{systemID}/billableitem
func (c *Client) EntityBillableItemDelete(entityType string, entityID int) {
	endpoint := fmt.Sprintf("%s/%v/billableitem", entityType, entityID)

	err := c.invokeAPI("DELETE", endpoint, nil, nil)
	AssertApiError(err, "EntityBillableItemDelete")
}