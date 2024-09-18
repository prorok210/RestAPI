package main

import (
	"context"
	"fmt"
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
	users := &TableUsers{BaseTable{TableName: "users"}}

	if err != nil {
		fmt.Println(err)
	}
	// user := newUser("Василий2 Пупкинн", "vvas2yapypkin@gmail.com")
	// Create(user)
	// users.GetAll()
	user, err := users.getById(255)
	fmt.Println(user)
	// if err != nil {
	// 	panic(err)
	// }
	// user.Name = "Виталий Крючков2"
	// Update(user)

	// users.GetAll()
	// users.GetAll()
}
