package app

/*
	Функция MainApplication() - основное приложение, в котором будет  производиться обработка запросов
*/

import (
	"RestAPI/server"
	"strconv"
)

func MainApplication(request *server.HttpRequest) ([]byte, error) {
		response := server.HttpResponse{
			Version: "HTTP/1.1",
			Status: 200,
			Reason: "OK",
			Headers: make(map[string]string),
			Body: "Hello, World!",
		}

		response.Headers["Content-Type"] = "text/plain"
		response.Headers["Content-Length"] = strconv.Itoa(len(response.Body))

		return response.ToBytes(), nil
}
