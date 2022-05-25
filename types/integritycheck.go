package types

type IntegrityCreateRequest struct {
	Dojobs    bool `json:"dojobs"`
	Forcejobs bool `json:"forcejobs"`
}

type IntegrityCheck struct {
	Id          int    `json:"id"`
	DtRequested string `json:"dtRequested"`
	Status      string `json:"status"`
}