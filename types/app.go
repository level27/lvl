package types

type StructApp struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Status        string `json:"status"`
	Organisations struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"organisations"`
	DtExpires     int    `json:"dtExpires"`
	BillingStatus string `json:"billingStatus"`
	Components    []struct {
		ID               int    `json:"id"`
		Name             string `json:"name"`
		Category         string `json:"category"`
		AppComponentType string `json:"appcomponenttype"`
	} `json:"components"`
}

type App struct {
	Data StructApp `json:"app"`
}

type Apps struct {
	Data []StructApp `json:"apps"`
}
