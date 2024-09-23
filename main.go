package main

import (
	"RestAPI/app"
	"RestAPI/server"
	"log"
)

func main() {
	serv, er := server.CreateServer(app.MainApplication)
	if er != nil {
		log.Println("Error creating server", er)
		return
	}

	app.InitViews()

	er = serv.Start()
	if er != nil {
		log.Println("Error starting server", er)
		return
	}
	select {}
}
