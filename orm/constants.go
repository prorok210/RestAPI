package orm

import (
	"RestAPI/core"

	"github.com/jackc/pgx/v5"
)

// Connection to the database
var (
	conData = core.CONNECTIONDATA
	conn    *pgx.Conn

	tagToSqlType = map[string]string{
		"varchar":            "character varying",
		"serial_primary_key": "integer",
		"ref":                "integer",
		"int":                "integer",
		"uint":               "integer",
		"bool":               "boolean",
		"not null":           "NO",
		"timestamp":          "timestamp without time zone",
	}
)
