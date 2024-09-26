package main

import (
	"RestAPI/app"
	"RestAPI/orm"
	"RestAPI/server"
	"log"
)

func main() {
	serv, er := server.CreateServer(app.MainApplication)
	if er != nil {
		log.Println("Error creating server", er)
		return
	}

	er = orm.InitDB()

	if er != nil {
		log.Println("Error inittializing DB", er)
	}

	app.InitHandlers()

	er = server.InitEnv()
	if er != nil {
		log.Println("Error initializing environment", er)
		return
	}

	er = serv.Start()
	if er != nil {
		log.Println("Error starting server", er)
		return
	}

	select {}
}
