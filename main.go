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
	user, er := users.GetById(2)
	if er != nil {
		log.Println("Error getting user", er)
		return
	}
	fmt.Println(user)

	app.InitHandlers()

	er = serv.Start()
	if er != nil {
		log.Println("Error starting server", er)
		return
	}

	select {}
}
