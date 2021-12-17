package types

// ALL USER RELATED TYPES

type User struct {
	ID        int      		`json:"id"`
	Username  string   		`json:"username"`
	Email     string   		`json:"email"`
	FirstName string   		`json:"test"`
	LastName  string   		`json:"user"`
	Roles     []string 		`json:"roles"`
	Status    string   		`json:"status"`
	StatusCategory string	`json:"statusCategory"`
	Language	string		`json:"language"`
	Fullname  	string   	`json:"fullname"`
}

type Contact struct{
	ID				int 	`json:"id"`
	DtStamp			string	`json:"dtStamp"`
	FullName		string 	`json:"fullName"`
	Language		string 	`json:"language"`
	Message			string 	`json:"message"`
	Status			int 	`json:"status"`
	Type			string	`json:"type"`
	Value			string	`json:"value"`
	ContactID		int		`json:"contactId"`	
}