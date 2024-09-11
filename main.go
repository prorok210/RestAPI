package main

import (
	"RestAPI/app"
	"RestAPI/server"
	"fmt"
)

func main() {
	listener, er := server.StartServer()
	if er != nil {
		fmt.Println(er)
		return
	}

	server.Listen(listener, app.MainApplication)
}