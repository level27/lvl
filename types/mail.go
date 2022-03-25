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

type MailboxShort struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Status     string `json:"status"`
	OooEnabled bool   `json:"oooEnabled"`
	OooSubject string `json:"oooSubject"`
	OooText    string `json:"oooText"`
	Mailgroup  struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"mailgroup"`
	StatusCategory string `json:"statusCategory"`
	PrimaryAddress string `json:"primaryAddress"`
	Aliases        int    `json:"aliases"`
}

type Mailbox struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Status     string `json:"status"`
	OooEnabled bool   `json:"oooEnabled"`
	OooSubject string `json:"oooSubject"`
	OooText    string `json:"oooText"`
	Source     string `json:"source"`
	Mailgroup  struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"mailgroup"`
	System struct {
		ID       int    `json:"id"`
		Fqdn     string `json:"fqdn"`
		Hostname string `json:"hostname"`
	} `json:"system"`
	BillableitemDetail struct {
		ID int `json:"id"`
	} `json:"billableitemDetail"`
	StatusCategory string `json:"statusCategory"`
	PrimaryAddress string `json:"primaryAddress"`
	Aliases        int    `json:"aliases"`
}

type MailboxCreate struct {
	Name       string `json:"name"`
	Password   string `json:"password"`
	OooEnabled bool   `json:"oooEnabled"`
	OooSubject string `json:"oooSubject"`
	OooText    string `json:"oooText"`
}

type MailboxPut struct {
	Name       string `json:"name"`
	Password   string `json:"password"`
	OooEnabled bool   `json:"oooEnabled"`
	OooSubject string `json:"oooSubject"`
	OooText    string `json:"oooText"`
}

type MailboxDescribe struct {
	Mailbox
	Addresses []MailboxAddress `json:"addresses"`
}

type MailboxAddress struct {
	ID      int    `json:"id"`
	Address string `json:"address"`
	Status  string `json:"status"`
}

type MailboxAddressCreate struct {
	Address string `json:"address"`
}

type Mailforwarder struct {
	ID          int      `json:"id"`
	Address     string   `json:"address"`
	Destination []string `json:"destination"`
	Status      string   `json:"status"`
	Mailgroup   struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"mailgroup"`
	Domain struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Domaintype struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"domaintype"`
	} `json:"domain"`
}

type MailforwarderCreate struct {
	Address     string `json:"address"`
	Destination string `json:"destination"`
}

type MailforwarderPut struct {
	Address     string `json:"address"`
	Destination string `json:"destination"`
}