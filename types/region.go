package types

type Region struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"country"`
	Systemprovider struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		API  string `json:"api"`
	} `json:"systemprovider"`
}

type Zone struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
}

type Image struct {
	ID                     int    `json:"id"`
	Name                   string `json:"name"`
	OperatingsystemVersion struct {
		ID              int    `json:"id"`
		Version         string `json:"version"`
		Type            string `json:"type"`
		Operatingsystem struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"operatingsystem"`
	} `json:"operatingsystemVersion"`
}