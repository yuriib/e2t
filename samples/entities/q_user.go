/*
	CAUTION: GENERATED FILE!!! DO NOT EDIT!!!
*/
package entities

type qUser struct {
	Id string
	FirstName string
	LastName string
	Address qAddress
}

var QUser = qUser {
	Id: "ID",
	FirstName: "FIRST_NAME",
	LastName: "LAST_NAME",
	Address: QAddress,
}

func (qUser) TableName() string {
	return "USERS"
}
