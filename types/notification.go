package types

import "encoding/json"

type Notification struct {
	// Data here is stored in an anonymous struct
	// so we can deserialize it separately from Entity in UnmarshalJSON() below.
	NotificationData
	Entity interface{} `json:"entity"`
}

type NotificationData struct {
	ID                int         `json:"id"`
	EnitityIndex      string      `json:"entityIndex"`
	EntityName        string      `json:"entityName"`
	DtStamap          string      `json:"dtStamp"`
	NotificationGroup string      `json:"notificationGroup"`
	Type              string      `json:"type"`
	EntityClass       string      `json:"entityClass"`
	EntityID          int         `json:"entityId"`
	RootEntityClass   string      `json:"rootEntityClass"`
	RootEntityID      int         `json:"rootEntityId"`
	Status            int         `json:"status"`
	StatusDisplay     string      `json:"statusDisplay"`
	StatusCategory    string      `json:"statusCategory"`
	SendMode          int         `json:"sendMode"`
	Priority          int         `json:"priority"`
	Subject           interface{} `json:"subject"`
	Params            interface{} `json:"params"`
	UserID            int         `json:"userId"`
	Contacts          []struct {
		ID        int    `json:"id"`
		DtStamp   string `json:"dtStamp"`
		FullName  string `json:"fullName"`
		Language  string `json:"language"`
		Message   string `json:"message"`
		Status    int    `json:"status"`
		Type      string `json:"type"`
		Value     string `json:"value"`
		ContactID int    `json:"contactId"`
	} `json:"contacts"`
	ExtraRecipients []string    `json:"extraRecipients"`
	User            struct {
		ID               int      `json:"id"`
		Username         string   `json:"username"`
		Email            string   `json:"email"`
		FirstName        string   `json:"firstName"`
		LastName         string   `json:"lastName"`
		Fullname         string   `json:"fullname"`
		Roles            []string `json:"roles"`
		Status           string   `json:"status"`
		StatusCategory   string   `json:"statusCategory"`
		Language         string   `json:"language"`
		WebsiteOrderInfo string   `json:"websiteOrderInfo"`
		Organisation     struct {
			ID                 int    `json:"id"`
			Name               string `json:"name"`
			Street             string `json:"street"`
			HouseNumber        string `json:"houseNumber"`
			Zip                string `json:"zip"`
			City               string `json:"city"`
			Reseller           string `json:"reseller"`
			UpdateEntitiesOnly bool   `json:"updateEntitiesOnly"`
		} `json:"organisation"`
		Country struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"country"`
	} `json:"user"`
}

func (n *Notification) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &n.NotificationData)
	if err != nil {
		return err
	}

	// The type of the Entity field is based on the value of EntityName.
	// We have to deserialize the main struct before we can deserialize Entity.

	switch (n.EntityName) {
	case "domain":
		var dat struct { Entity Domain `json:"entity"` }
		err = json.Unmarshal(data, &dat)
		n.Entity = dat.Entity
	}

	return err
}