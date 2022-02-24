package types

type Job struct {
	Action  string        `json:"action"`
	Dt      interface{}   `json:"dt"`
	Eclass  string        `json:"eClass"`
	Eid     int           `json:"eId"`
	Estring string        `json:"eString"`
	ExcCode int           `json:"excCode"`
	ExcMsg  string        `json:"excMsg"`
	Hoe     int           `json:"hoe"`
	Id      int           `json:"id"`
	Jobs    []Job         `json:"jobs"`
	Logs    []interface{} `json:"logs"`
	Message string        `json:"msg"`
	Service string        `json:"service"`
	Status  int           `json:"status"`
	System  int           `json:"system"`
}
