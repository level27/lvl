package types

type AclAdd struct {
	Organisation int `json:"organisation"`
}

type Acl struct {
	ID           int             `json:"id"`
	Object       string          `json:"object"`
	ObjectID     int             `json:"objectId"`
	Permissions  interface{}     `json:"permissions"`
	Extra        interface{}     `json:"extra"`
	Type         string          `json:"type"`
	Organisation OrganisationRef `json:"organisation"`
}