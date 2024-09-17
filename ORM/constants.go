package main

import "github.com/jackc/pgx/v5"

const CONNECTIONDATA string = "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable"

var (
	conn        *pgx.Conn
	InitDBError error
)
