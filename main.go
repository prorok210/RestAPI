package main

import (
	"RestAPI/app"
	"RestAPI/core"
	"RestAPI/db"
	"RestAPI/orm"
	"fmt"
	"log"
)

func main() {
	serv, er := core.CreateServer(app.MainApplication)
	if er != nil {
		log.Println("Error creating server", er)
		return
	}
	er = core.InitEnv()
	if er != nil {
		log.Println("Error initializing environment", er)
		return
	}
	er = db.Register()
	if er != nil {
		log.Println("Error registering models", er)
		return
	}
	er = orm.InitDB()
	if er != nil {
		log.Println("Error inittializing DB", er)
		return
	}

	users := orm.BaseTable{TableName: "users"}
	users.GetAll()
	userInterface, er := users.GetById(2)
	if er != nil {
		log.Println("Error getting user", er)
		return
	}
	fmt.Println("USER", userInterface)
	user, ok := userInterface.(db.User)
	if !ok {
		// Обработка ошибки, если значение не является db.User
		panic("не удалось привести interface{} к db.User")
	}
	// user.Interface().(*db.User).Name = "Vasya"
	er = orm.Update(user)
	if er != nil {
		log.Println("Error updating user", er)
		return
	}
	fmt.Println("USER", user)

	app.InitHandlers()

	er = serv.Start()
	if er != nil {
		log.Println("Error starting server", er)
		return
	}

	select {}
}
