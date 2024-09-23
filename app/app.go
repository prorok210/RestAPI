package app

/*
	Функция MainApplication() - основное приложение, в котором будет  производиться обработка запросов
*/

import (
	"RestAPI/server"
	"fmt"
	"strconv"
)

func MainApplication(request *server.HttpRequest) ([]byte, error) {
	if request == nil {
		return server.HTTP400.ToBytes(), nil
	}
	if request.Headers["Content-Type"] == "application/x-www-form-urlencoded" || request.Headers["Content-Type"] == "multipart/form-data" {
		er := request.ParseFormData()
		if er != nil {
			fmt.Println("Error parsing form data:", er)
			return server.HTTP400.ToBytes(), nil
		}
	}

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
