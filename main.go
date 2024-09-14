package main

import (
	"RestAPI/app"
	"RestAPI/server"
	"fmt"
	"log"
)

func main() {
	serv, er := server.CreateServer(app.MainApplication)
	if er != nil {
		fmt.Println("Error creating server", er)
		return
	}

	er = serv.Start()
	if er != nil {
		log.Println("Error starting server", er)
		return
	}
	select {}
}	