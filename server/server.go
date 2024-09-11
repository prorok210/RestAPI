package server

import (
	"errors"
	"log"
	"net"
)


type RequestHandler func(requestChan <-chan HttpRequest, responseChan chan<- []byte)


func StartServer() (*net.TCPListener, error) {
	log.Println("Starting server...")
	addrSet := &net.TCPAddr {
		IP: net.ParseIP(HOST),
		Port: PORT,
	}

	listener, er := net.ListenTCP("tcp", addrSet)
	if er != nil {
		return nil, errors.New("Error starting server")
	}
	log.Println("Server successfully started at: ", HOST, ":", PORT)
	return listener, nil
}


func Listen(listener *net.TCPListener, mainApplication RequestHandler){
	log.Println("Listening for incoming connections...")
	defer listener.Close()

	for {
		clientConn, er := listener.Accept()
		if er != nil {
			log.Println("Error accepting connection")
			continue
		}
		clientAddr, _, er := net.SplitHostPort(clientConn.RemoteAddr().String())
		if er != nil {
			log.Println("Error splitting host and port")
			clientConn.Close()
			continue
		}

		if isAllowedHost(clientAddr) {
			requestChan := make(chan HttpRequest, CHANNEL_BUFSIZE)
			responseChan := make(chan []byte, CHANNEL_BUFSIZE)
			go handleConnection(clientConn, requestChan, responseChan)
			go mainApplication((<-chan HttpRequest)(requestChan), (chan<- []byte)(responseChan))
		} else {
			log.Printf("Host: %s is not allowed\n", clientConn.RemoteAddr().String())
			clientConn.Close()
		}
	}
}


func isAllowedHost(clientAddr string) bool {
	for _, allowedHost := range ALLOWED_HOSTS {
		if clientAddr == allowedHost {
			return true
		}
	}
	return false
}