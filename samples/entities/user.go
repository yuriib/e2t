package entities

//go:generate e2t -entity=$GOFILE -table=USERS

// Simple struct that represents User entity that will be mapped to table USERS
// (see go:generate code above)
type User struct {
	Id int64			`entity:"ID"`

	FirstName string	`entity:"FIRST_NAME"`
	LastName string		`entity:"LAST_NAME"`

	Address Address 	`join-entity:"Address"`

	SomeUnmappedField string
	someUnmappedField string
}
