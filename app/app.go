package app

import (
	"RestAPI/server"
	"log"
	"time"
)

func MainApplication(requestChan <-chan server.HttpRequest, responseChan chan<- []byte) {
	for {
		select {
		case request, ok := <-requestChan:
			if !ok {
				log.Println("Request channel closed, exiting MainApplication")
				return
			}

			response := server.HttpResponse{
				Version: request.Version,
				Status:  200,
				Reason:  "OK",
				Headers: make(map[string]string),
				Body:    "Hello, World!",
			}
			response.Headers["Content-Type"] = "text/plain"
			
			select {
			case responseChan <- response.ToBytes():
				log.Println("Response sent successfully")
			case <-time.After(time.Duration(server.RESP_TIMEOUT) * time.Second):
				log.Println("Timeout sending response")
			}

		case <-time.After(time.Duration(server.RESP_TIMEOUT) * time.Second):
			continue
		}
	}
}
