package types

type Mailgroup struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Systemgroup struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"systemgroup"`
	Organisation  OrganisationRef `json:"organisation"`
	BillingStatus string          `json:"billingStatus"`
	DtExpires     int             `json:"dtExpires"`
	Domains       []struct {
		ID          int             `json:"id"`
		Name        string          `json:"name"`
		MailPrimary bool            `json:"mailPrimary"`
		Domaintype  DomainExtension `json:"domaintype"`
	} `json:"domains"`
	ExternalInfo       interface{} `json:"externalInfo"`
	StatusCategory     string      `json:"statusCategory"`
	MailboxCount       int         `json:"mailboxCount"`
	MailforwarderCount int         `json:"mailforwarderCount"`
}

type MailgroupCreate struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Organisation int    `json:"organisation"`
	Systemgroup  int    `json:"systemgroup"`
	AutoTeams    string `json:"autoTeams"`
	ExternalInfo string `json:"externalInfo"`
}

type MailgroupPut struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Organisation int    `json:"organisation"`
	Systemgroup  int    `json:"systemgroup"`
	AutoTeams    string `json:"autoTeams"`
}

type MailgroupDomainAdd struct {
	Domain        int  `json:"domain"`
	HandleMailDns bool `json:"handleMailDns"`
}