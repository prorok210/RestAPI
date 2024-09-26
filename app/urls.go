package app

import (
	"RestAPI/core"
	"RestAPI/user"
)

type HandlerFunc func(core.HttpRequest) core.HttpResponse

var viewsList = make(map[string]HandlerFunc)

func registerHandler(url string, f HandlerFunc) {
	viewsList[url] = f
}

func router(url string) HandlerFunc {
	return viewsList[url]
}

/*
	Функция InitHandler() - инициализация списка представлений
	После создания представлений их необходимо зарегистрировать в этой функции, чтобы они были доступны для обработки запросов
	Роутер выдаст указатель на функцию, которая будет обрабатывать запрос или nil, если функции не нашлось
*/

func InitHandlers() {
	registerHandler("/images", user.ImageHandler)
	registerHandler("/createUser", user.CreateUserHandler)
	registerHandler("/verifyUser", user.VerifyUserHandler)
	registerHandler("/createUserForm", user.CreateUserFormdataHandler)
}
