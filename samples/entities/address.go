package entities

//go:generate e2t -entity=$GOFILE -table=ADDRESSES

// Simple struct that represents Address entity that will be mapped to table ADDRESSES
// (see go:generate code above)
type Address struct {
	Id int64			`entity:"ID"`

	Zip string			`entity:"ZIP_CODE"`
	StreetName string	`entity:"STREET_NAME"`
}
