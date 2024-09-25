package main

import (
	"reflect"

	"github.com/jackc/pgx/v5"
)

const CONNECTIONDATA string = "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable"

// Connection to the database
var (
	conn *pgx.Conn
	// mapp for create child objects
	tableRegistry = map[string]reflect.Type{
		"users":    reflect.TypeOf(User{TableName: "users"}),
		"dialogs":  reflect.TypeOf(Dialog{TableName: "dialogs"}),
		"messages": reflect.TypeOf(Message{TableName: "messages"}),
	}

	typeMap = map[string]reflect.Type{
		"users":   reflect.TypeOf((*User)(nil)).Elem(),
		"dialogs": reflect.TypeOf((*Dialog)(nil)).Elem(),
	}
	tagToSqlType = map[string]string{
		"varchar":  "character varying",
		"serial":   "integer",
		"int":      "integer",
		"uint":     "integer",
		"bool":     "boolean",
		"not null": "NO",
	}
)
