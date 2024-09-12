package server

import (
	"context"
	"log"
	"net"
	"strings"
	"time"
)


func handleConnection(clientConn net.Conn, requestChan chan <- HttpRequest, responseChan <- chan[]byte) {
	defer clientConn.Close()
	defer log.Println("Connection closed")

	ctx, cancel := context.WithTimeout(context.Background(), CONN_TIMEOUT * time.Second)
	defer cancel()

	log.Println("Connection accepted from: ", clientConn.RemoteAddr().String())


	for {
		buffer := make([]byte, BUFSIZE)

		bytesRead, er := clientConn.Read(buffer)
		if er != nil {
			log.Println("Error reading")
			return
		}
		if bytesRead == 0 {
			log.Println("Connection closed by client")
			return
		}

		line:= strings.Split(string(buffer[:bytesRead]), "\n")[0]
		log.Println(line)

		requestSrct := new(HttpRequest)
		er = requestSrct.ParseRequest(buffer[:bytesRead])
		if er != nil {
			log.Println("Error parsing request")
			return
		}

		select {
		case requestChan <- *requestSrct:
			log.Println("Request sent to channel")
			select {
			case response := <-responseChan:
				_, er = clientConn.Write(response)
				if er != nil {
					log.Println("Error writing response")
					return
				}
			log.Println("Response sent to client")
			}
			case <-ctx.Done():
				log.Println("Response timeout")
				return
		
		case <-ctx.Done():
			log.Println("Connection timeout")
			return
		}
		
	
		// keepAlive = false
		// for key, value := range requestSrct.Headers {
		// 	if key == "Connection" && value == "keep-alive" {
		// 		keepAlive = true
		// 	}
		// 	if key == "Content-Length" {
		// 		reqLen, er := strconv.Atoi(value)
		// 		if er != nil {
		// 			log.Println("Error parsing content length")
		// 			continue
		// 		}
		// 		if len(requestSrct.Body) < reqLen {
		// 			log.Println("Waiting for more data")
		// 			keepAlive = true
		// 			break
		// 		}
		// 	}
		// }

	}
}