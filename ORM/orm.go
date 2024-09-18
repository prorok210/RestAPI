package main

import (
	"context"
)

func main() {
	// Таблицу в объект языка
	InitDB()
	defer conn.Close(context.Background())
	users := &TableUsers{BaseTable{TableName: "users"}}

	// когда появится метод Create, тогда не нужно будет передавать conn
	// Либо сделать conn отдельным полем BaseModel
	users.GetAll()
	// user := newUser("Василий2 Пупкинн", "vvas2yapypkin@gmail.com")
	// Create(user)
	// users.GetAll()
	// user, err := users.getById(25)
	// if err != nil {
	// 	panic(err)
	// }
	// user.Name = "Виталий Крючков2"
	// Update(user)

	// users.GetAll()
	// users.GetAll()
}
