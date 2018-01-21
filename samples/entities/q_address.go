/*
	CAUTION: GENERATED FILE!!! DO NOT EDIT!!!
*/
package entities

type qAddress struct {
	Id string
	Zip string
	StreetName string
}

var QAddress = qAddress {
	Id: "ID",
	Zip: "ZIP_CODE",
	StreetName: "STREET_NAME",
}

func (qAddress) TableName() string {
	return "ADDRESSES"
}
