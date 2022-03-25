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

// POST /mailgroups/{mailgroupID}/domains
func (c *Client) MailgroupsDomainsLink(mailgroupID int, data types.MailgroupDomainAdd) types.Mailgroup {
	var response struct {
		Mailgroup types.Mailgroup `json:"mailgroup"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/domains", mailgroupID)
	err := c.invokeAPI("POST", endpoint, data, &response)
	AssertApiError(err, "MailgroupsDomainsAdd")

	return response.Mailgroup
}

// DELETE /mailgroups/{mailgroupID}/domains/{domainId}
func (c *Client) MailgroupsDomainsUnlink(mailgroupID int, domainId int) {
	endpoint := fmt.Sprintf("mailgroups/%d/domains/%d", mailgroupID, domainId)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)
	AssertApiError(err, "MailgroupsDomainsRemove")
}

// PATCH /mailgroups/{mailgroupID}/domains/{domainId}/setprimary
func (c *Client) MailgroupsDomainsSetPrimary(mailgroupID int, domainId int) {
	endpoint := fmt.Sprintf("mailgroups/%d/domains/%d/setprimary", mailgroupID, domainId)
	err := c.invokeAPI("PATCH", endpoint, nil, nil)
	AssertApiError(err, "MailgroupsDomainsSetPrimary")
}

// PATCH /mailgroups/{mailgroupID}/domains/{domainID}
func (c *Client) MailgroupsDomainsPatch(mailgroupID int, domainID int, data map[string]interface{}) {
	endpoint := fmt.Sprintf("mailgroups/%d/domains/%d", mailgroupID, domainID)
	err := c.invokeAPI("PATCH", endpoint, data, nil)
	AssertApiError(err, "MailgroupsDomainsPatch")
}


// GET /mailgroups/{mailgroupId}/mailboxes
func (c *Client) MailgroupsMailboxesGetList(mailgroupID int, get types.CommonGetParams) []types.MailboxShort {
	var response struct {
		Mailboxes []types.MailboxShort `json:"mailboxes"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes", mailgroupID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "MailgroupsMailboxesGetList")

	return response.Mailboxes
}

// POST /mailgroups/{mailgroupId}/mailboxes
func (c *Client) MailgroupsMailboxesCreate(mailgroupID int, data types.MailboxCreate) types.Mailbox {
	var response struct {
		Mailbox types.Mailbox `json:"mailbox"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes", mailgroupID)
	err := c.invokeAPI("POST", endpoint, data, &response)
	AssertApiError(err, "MailgroupsMailboxesCreate")

	return response.Mailbox
}

// GET /mailgroups/{mailgroupId}/mailboxes/{mailboxId}
func (c *Client) MailgroupsMailboxesGetSingle(mailgroupID int, mailboxID int) types.Mailbox {
	var response struct {
		Mailbox types.Mailbox `json:"mailbox"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d", mailgroupID, mailboxID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "MailgroupsMailboxesGetSingle")

	return response.Mailbox
}

// DELETE /mailgroups/{mailgroupId}/mailboxes/{mailboxId}
func (c *Client) MailgroupsMailboxesDelete(mailgroupID int, mailboxID int)  {
	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d", mailgroupID, mailboxID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)
	AssertApiError(err, "MailgroupsMailboxesDelete")
}

// PUT /mailgroups/{mailgroupId}/mailboxes
func (c *Client) MailgroupsMailboxesUpdate(mailgroupID int, mailboxID int, data map[string]interface{}) {
	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d", mailgroupID, mailboxID)
	err := c.invokeAPI("PUT", endpoint, data, nil)
	AssertApiError(err, "MailgroupsMailboxesUpdate")
}


func (c *Client) MailgroupsMailboxesLookup(mailgroupID int, name string) *types.MailboxShort {
	mailgroups := c.MailgroupsMailboxesGetList(mailgroupID, types.CommonGetParams{Filter: name})
	for _, val := range mailgroups {
		if val.Name == name || val.Username == name {
			return &val
		}
	}

	return nil;
}

// GET /mailgroups/{mailgroupId}/mailboxes/{mailboxId}/addresses
func (c *Client) MailgroupsMailboxesAddressesGetList(mailgroupID int, mailboxID int, get types.CommonGetParams) []types.MailboxAddress {
	var response struct {
		MailboxAddresses []types.MailboxAddress `json:"mailboxAddresses"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d/addresses", mailgroupID, mailboxID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "MailgroupsMailboxesAddressesGetList")

	return response.MailboxAddresses
}

// POST /mailgroups/{mailgroupId}/mailboxes/{mailboxId}/addresses
func (c *Client) MailgroupsMailboxesAddressesCreate(mailgroupID int, mailboxID int, data types.MailboxAddressCreate) types.MailboxAddress {
	var response struct {
		MailboxAddress types.MailboxAddress `json:"mailboxAdress"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d/addresses", mailgroupID, mailboxID)
	err := c.invokeAPI("POST", endpoint, data, &response)
	AssertApiError(err, "MailgroupsMailboxesAddressesCreate")

	return response.MailboxAddress
}

// GET /mailgroups/{mailgroupId}/mailboxes/{mailboxId}/addresses/{addressId}
func (c *Client) MailgroupsMailboxesAddressesGetSingle(mailgroupID int, mailboxID int, addressID int) types.MailboxAddress {
	var response struct {
		MailboxAddress types.MailboxAddress `json:"mailboxAddress"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d/addresses/%d", mailgroupID, mailboxID, addressID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "MailgroupsMailboxesAddressesGetSingle")

	return response.MailboxAddress
}

// DELETE /mailgroups/{mailgroupId}/mailboxes/{mailboxId}/addresses/{addressId}
func (c *Client) MailgroupsMailboxesAddressesDelete(mailgroupID int, mailboxID int, addressID int)  {
	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d/addresses/%d", mailgroupID, mailboxID, addressID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)
	AssertApiError(err, "MailgroupsMailboxesAddressesDelete")
}

// PUT /mailgroups/{mailgroupId}/mailboxes/addresses/{addressId}
func (c *Client) MailgroupsMailboxesAddressesUpdate(mailgroupID int, mailboxID int, addressID int, data map[string]interface{}) {
	endpoint := fmt.Sprintf("mailgroups/%d/mailboxes/%d/addresses/%d", mailgroupID, mailboxID, addressID)
	err := c.invokeAPI("PUT", endpoint, data, nil)
	AssertApiError(err, "MailgroupsMailboxesAddressesUpdate")
}

func (c *Client) MailgroupsMailboxesAddressesLookup(mailgroupID int, mailboxID int, address string) *types.MailboxAddress {
	addresses := c.MailgroupsMailboxesAddressesGetList(mailgroupID, mailboxID, types.CommonGetParams{Filter: address})
	for _, val := range addresses {
		if val.Address == address {
			return &val
		}
	}

	return nil;
}


// GET /mailgroups/{mailgroupId}/mailforwarders
func (c *Client) MailgroupsMailforwardersGetList(mailgroupID int, get types.CommonGetParams) []types.Mailforwarder {
	var response struct {
		Mailforwarders []types.Mailforwarder `json:"mailforwarders"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailforwarders", mailgroupID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "MailgroupsMailforwardersGetList")

	return response.Mailforwarders
}

// POST /mailgroups/{mailgroupId}/mailforwarders
func (c *Client) MailgroupsMailforwardersCreate(mailgroupID int, data types.MailforwarderCreate) types.Mailforwarder {
	var response struct {
		Mailforwarder types.Mailforwarder `json:"mailforwarder"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailforwarders", mailgroupID)
	err := c.invokeAPI("POST", endpoint, data, &response)
	AssertApiError(err, "MailgroupsMailforwardersCreate")

	return response.Mailforwarder
}

// GET /mailgroups/{mailgroupId}/mailforwarders/{mailforwarderId}
func (c *Client) MailgroupsMailforwardersGetSingle(mailgroupID int, mailforwarderID int) types.Mailforwarder {
	var response struct {
		Mailforwarder types.Mailforwarder `json:"mailforwarder"`
	}

	endpoint := fmt.Sprintf("mailgroups/%d/mailforwarders/%d", mailgroupID, mailforwarderID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "MailgroupsMailforwardersGetSingle")

	return response.Mailforwarder
}

// DELETE /mailgroups/{mailgroupId}/mailforwarders/{mailforwarderId}
func (c *Client) MailgroupsMailforwardersDelete(mailgroupID int, mailforwarderID int)  {
	endpoint := fmt.Sprintf("mailgroups/%d/mailforwarders/%d", mailgroupID, mailforwarderID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)
	AssertApiError(err, "MailgroupsMailforwardersDelete")
}

// PUT /mailgroups/{mailgroupId}/mailforwarders
func (c *Client) MailgroupsMailforwardersUpdate(mailgroupID int, mailforwarderID int, data map[string]interface{}) {
	endpoint := fmt.Sprintf("mailgroups/%d/mailforwarders/%d", mailgroupID, mailforwarderID)
	err := c.invokeAPI("PUT", endpoint, data, nil)
	AssertApiError(err, "MailgroupsMailforwardersUpdate")
}


func (c *Client) MailgroupsMailforwardersLookup(mailgroupID int, name string) *types.Mailforwarder {
	mailgroups := c.MailgroupsMailforwardersGetList(mailgroupID, types.CommonGetParams{Filter: name})
	for _, val := range mailgroups {
		if val.Address == name {
			return &val
		}
	}

	return nil;
}

