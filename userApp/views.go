package userApp

import (
	"RestAPI/server"
	"strconv"
)

func HelloView(request server.HttpRequest) (server.HttpResponse) {
	if request.Method == "GET" {
		response := server.HTTP200
		response.Body = `{"Hello, world!"}`
		response.Headers["Content-Length"] = strconv.Itoa(len(response.Body))
		return response
	} else {
		response := server.HTTP405
		return response
	}
}

func GoodbyeView(request server.HttpRequest) (server.HttpResponse) {
	if request.Method == "GET" {
		response := server.HTTP200
		response.Body = `{"Goodbye, world!"}`
		response.Headers["Content-Length"] = strconv.Itoa(len(response.Body))
		return response
	} else {
		response := server.HTTP405
		return response
	}
}

func AddView(request server.HttpRequest) (server.HttpResponse) {
	if request.Method == "POST" {
		response := server.HTTP201
		response.Body = `{"Addition"}`
		response.Headers["Content-Length"] = strconv.Itoa(len(response.Body))
		return response
	} else {
		response := server.HTTP405
		return response
	}
}