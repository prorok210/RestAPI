package main

""" 1. Init your database in func InitDB in objFunctions file and append its in CreateTables slice
    2. Create new table Model in file Models.go and a new model with objects of this table
  	3. Add new record in tableRegistry in file constants.go
 	4. add constructor of new model in file constructors.go
 	5. All structure fields starting with a capital letter, all other letters are small
  	6. in the start of file objFunctions create a ToFields func """

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
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

	// Making a Go object from a sql table
	dialogs := &TableDialogs{BaseTable{TableName: "dialogs"}}
	// dialog := newDialog("Mee", "MyFriends")
	// Create(dialog)
	// dialogs.GetAll()
	users := &TableUsers{BaseTable{TableName: "users"}}
	// user := newUser("Vasya", "Piupiu")
	// Create(user)
	// dialogs.GetAll()
	dialog1, err := dialogs.getById(1)

	user, err := users.getById(1)
	fmt.Println(reflect.TypeOf(user))

	fmt.Println(reflect.TypeOf(dialog1))

	if err != nil {
		log.Printf("Error getting dialog by id: %v", err)
	} else {
		log.Printf("Dialog by id: %v", dialog1)
	}
	// users.GetAll()

}
