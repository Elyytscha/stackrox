package main

import (
	"fmt"
	"strings"

	. "github.com/dave/jennifer/jen"
)

func dataTypeToSQLType(dataType DataType) string {
	var sqlType string
	switch dataType {
	case BOOL:
		sqlType = "bool"
	case NUMERIC:
		sqlType = "numeric"
	case STRING:
		sqlType = "varchar"
	case DATETIME:
		sqlType = "timestamp"
	case MAP, STRING_ARRAY:
		sqlType = "jsonb"
	case ENUM:
		sqlType = "integer"
	case JSONB:
		sqlType = "jsonb"
	default:
		panic(dataType.String())
	}
	return sqlType
}

func fieldsFromPath(b *strings.Builder, table *Path) {
	for i, elem := range table.Elems {
		if !(table.Parent == nil && i == 0) {
			fmt.Fprint(b, ", ")
		}
		fmt.Fprintf(b, "%s %s", elem.SQLPath(), dataTypeToSQLType(elem.DataType))
	}
	for _, child := range table.Children {
		fieldsFromPath(b, child)
	}
}

func generateTableCreation(f *File, tableName string, table *Path) {
	var b strings.Builder
	fmt.Fprintf(&b, "create table if not exists %s (", tableName)
	fieldsFromPath(&b, table)
	fmt.Fprintf(&b, ")")

	f.Const().Id("createTableQuery").Op("=").Lit(b.String())
	f.Func().Id("CreateTable").Params(Id("db").Op("*").Qual("database/sql", "DB")).Error().Block(
		List(Id("_"), Err()).Op(":=").Id("db").Dot("Exec").Call(Id("createTableQuery")),
		Return(Err()),
	)
}
