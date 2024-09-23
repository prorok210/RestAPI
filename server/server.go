package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type RequestHandler func(request *HttpRequest) ([]byte, error)

//go:generate moq --out=mocks/Conn_moq_test.go . Conn
type Conn interface {
	net.Conn
}

type Server struct {
	servAddr  string
	listener  net.Listener
	handleApp RequestHandler
}

func CreateServer(mainApplication RequestHandler) (*Server, error) {
	regex := regexp.MustCompile(`^([a-zA-Z0-9.-]+):([0-9]{1,5})$`)
	addr := HOST + ":" + strconv.Itoa(PORT)
	if !regex.MatchString(addr) {
		log.Println("Invalid address format")
		return nil, errors.New("Invalid address format")
	}
	if mainApplication == nil {
		log.Println("Main application handler is not set")
		return nil, errors.New("Main application handler is not set")
	}

	server := &Server{
		servAddr:  addr,
		handleApp: mainApplication,
	}
	return server, nil
}

func (s *Server) Start() error {
	log.Println("Starting server ...")
	if s == nil {
		log.Println("Server is not created")
		return errors.New("Server is not created")
	}

	if s.servAddr == "" || s.handleApp == nil {
		log.Println("Server address or handle func is not set")
		return errors.New("Server address or handle func is not set")
	}

	listener, er := net.Listen("tcp", s.servAddr)
	if er != nil {
		log.Println("Error starting server", er)
		return er
	}
	s.listener = listener

	log.Println("Server started successfully on ", s.servAddr)

	go s.Listen()

	return nil
}

func (s *Server) Stop() {
	if s == nil {
		log.Println("Server is not created")
		return
	}

	if s.listener == nil {
		log.Println("Server is not started")
		return
	}

	log.Println("Stopping server ...")
	s.listener.Close()
}

func (s *Server) Listen() {
	log.Println("Listening for incoming connections")
	defer s.listener.Close()
	defer log.Println("Server stopped")

	for {
		clientConn, er := s.listener.Accept()
		if er != nil {
			log.Println("Error accepting connection", er)
			continue
		}

		if isAllowedHostMiddleware(clientConn.RemoteAddr().String()) {
			log.Println("Connection accepted from: ", clientConn.RemoteAddr().String())
			go s.ConnProcessing(clientConn)
		} else {
			log.Println("Connection refused from: ", clientConn.RemoteAddr().String())
			clientConn.Close()
		}
	}
}

func (s *Server) ConnProcessing(clientConn Conn) {
	defer clientConn.Close()
	defer log.Println("Connection closed with: ", clientConn.RemoteAddr().String())

	clientConn.SetDeadline(time.Now().Add(CONN_TIMEOUT * time.Second))

	buf := make([]byte, BUFSIZE)

	for {
		bytesRead, er := clientConn.Read(buf)
		if er != nil {
			if netErr, ok := er.(net.Error); ok && netErr.Timeout() {
				log.Println("Read timeout", netErr)
				clientConn.Write(HTTP408.ToBytes())
				return
			}
			if er.Error() == "EOF" {
				log.Println("Connection closed by client")
				return
			}
			log.Println("Error reading request", er)
			clientConn.Write(HTTP400.ToBytes())
			continue
		}
		fmt.Println("Request:", string(buf[:bytesRead]))

		// log.Println(clientConn.RemoteAddr().String(), strings.Split(string(buf[:bytesRead]), "\n")[0])

		request := new(HttpRequest)
		er = request.ParseRequest(buf[:bytesRead])
		if er != nil {
			log.Println("Error parsing request", er)
			clientConn.Write(HTTP400.ToBytes())
			continue
		}
		fmt.Println("Body: ", request.Body)

		er = reqMiddleware(request, clientConn)
		if er != nil {
			log.Println("Error in request middleware", er)
			return
		}

		response, er := s.handleApp(request)
		if er != nil {
			log.Println("Error handling request", er)
			clientConn.Write(HTTP500.ToBytes())
			continue
		}

		go func() {
			_, er = clientConn.Write(response)
			if er != nil {
				if netErr, ok := er.(net.Error); ok && netErr.Timeout() {
					log.Println("Write timeout", netErr)
					clientConn.Write(HTTP408.ToBytes())
					return
				}
				log.Println("Error writing response", er)
				return
			} else {
				log.Println(clientConn.RemoteAddr().String(), strings.Split(string(response), "\n")[0])
			}
		}()

		er = keepAliveMiddleware(request, clientConn)
		if er != nil {
			log.Println("Error in keep-alive middleware", er)
			return
		}
	}
}
