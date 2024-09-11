package server

import (
	"log"
	"net"
	"strconv"
	"time"
)


func handleConnection(clientConn net.Conn, requestChan chan <- HttpRequest, responseChan <- chan []byte) {
	defer clientConn.Close()
	defer log.Println("Connection closed")
	log.Println("Connection accepted from: ", clientConn.RemoteAddr().String())

	timeoutDuration := time.Duration(CONN_TIMEOUT) * time.Second

	keepAlive := true
	for keepAlive {
		timer := time.NewTimer(timeoutDuration)
		defer timer.Stop()
		
		buffer := make([]byte, BUFSIZE)
		readChan := make(chan int)
		readErrChan := make(chan error)
		
		go func() {
			n, er := clientConn.Read(buffer)
			if er != nil {
				readErrChan <- er
			} else {
				readChan <- n
			}
		}()

		var bytesRead int
		select {
		case bytesRead := <-readChan:
			if bytesRead == 0 {
				log.Println("Closing connection")
				return
			}
			log.Println("Received: ", string(buffer[:bytesRead]))
		case er := <-readErrChan:
			log.Println("Error reading from connection: ", er)
			return
		case <-timer.C:
			log.Println("Connection timeout")
			return
		}

		requestSrct := new(HttpRequest)
		er := requestSrct.ParseRequest(buffer[:bytesRead])
		if er != nil {
			log.Println("Error parsing request")
			return
		}
		
		select {
		case requestChan <- *requestSrct:
		case <-timer.C:
			log.Println("Connection timeout")
			return
		}

		select {
		case response := <-responseChan:
			_, err := clientConn.Write(response)
			if err != nil {
				log.Println("Error writing to connection:", err)
				return
			}
		case <-timer.C:
			log.Println("Timeout waiting for response")
			return
		}
	
		keepAlive = false
		for key, value := range requestSrct.Headers {
			if key == "Connection" && value == "keep-alive" {
				keepAlive = true
			}
			if key == "Content-Length" {
				reqLen, er := strconv.Atoi(value)
				if er != nil {
					log.Println("Error parsing content length")
					continue
				}
				if len(requestSrct.Body) < reqLen {
					log.Println("Waiting for more data")
					keepAlive = true
					break
				}
			}
		}

		timer.Reset(timeoutDuration)
	}
}