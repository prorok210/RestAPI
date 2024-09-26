package app

/*
	Функция MainApplication() - основное приложение, в котором будет  производиться обработка запросов
*/

import (
	"RestAPI/core"
	"fmt"
	"mime"
	"strconv"
	"strings"
)

func MainApplication(request *core.HttpRequest) ([]byte, error) {
	if request == nil {
		return core.HTTP400.ToBytes(), nil
	}
	if request.Method == "OPTIONS" {
		response := core.HTTP200
		allowedOrigins := strings.Join(core.ALLOWED_HOSTS, ", ")
		allowedMethods := strings.Join(core.ALLOWED_METHODS, ", ")
		allowedContentTypes := strings.Join(core.SUPPORTED_MEDIA_TYPES, ", ")
		response.SetHeader("Access-Control-Allow-Origin", allowedOrigins)
		response.SetHeader("Access-Control-Allow-Methods", allowedMethods)
		response.SetHeader("Access-Control-Allow-Headers", "*")
		response.SetHeader("Access-Control-Allow-Content-Type", allowedContentTypes)
		response.SetHeader("Access-Control-Allow-Credentials", "true")
		return response.ToBytes(), nil
	}

	contentType, _, _ := mime.ParseMediaType(request.Headers["Content-Type"])
	if contentType == "application/x-www-form-urlencoded" || contentType == "multipart/form-data" {
		er := request.ParseFormData()
		if er != nil {
			fmt.Println("Error parsing form data:", er)
			return core.HTTP400.ToBytes(), nil
		}
	}

	view := router(request.Url)
	if view == nil {
		return core.HTTP404.ToBytes(), nil
	}

	response := view(*request)
	if response.Body != "" {
		response.SetHeader("Content-Length", strconv.Itoa(len(response.Body)))
	}

	return response.ToBytes(), nil
}
