package core

/*
	Настройки сервера (и приложения)
*/

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var APPS = []string{
	"user",
	"chats",
}

const (
	// Настройки сервера
	HOST          string        = "localhost"
	HTTP_PORT     int           = 8082
	HTTPS_PORT    int           = 8446
	CERT_FILE     string        = "/home/user/etc/ssl/certs/dev.crt"
	KEY_FILE      string        = "/home/user/etc/ssl/private/dev.key"
	CONN_TIMEOUT  time.Duration = 20
	WRITE_TIMEOUT time.Duration = 20
	BUFSIZE       int           = 5 * 1024 * 1024
	IMAGES_DIR    string        = "/media/images"
	// Настройки мидлваров
	IS_ALLOWED_HOSTS bool = true
	REQ_MIDDLEWARE   bool = true
	KEEP_ALIVE       bool = true
	// Настройки БД
	DB_USER string = "admin"
	DB_NAME string = "dev"
	DB_PORT int    = 5432
)

// Настройки БД
var (
	DB_PASSWORD    string
	CONNECTIONDATA string
)

var ALLOWED_HOSTS = []string{
	"localhost",
	"127.0.0.1",
	"77.232.37.23",
	"::1",
}

var ALLOWED_METHODS = []string{
	"OPTIONS",
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

/*
Настройки JWT
*/
var (
	JWT_ACCESS_SECRET_KEY       string
	JWT_REFRESH_SECRET_KEY      string
	JWT_ACCESS_EXPIRATION_TIME  time.Duration = time.Hour * 24
	JWT_REFRESH_EXPIRATION_TIME time.Duration = time.Hour * 336
)

/*
MTS API KEY
*/
var (
	MTS_API_KEY    string
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
	if MTS_API_KEY == "" || MTS_API_NUMBER == "" {
		log.Fatalf("Error env load %v", err)
		return err
	}
	JWT_ACCESS_SECRET_KEY = os.Getenv("JWT_ACCESS_SECRET_KEY")
	JWT_REFRESH_SECRET_KEY = os.Getenv("JWT_REFRESH_SECRET_KEY")
	if JWT_ACCESS_SECRET_KEY == "" || JWT_REFRESH_SECRET_KEY == "" {
		log.Fatalf("Error env load %v", err)
		return err
	}
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	CONNECTIONDATA = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)
	if CONNECTIONDATA == "" || DB_PASSWORD == "" {
		log.Fatalf("Error env load %v", err)
		return err
	}
	return nil
}
