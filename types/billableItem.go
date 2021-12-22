package types

type BillableItem struct {
	ID           int `json:"id"`
	Organisation struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"organisation"`
	PreventDeactivation bool        `json:"preventDeactivation"`
	Status              int         `json:"status"`
	StatusDisplay       string      `json:"statusDisplay"`
	Description         string      `json:"description"`
	AutoRenew           bool        `json:"autoRenew"`
	DtExpires           interface{} `json:"dtExpires"`
	DtNextRenewal       int         `json:"dtNextRenewal"`
	DocumentsExist      bool        `json:"documentsExist"`
	TotalPrice          int         `json:"totalPrice"`
	Details             []struct {
		ManuallyAdded        interface{} `json:"manuallyAdded"`
		AllowToSkipInvoicing bool        `json:"allowToSkipInvoicing"`
		ID                   int         `json:"id"`
		Price                int         `json:"price"`
		DtExpires            interface{} `json:"dtExpires"`
		Quantity             int         `json:"quantity"`
		Description          string      `json:"description"`
		Product              struct {
			ID                  string `json:"id"`
			Description         string `json:"description"`
			AllowQuantityChange bool   `json:"allowQuantityChange"`
		} `json:"product"`
		ProductPrice struct {
			ID       int    `json:"id"`
			Period   int    `json:"perion"`
			Currency string `json:"currency"`
			Price    string    `json:"price"`
			Timing   string `json:"timing"`
			Status   int    `json:"status"`
		} `json:"productPrice"`
		Type int `json:"Type"`
	} `json:"details"`
	Extra1       string `json:"extra1"`
	Extra2       string `json:"extra2"`
	ExternalInfo string `json:"externalInfo"`
	Agreement    struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"agreement"`
}

// returns the billable item for GET call
type BillableItemGet struct {
	BillableItem BillableItem `json:"billableitem"`
}
