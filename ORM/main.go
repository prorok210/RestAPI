package main

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

	// Making a Go object from a sql table
	dialogs := &TableDialogs{BaseTable{TableName: "dialogs"}}
	dialog := newDialog("Mee", "MyFriends")
	Create(dialog)
	dialogs.GetAll()
	// users := &TableUsers{BaseTable{TableName: "users"}}
	user := newUser("Vasya", "Piupiu")
	Create(user)
	dialogs.GetAll()
	// users.GetAll()

}
