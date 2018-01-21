package main

import (
	"testing"
	"go/parser"
	"go/token"
	"fmt"
	"reflect"
	"e2t/samples/entities"
	"strings"
	"os"
)

const samplePackageName = "entities"

func TestParseSource(t *testing.T) {
	f, _ := parser.ParseFile(
		token.NewFileSet(),
		fmt.Sprintf(fmt.Sprintf("samples/%s/user.go", samplePackageName)),
		nil,
		parser.ParseComments,
	)

	userType := reflect.TypeOf(entities.User{})
	sm := parseSource(f)

	if !strings.Contains(userType.PkgPath(), sm.Package) {
		t.Fatal(
			fmt.Sprintf("Invalid struct package: expected '%s', got '%s'",
				samplePackageName,
				sm.Package,
			),
		)
	}

	if sm.Name != userType.Name() {
		t.Fatal(
			fmt.Sprintf("Invalid struct name: expected '%s', got '%s'",
				userType.Name(),
				sm.Name,
			),
		)
	}

	if len(sm.Fields) != userType.NumField() {
		t.Fatal(
			fmt.Sprintf("Invalid field count: expected '%s', got '%s'",
				userType.NumField(),
				len(sm.Fields),
			),
		)
	}

	for _, fieldMetadata := range sm.Fields {
		userTypeField, ok := userType.FieldByName(fieldMetadata.Name)

		if !ok {
			t.Fatal(
				fmt.Sprintf("Field '%s' is absent in struct",
					fieldMetadata.Name,
				),
			)
		}

		for _, tagMetadata := range fieldMetadata.Tags {
			tag, ok := userTypeField.Tag.Lookup(tagMetadata.Name)

			if !ok {
				t.Fatal(
					fmt.Sprintf("Struct field '%s' does not have tag '%s'",
						fieldMetadata.Name,
						tagMetadata.Name,
					),
				)
			} else if tag != tagMetadata.Value {
				t.Fatal(
					fmt.Sprintf("Invalid tag value was detected: expected '%s', got '%s'",
						tag,
						tagMetadata.Value,
					),
				)
			}
		}
	}
}

func ExampleGeneratedUserEntityMetadataSourceCode() {
	tableName := "USERS"

	f, _ := parser.ParseFile(
		token.NewFileSet(),
		fmt.Sprintf(fmt.Sprintf("samples/%s/user.go", samplePackageName)),
		nil,
		parser.ParseComments,
	)

	tplt.Execute(os.Stdout, &entityMetadata {
		TableName: &tableName,
		Struct: parseSource(f),
	})

	// Output:
	///*
	//	CAUTION: GENERATED FILE!!! DO NOT EDIT!!!
	//*/
	//package entities
	//
	//type qUser struct {
	//	Id string
	//	FirstName string
	//	LastName string
	//	Address qAddress
	//}
	//
	//var QUser = qUser {
	//	Id: "ID",
	//	FirstName: "FIRST_NAME",
	//	LastName: "LAST_NAME",
	//	Address: QAddress,
	//}
	//
	//func (qUser) TableName() string {
	//	return "USERS"
	//}
}