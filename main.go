package main

import (
	"RestAPI/app"
	"RestAPI/server"
	"RestAPI/user"
	"fmt"
	"log"
)

func main() {
	serv, er := server.CreateServer(app.MainApplication)
	if er != nil {
		fmt.Println("Error creating server", er)
		return
	}
	
	app.InitViews()

	er = server.InitEnv()
	if er != nil {
		log.Println("Error initializing environment", er)
		return
	}
	fmt.Print(user.GenerateSecretKey(64))
	er = serv.Start()
	if er != nil {
		log.Println("Error starting server", er)
		return
	}
	
	select {}
}
