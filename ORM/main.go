package main

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
	users := &TableUsers{BaseTable{TableName: "users"}}

	if err != nil {
		fmt.Println(err)
	}

	cell, err := users.getById(36)
	user := cell.(*User)
	fmt.Println(user)
	user.Name = "Василий Пупкин"
	fmt.Println(reflect.TypeOf(user))

	fmt.Println(user)
	Update(user)
	cell, err = users.getById(36)
	user = cell.(*User)
	fmt.Println(user)
	// users.GetAll()

}
