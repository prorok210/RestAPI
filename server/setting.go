package server

/*
	Настройки сервера
*/

import (
	"time"
)

var ALLOWED_HOSTS = []string{
	"localhost",
	"127.0.0.1",
	"::1",
}

var ALLOWED_METHODS = []string{
	"GET",
	"POST",
	"PUT",
	"PATCH",
	"DELETE",
}

var SUPPORTED_MEDIA_TYPES = []string{
	"application/json",
	"application/x-www-form-urlencoded",
	"multipart/form-data",
	"text/plain",
	"image/jpeg",
	"image/png",
}

const (
	HOST string = "localhost"
	PORT int = 8081
	CONN_TIMEOUT time.Duration = 20
	BUFSIZE int = 2048
)

