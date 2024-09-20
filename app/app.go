package app

/*
	Функция MainApplication() - основное приложение, в котором будет  производиться обработка запросов
*/

import (
	"RestAPI/server"
	"strconv"
)

func MainApplication(request *server.HttpRequest) ([]byte, error) {
	view := router(request.Url)

	if view == nil {	
		return server.HTTP404.ToBytes(), nil
	}

	response := view(*request)
	if response.Body != "" {
		response.SetHeader("Content-Length", strconv.Itoa(len(response.Body)))
	}

	return response.ToBytes(), nil
}
