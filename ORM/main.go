package main

/* 1. Init your database in func InitDB in objFunctions file and append its in CreateTables slice
   2. Create new table Model in file Models.go and a new model with objects of this table
 	3. Add new record in tableRegistry in file constants.go
	4. add constructor of new model in file constructors.go
	5. All structure fields starting with a capital letter, all other letters are small
 	6. in the start of file objFunctions create a ToFields func */

import (
	"context"
	"log"
	"os"
)

func main() {
	InitDB()
	defer conn.Close(context.Background())

	// Откройте файл для записи логов
	logFile, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer logFile.Close()

	// log to file
	log.SetOutput(logFile)

	CreateTable("users", User{})

}
