package types

// main structure of an app
type App struct {
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

//type to create an app (post request)
type AppPostRequest struct {
	Name         string      `json:"name"`
	Organisation int         `json:"organisation"`
	AutoTeams    interface{} `json:"autoTeams"`
	ExternalInfo *string     `json:"externalInfo"`
}
