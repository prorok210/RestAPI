package app

import (
	"RestAPI/server"
	"RestAPI/user"
)

type viewFunc func(server.HttpRequest) (server.HttpResponse)


var viewsList = make(map[string]viewFunc)


func registerView(url string, f viewFunc) {
	viewsList[url] = f
}

func router(url string) (viewFunc) {
	return viewsList[url]
}

/*
	Функция InitViews() - инициализация списка представлений.
	После создания представлений их необходимо зарегистрировать в этой функции, чтобы они были доступны для обработки запросов
	Роутер выдаст указатель на функцию, которая будет обрабатывать запрос или nil, если такой функции нет
*/

func InitViews() {
	registerView("/hello", user.HelloView)
	registerView("/goodbye", user.GoodbyeView)
	registerView("/add", user.AddView)
}
