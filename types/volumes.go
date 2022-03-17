package types

type Volume struct {
	ID             int             `json:"id"`
	Name           string          `json:"name"`
	Status         string          `json:"status"`
	Space          int             `json:"space"`
	UID            string          `json:"uid"`
	Remarks        interface{}     `json:"remarks"`
	AutoResize     bool            `json:"autoResize"`
	DeviceName     string          `json:"deviceName"`
	Organisation   OrganisationRef `json:"organisation"`
	System         SystemRef       `json:"system"`
	Volumegroup    VolumegroupRef  `json:"volumegroup"`
	StatusCategory string          `json:"statusCategory"`
}

type VolumegroupRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type VolumeCreate struct {
	Name         string `json:"name"`
	Space        int    `json:"space"`
	Organisation int    `json:"organisation"`
	System       int    `json:"system"`
	Volumegroup  *int   `json:"volumegroup"`
	AutoResize   bool   `json:"autoResize"`
	DeviceName   string `json:"deviceName"`
}

type VolumePut struct {
	Name         string      `json:"name"`
	DeviceName   string      `json:"deviceName"`
	Space        int         `json:"space"`
	Organisation int         `json:"organisation"`
	AutoResize   bool        `json:"autoResize"`
	Remarks      interface{} `json:"remarks"`
	System       int         `json:"system"`
	Volumegroup  int         `json:"volumegroup"`
}