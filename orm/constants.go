package orm

import (
	"github.com/jackc/pgx/v5"
)

// Подключение к базе данных
var (
	conn *pgx.Conn

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
