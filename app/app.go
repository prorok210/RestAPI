package app

import (
	"RestAPI/server"
	"context"
	"log"
	"time"
)

func MainApplication(requestChan  <- chan server.HttpRequest, response chan <- []byte) {
	
	for {
		ctx, cancel := context.WithTimeout(context.Background(), server.CONN_TIMEOUT * time.Second)
		defer cancel()

		select {
		case request := <-requestChan:
			log.Println(request)

			resp := server.HttpResponse{
				Version: "HTTP/1.1",
				Status: 200,
				Reason: "OK",
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
				Body: "Hello, World!",
			}

			select {
			case response <- resp.ToBytes():
				log.Println("Response sent to channel")
			case <-ctx.Done():
				log.Println("Response timeout")
				return
			}
		default:
			log.Println("No request received")
		}
	}
}
