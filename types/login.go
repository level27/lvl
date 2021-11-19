package types

type Login struct {
	Success bool `json:"success"`
	User    struct {
		ID           int      `json:"id"`
		Username     string   `json:"username"`
		Email        string   `json:"email"`
		Firstname    string   `json:"firstName"`
		Lastname     string   `json:"lastName"`
		Roles        []string `json:"roles"`
		Status       string   `json:"status"`
		Language     string   `json:"language"`
		Organisation struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Street      string `json:"street"`
			Housenumber string `json:"houseNumber"`
			Zip         string `json:"zip"`
			City        string `json:"city"`
		} `json:"organisation"`
		Country struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"country"`
		Fullname string `json:"fullname"`
	} `json:"user"`
	Hash string `json:"hash"`
}
