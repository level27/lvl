package utils

import (
	"fmt"

	"bitbucket.org/level27/lvl/types"
)

// GET /mailgroups
func (c *Client) MailgroupsGetList(get types.CommonGetParams) []types.Mailgroup {
	var response struct {
		Mailgroups []types.Mailgroup `json:"mailgroups"`
	}

	endpoint := fmt.Sprintf("mailgroups?%s", formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "MailgroupsGetList")

	return response.Mailgroups
}

// GET /mailgroups/{mailgroupID}
func (c *Client) MailgroupsGetSingle(mailgroupID int) types.Mailgroup {
	var response struct {
		Mailgroup types.Mailgroup `json:"mailgroup"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d", mailgroupID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "MailgroupsGetSingle")

	return response.Mailgroup
}

func (c *Client) MailgroupsLookup(name string) *types.Mailgroup {
	mailgroups := c.MailgroupsGetList(types.CommonGetParams{Filter: name})
	for _, val := range mailgroups {
		if val.Name == name {
			return &val
		}

		// Check domain names
		for _, domain := range val.Domains {
			fullName := fmt.Sprintf("%s.%s", domain.Name, domain.Domaintype.Extension)
			if fullName == name {
				return &val
			}
		}
	}

	return nil;
}

// POST /mailgroups
func (c *Client) MailgroupsCreate(create types.MailgroupCreate) types.Mailgroup {
	var response struct {
		Mailgroup types.Mailgroup `json:"mailgroup"`
	}

	endpoint := "mailgroups"
	err := c.invokeAPI("POST", endpoint, create, &response)
	AssertApiError(err, "MailgroupsCreate")

	return response.Mailgroup
}

// DELETE /mailgroups/{mailgroupID}
func (c *Client) MailgroupsDelete(mailgroupID int) {
	endpoint := fmt.Sprintf("mailgroups/%d", mailgroupID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)
	AssertApiError(err, "MailgroupsDelete")
}

// PUT /mailgroups/{mailgroupID}
func (c *Client) MailgroupsUpdate(mailgroupID int, data map[string]interface{}) {
	endpoint := fmt.Sprintf("mailgroups/%d", mailgroupID)
	err := c.invokeAPI("PUT", endpoint, data, nil)
	AssertApiError(err, "MailgroupsUpdate")
}

// POST /mailgroups/{mailgroupID}/actions
func (c *Client) MailgroupsAction(mailgroupID int, action string) types.Mailgroup {
	var response struct {
		Mailgroup types.Mailgroup `json:"mailgroup"`
	}

	var request struct {
		Type string `json:"type"`
	}
	request.Type = action;

	endpoint := fmt.Sprintf("mailgroups/%d/actions", mailgroupID)
	err := c.invokeAPI("POST", endpoint, request, &response)
	AssertApiError(err, "MailgroupsAction")

	return response.Mailgroup
}
