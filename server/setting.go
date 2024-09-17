package server

/*
	Настройки сервера (и приложения)
*/

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var ALLOWED_HOSTS = []string{
	"localhost",
	"127.0.0.1",
	// "::1",
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
	HOST string 				= "localhost"
	PORT int 					= 8081
	CONN_TIMEOUT time.Duration	= 20
	BUFSIZE int 				= 2048	
)
/*
	MTS API KEY
*/
var (
	MTS_API_KEY string
	MTS_API_NUMBER string
)


func InitEnv() error {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error env load %v", err)
		return err
    }
	MTS_API_KEY = os.Getenv("MTS_API_KEY")
	MTS_API_NUMBER = os.Getenv("MTS_API_NUMBER")
	if (MTS_API_KEY == "" || MTS_API_NUMBER == "") {
		log.Fatalf("Error env load %v", err)
		return err
	}
	return nil
}