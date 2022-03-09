package types

type Network struct {
	NetworkRef
	UID             interface{}     `json:"uid"`
	Remarks         interface{}     `json:"remarks"`
	Status          string          `json:"status"`
	Vlan            interface{}     `json:"vlan"`
	Ipv4            string          `json:"ipv4"`
	Netmaskv4       int             `json:"netmaskv4"`
	Gatewayv4       string          `json:"gatewayv4"`
	Ipv6            string          `json:"ipv6"`
	Netmaskv6       int             `json:"netmaskv6"`
	Gatewayv6       string          `json:"gatewayv6"`
	PublicIP4Native interface{}     `json:"publicIp4Native"`
	PublicIP6Native interface{}     `json:"publicIp6Native"`
	Full            interface{}     `json:"full"`
	Systemgroup     interface{}     `json:"systemgroup"`
	Organisation    OrganisationRef `json:"organisation"`
	Zone            struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Region struct {
			ID int `json:"id"`
		} `json:"region"`
	} `json:"zone"`
	Systemprovider struct {
		ID                 int    `json:"id"`
		API                string `json:"api"`
		Name               string `json:"name"`
		AdvancedNetworking bool   `json:"advancedNetworking"`
	} `json:"systemprovider"`
	Rzone4         interface{}   `json:"rzone4"`
	Rzone6         interface{}   `json:"rzone6"`
	Zones          []interface{} `json:"zones"`
	StatusCategory string        `json:"statusCategory"`
}

type NetworkRef struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Description interface{} `json:"description"`
	Public      bool        `json:"public"`
	Customer    bool        `json:"customer"`
	Internal    bool        `json:"internal"`
}