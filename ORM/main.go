package main

/*	1. Create new table Model in file Models.go and a new model with objects of this table
 	2. Add new record in tableRegistry in file constants.go
	3. add constructor of new model in file constructors.go
	4. All structure fields starting with a capital letter, all other letters are small
 	5. in the start of file objFunctions create a ToFields func
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

	err = CreateTable(User{TableName: "usrees"})
	user1 := User{TableName: "usrees", Name: "Alice", Email: "alice@example.com"}
	user2 := User{TableName: "usrees", Name: "Bob", Email: "bob@example.com"}
	user3 := User{TableName: "usrees", Name: "Charlie", Email: "charlie@example.com"}
	user4 := User{TableName: "usrees", Name: "David", Email: "david@example.com"}
	user5 := User{TableName: "usrees", Name: "Eve", Email: "eve@example.com"}
	user6 := User{TableName: "usrees", Name: "Frank", Email: "frank@example.com"}
	user7 := User{TableName: "usrees", Name: "Grace", Email: "grace@example.com"}
	user8 := User{TableName: "usrees", Name: "Heidi", Email: "heidi@example.com"}
	user9 := User{TableName: "usrees", Name: "Ivan", Email: "ivan@example.com"}
	user10 := User{TableName: "usrees", Name: "Judy", Email: "judy@example.com"}
	Create(user1)
	Create(user2)
	Create(user3)
	Create(user4)
	Create(user5)
	Create(user6)
	Create(user7)
	Create(user8)
	Create(user9)
	Create(user10)

	// Get all users

	fmt.Println(err)
	users := TableUsers{BaseTable{"users"}}
	user, err := users.getById(6)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(user)

}
