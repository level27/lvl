package types

type system struct {
	Id                    int    `json:"id"`
	Uid                   string `json:"uid"`
	Hostname              string `json:"hostname"`
	Fqdn                  string `json:"fqdn"`
	Name                  string `json:"name"`
	Type                  string `json:"type"`
	Status                string `json:"status"`
	StatusCategory        string `json:"statusCategory"`
	RunningStatus         string `json:"runningStatus"`
	RunningStatusCategory string `json:"runningStatusCategory"`
	Cpu                   int    `json:"cpu"`
	Memory                int    `json:"memory"`
	Disk                  string `json:"disk"`
	MonitoringEnabled     bool   `json:"monitoringEnabled"`
	ManagementType        string `json:"managementType"`
	Organisation          struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"organisation"`
	
}
