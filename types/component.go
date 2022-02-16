package types

import "github.com/Jeffail/gabs/v2"

type Component struct {
	ID                     int    `json:"id"`
	Name                   string `json:"name"`
	AppComponentType       string `json:"appcomponenttype"`
	Status                 string `json:"status"`
	StatusCategory         string `json:"statuscategory"`
	AppComponentParameters gabs.Container `json:"appcomponentparameters"`
	App                    struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Status string `json:"status"`
	} `json:"apps"`
	Systems []struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Fqdn      string `json:"fqdn"`
		Cookbooks []struct {
			ID             int      `json:"id"`
			CookbookType   string   `json:"cookbooktype"`
			Status         string   `json:"status"`
			StatusCategory string   `json:"statuscategory"`
			Versions       []string `json:"versions"`
		} `json:"cookbooks"`
	} `json:"systems"`
	// SystemGroup
	// ContainerSystem
	// Provider
	LinkedUrlsCount int `json:"linkedurlscount"`
}
