package app

import "RestAPI/user"

/*
	Функция InitHandler() - инициализация списка представлений
	После создания представлений их необходимо зарегистрировать в этой функции, чтобы они были доступны для обработки запросов
	Для регистрации нужно передать url, по которому будет доступно представление, указатель на функцию-обработчик и имя предсталвения(оно должно совпадать с именем в документации для корректной работы)
	Роутер выдаст указатель на функцию, которая будет обрабатывать запрос или nil, если функции не нашлось
*/

func InitHandlers() {
	registerHandler("/images", user.ImageHandler, "images")
	registerHandler("/createUser", user.CreateUserHandler, "createUser")
	registerHandler("/verifyUser", user.VerifyUserHandler, "verifyUser")
	registerHandler("/createUserForm", user.CreateUserFormdataHandler, "createUserForm")
}
