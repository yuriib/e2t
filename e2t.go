package main

import (
	"go/parser"
	"go/token"
	"go/ast"
	"regexp"
	"strings"
	"os"
	"flag"
	"log"
	"fmt"
	"text/template"
)

const entityFlag = "entity"
const tableFlag = "table"

const entityTag = "entity"
const joinEntityTag = "join-entity"

var tagParser = regexp.MustCompile(`(.+):"(.+)"`)

type entityMetadata struct {
	TableName *string

	Struct *structMetadata
}

type structMetadata struct {
	Package string

	Name string
	Fields structFieldsMetadata
}

type structFieldMetadata struct {
	Name string
	Type string
	Tags structFieldTagsMetadata
}

type structFieldTagMetadata struct {
	Name string
	Value string
}

type structFieldsMetadata []structFieldMetadata
type structFieldTagsMetadata []structFieldTagMetadata

var qTemplate =
`/*
	CAUTION: GENERATED FILE!!! DO NOT EDIT!!!
*/

{{- $tblName := .TableName}}
{{- with .Struct}}
package {{.Package}}

type q{{.Name}} struct {
	{{- range .Fields}}
	{{- if .Tags}}
	{{.Name}} {{.Tags | extractType}}
	{{- end}}
	{{- end}}
}

var Q{{.Name}} = q{{.Name}} {
	{{- range .Fields}}
	{{- if .Tags}}
	{{.Name}}: {{.Tags | extractValue}},
	{{- end}}
	{{- end}}
}

func (q{{.Name}}) TableName() string {
	return "{{$tblName}}"
}
{{- end}}
`

var tplt = template.Must(
	template.
	New("qTemplate").
		Funcs(
			template.FuncMap {
				"extractType": extractType,
				"extractValue": extractValue,
			},
		).
		Parse(qTemplate),
)

var (
	entityFileName = flag.String(entityFlag, "", "Entity file name (i.e.: user.go or $GOFILE)")
	tableName = flag.String(tableFlag, "", "Table which will be mapped to entity (i.e.: USERS)")
)

func main() {
	flag.Parse()

	validateInput()

	path := handleGetpw(os.Getwd())

	f := handleParseFile(
		parser.ParseFile(
			token.NewFileSet(),
			fmt.Sprintf("%s/%s", path, *entityFileName),
			nil,
			parser.ParseComments,
		),
	)

	qf := handleCreateFile(
		os.Create(
			fmt.Sprintf(
				"%s/q_%s",
				path,
				*entityFileName,
			),
		),
	)

	tplt.Execute(qf, &entityMetadata {
		TableName: tableName,
		Struct: parseSource(f),
	})
}

func validateInput() {
	if *entityFileName == "" {
		handleMissingFlag(entityFlag)
	}

	if *tableName == "" {
		handleMissingFlag(tableFlag)
	}

	if strings.Contains(*entityFileName, "/") {
		handleInvalidFileName(*entityFileName)
	}
}

func handleMissingFlag(flag string) {
	log.Fatal("Missing flag: -", flag)

	os.Exit(1)
}

func handleInvalidFileName(fileName string) {
	log.Fatal("Invalid file name was provided: ", fileName)

	os.Exit(1)
}

func handleGetpw(dir string, err error) string {
	if err != nil {
		log.Fatal(err)

		os.Exit(1)
	}

	return dir
}

func handleParseFile(f *ast.File, err error) *ast.File {
	if err != nil {
		log.Fatal(err)

		os.Exit(1)
	}

	return f
}

func handleCreateFile(f *os.File, err error) *os.File {
	if err != nil {
		log.Fatal(err)

		os.Exit(1)
	}

	return f
}

func parseSource(f *ast.File) *structMetadata {
	sm := new(structMetadata)

	ast.Inspect(f, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.File:
			sm.Package = x.Name.Name
			break
		case *ast.TypeSpec:
			structType, ok := x.Type.(*ast.StructType)

			// we're ok as we know how to parse struct
			if ok {
				sm.Name = x.Name.Name
				sm.Fields = parseStructFields(structType.Fields.List)
			}

			// we've done with parsing structure
			return false
		}

		return true
	})

	return sm
}

func parseStructFields(fields []*ast.Field) []structFieldMetadata {
	sfm := make([]structFieldMetadata, len(fields))

	for i, field := range fields {
		aType := field.Type.(*ast.Ident)

		// grab first name only. Do not know when we can have more than 1 field name
		for _, name := range field.Names {
			sfm[i] = structFieldMetadata {
				Name: name.Name,
				Type: aType.Name,
				Tags: parseStructFieldTags(field.Tag),
			}

			break
		}
	}

	return sfm
}

func parseStructFieldTags(rawTags *ast.BasicLit) []structFieldTagMetadata {
	if rawTags == nil {
		return nil
	}

	tags := strings.Split(strings.Trim(rawTags.Value, "`"), " ")
	sftm := make([]structFieldTagMetadata, len(tags))

	for i, tag := range tags {
		parsedTag := tagParser.FindAllStringSubmatch(tag, -1)

		sftm[i] = structFieldTagMetadata {
			Name: parsedTag[0][1],
			Value: parsedTag[0][2],
		}
	}

	return sftm
}

func extractType(tags structFieldTagsMetadata) (aType string) {
	aType = "string"

	if tag, ok := tags.getTagByName(joinEntityTag); ok {
		aType = "q" + tag.Value
	}

	return aType
}

func extractValue(tags structFieldTagsMetadata) string {
	aValue := ""

	if tag, ok := tags.getTagByName(joinEntityTag); ok {
		aValue = "Q" + tag.Value
	} else if tag, ok := tags.getTagByName(entityTag); ok {
		aValue = `"` + tag.Value + `"`
	}

	return aValue
}

func (tags structFieldTagsMetadata) getTagByName(tagName string) (*structFieldTagMetadata, bool) {
	for _, tag := range tags {
		if tag.Name == tagName {
			return &tag, true
		}
	}

	return nil, false
}