package app

/*
	Функция MainApplication() - основное приложение, в котором будет  производиться обработка запросов
*/

import (
	"RestAPI/server"
)

func MainApplication(request *server.HttpRequest) ([]byte, error) {
	view := router(request.Url)

	if view == nil {	
		return server.HTTP404.ToBytes(), nil
	}

	response := view(*request)

	return response.ToBytes(), nil
}
