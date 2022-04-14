package types

// main structure of an app
type App struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	Organisation struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Reseller string `json:"reseller"`
	} `json:"organisation"`
	DtExpires     int    `json:"dtExpires"`
	BillingStatus string `json:"billingStatus"`
	Components    []struct {
		ID               int    `json:"id"`
		Name             string `json:"name"`
		Category         string `json:"category"`
		AppComponentType string `json:"appcomponenttype"`
	} `json:"components"`
	CountTeams int `json:"countTeams"`
	Teams      []struct {
		ID             int    `json:"id"`
		Name           string `json:"name"`
		AdminOnly      bool   `json:"adminOnly"`
		OrganisationID int    `json:"organisationId"`
	} `json:"teams"`
	ExternalInfo string `json:"externalInfo"`
}

//type to create an app (post request)
type AppPostRequest struct {
	Name         string `json:"name"`
	Organisation int    `json:"organisation"`
	AutoTeams    []int  `json:"autoTeams"`
	ExternalInfo string `json:"externalInfo"`
}

//type to update an app (put request)
type AppPutRequest struct {
	Name         string   `json:"name"`
	Organisation int      `json:"organisation"`
	AutoTeams    []string `json:"autoTeams"`
}

// type needed to do an action on a system
type AppActionRequest struct {
	Type string `json:"type"`
}

//type appcomponent
type AppComponent2 struct {
	App struct {
		ID             int64  `json:"id"`
		Status         string `json:"status"`
		Name           string `json:"name"`
		StatusCategory string `json:"statusCategory"`
	} `json:"app"`
	AppcomponentparameterDescriptions interface{} `json:"appcomponentparameterDescriptions"`
	Appcomponentparameters            interface{} `json:"appcomponentparameters"`
	Appcomponenttype                  string      `json:"appcomponenttype"`
	BillableitemDetailID              int64       `json:"billableitemDetailId"`
	Category                          string      `json:"category"`
	ID                                int64       `json:"id"`
	Name                              string      `json:"name"`
	Organisation                      struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"organisation"`
	Provider struct {
		ID   interface{} `json:"id"`
		Name interface{} `json:"name"`
	} `json:"provider"`
	SelectedSystem interface{} `json:"selectedSystem"`
	Status         string      `json:"status"`
	Systemgroup    interface{} `json:"systemgroup"`
	Systems        []struct {
		Cookbooks []interface{} `json:"cookbooks"`
		Fqdn      string        `josn:"fqdn"`
		ID        int64         `json:"id"`
		Name      string        `json:"name"`
	} `json:"systems"`
}

// type appcomponent category
type AppcomponentCategory struct {
	Name string
}

// type appcomponenttype
type Appcomponenttype map[string]AppcomponenttypeServicetype

type AppcomponenttypeServicetype struct{
	Servicetype struct {
		Name                    string        `json:"name"`
		Cookbook                string        `json:"cookbook"`
		DisplayName             string        `json:"displayName"`
		Description             string        `json:"description"`
		URLPossible             bool          `json:"urlPossible"`
		RestorePossible         bool          `json:"restorePossible"`
		MigrationPossible       bool          `json:"migrationPossible"`
		SelectingSystemPossible bool          `json:"selectingSystemPossible"`
		DisabledOnProduction    bool          `json:"disabledOnProduction"`
		InvisibleOnProduction   bool          `json:"invisibleOnProduction"`
		Runlist                 string        `json:"runlist"`
		AllowedActions          []interface{} `json:"allowedActions"`
		Category                string        `json:"category"`
		Parameters              []struct {
			Name         string `json:"name"`
			DisplayName  string `json:"displayName"`
			Description  string `json:"description"`
			Type         string `json:"type"`
			DefaultValue string `json:"defaultValue"`
			Readonly     bool   `json:"readonly"`
			DisableEdit  bool   `json:"disableEdit"`
			Required     bool   `json:"required"`
			Category     string `json:"category"`
		} `json:"parameters"`
	} `json:"servicetype"`
}


