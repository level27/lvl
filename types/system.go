package types

type SystemGet struct {
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
	SystemImage struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		ExternalId  string `json:"externalId"`
		OsId        int    `json:"osId"`
		OsName      string `json:"osName"`
		OsType      string `json:"osType"`
		OsVersion   string `json:"osVersion"`
		OsVersionId int    `json:"osVersionId"`
	} `json:"systemimage"`
	OperatingSystemVersion struct {
		Id        int    `json:"id"`
		OsId      int    `json:"osId"`
		OsName    string `json:"osName"`
		OsType    string `json:"osType"`
		OsVersion string `json:"osVersion"`
	} `json:"operatingsystemVersion"`
	ProvideId                   int    `json:"providerId"`
	Provider                    string `json:"provider"`
	ProviderApi                 string `json:"providerApi"`
	SystemProviderConfiguration struct {
		Id          int    `json:"id"`
		ExternalId  string `json:"externalId"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"systemproviderConfiguration"`
	Region string `json:"region"`
	Zone   struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"zone"`
	Networks         []interface{} `json:"networks"`
	PublicNetworking bool          `json:"publicNetworking"`
	StatsSummary     struct {
		DiskSpace struct {
			Unit  string `json:"unit"`
			Value string `json:"value"`
			Max   string `json:"max"`
		} `json:"diskspace"`
		Memory struct {
			Unit  string `json:"unit"`
			Value string `json:"value"`
			Max   string `json:"max"`
		} `json:"Memory"`
		Cpu struct {
			Unit  string `json:"unit"`
			Value string `json:"value"`
			Max   string `json:"max"`
		} `json:"cpu"`
	} `json:"statsSummary"`
	DtExpires     int    `json:"dtExpires"`
	BillingStatus string `json:"billingStatus"`
	ExternalInfo  string `json:"externalInfo"`
	Remarks       string `json:"remarks"`
}
