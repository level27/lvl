package types

// ----------------------------------- SYSTEMGROUPS ----------------------------------
// structure of a system group returned by API
type Systemgroup struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Shared  bool   `json:"shared"`
	Systems []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"systems"`
	Organisation struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"organisation"`
}

// request type for creating systemgroup.
type SystemgroupRequest struct {
	Name         string `json:"name"`
	Organisation int    `json:"organisation"`
}
