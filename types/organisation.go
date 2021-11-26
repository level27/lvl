package types

type StructOrganisation struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	TaxNumber   string `json:"taxNumber"`
	MustPayTax  bool   `json:"mustPayTax"`
	Street      string `json:"street"`
	HouseNumber string `json:"houseNumber"`
	Zip         string `json:"zip"`
	City        string `json:"city"`
	Country     struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"country"`
	// ResellerOrganisation
	Users []struct {
		ID        int      `json:"id"`
		Username  string   `json:"name"`
		Email     string   `json:"email"`
		FirstName string   `json:"firstName"`
		LastName  string   `json:"lastName"`
		Roles     []string `json:"roles"`
	} `json:"users"`
	// RemarksToprintInvoice
	UpdateEntitiesOnly bool `json:"updateEntitiesOnly"`
}

type Organisation struct {
	Data StructOrganisation `json:"organisation"`
}

type Organisations struct {
	Data []StructOrganisation `json:"organisations"`
}