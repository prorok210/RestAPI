package main

/*	1. Create new table Model in file Models.go and a new model with objects of this table
 	2. Add new record in tableRegistry and typeMap in file constants.go
	3. add constructor of new model in file constructors.go
	4. in the start of file objFunctions create a ToFields func
	5. Creating table for func "CreateTable(objOfTable{TableName: "tablename"})"
	6. All structure fields starting with a capital letter, all other letters are small
*/

import (
	"context"
	"fmt"
)

func main() {
	err := InitDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close(context.Background())

	// Get all users

	// fmt.Println(err)
	// users := TableUsers{BaseTable{"users"}}
	// user, err := users.getById(6)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(user)
	// messages := TableMessages{BaseTable{"messages"}}
	CreateTable(Message{TableName: "messagess"})

}
