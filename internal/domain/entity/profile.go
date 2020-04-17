package entity

type Profile struct {
	ID     string
	UserID string
	//
	Address1 string
	Address2 string
	City     string
	State    string
	Country  string
	Zip      string
	//
	PhotoURL  string
	FirstName string
	LastName  string
	BirthDate uint32
	//
	Language string
}
