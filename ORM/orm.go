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
	user := newUser("Василий Пупкинн", "vvasyapypkin@gmail.com")
	Create(user)
	users.GetAll()
}
