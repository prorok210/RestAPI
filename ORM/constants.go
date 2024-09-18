package main

import (
	"os"

	"github.com/jackc/pgx/v5"
)

const CONNECTIONDATA string = "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable"

// Connection to the database
var (
	conn    *pgx.Conn
	logFile *os.File
)
