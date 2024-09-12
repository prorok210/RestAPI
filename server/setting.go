package server

import (
	"time"
)

var ALLOWED_HOSTS = []string{
	"localhost",
	"127.0.0.1",
	"::1",
}

const (
	HOST string = "localhost"
	PORT int = 8080
	CONN_TIMEOUT time.Duration = 10
	RESP_TIMEOUT time.Duration = 10
	BUFSIZE int = 1024
	CHANNEL_BUFSIZE int = 256
)

