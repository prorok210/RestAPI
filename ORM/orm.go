package main

import "context"

func main() {
	// Таблицу в объект языка
	InitDB()
	defer conn.Close(context.Background())
	users := &TableUsers{BaseModel{TableName: "users"}}

	// когда появится метод Create, тогда не нужно будет передавать conn
	// Либо сделать conn отдельным полем BaseModel
	users.GetAll()
	testUser := User{
		Name:  "Иннокентий1 Лельков",
		Email: "keshalelkov1@gmail.com",
	}
	users.Create(&testUser)
	users.GetAll()
}
