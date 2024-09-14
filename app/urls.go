package app

import (
	"RestAPI/server"
	"RestAPI/userApp"
)

type viewFunc func(server.HttpRequest) (server.HttpResponse)


var viewsList = make(map[string]viewFunc)


func registerView(url string, f viewFunc) {
	viewsList[url] = f
}

func router(url string) (viewFunc) {
	return viewsList[url]
}

func InitViews() {
	registerView("/hello", userApp.HelloView)
	registerView("/goodbye", userApp.GoodbyeView)
	registerView("/add", userApp.AddView)
}
