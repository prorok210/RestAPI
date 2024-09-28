package main

import (
	"RestAPI/app"
	"RestAPI/core"
	"RestAPI/docs"
	"RestAPI/orm"
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

	er = orm.InitDB()
	if er != nil {
		log.Println("Error inittializing DB", er)
		return
	}

	app.InitHandlers()

	er = docs.GenerateDocs()
	if er != nil {
		log.Println("Error generating docs", er)
		return
	}

	er = serv.Start()
	if er != nil {
		log.Println("Error starting server", er)
		return
	}

	select {}
}
