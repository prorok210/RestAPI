package server

import (
	"errors"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type RequestHandler func(request *HttpRequest) ([]byte, error)

type Server struct {
	servAddr 		string
	listener 		net.Listener
	handleApp 		RequestHandler
}

func CreateServer(mainApplication RequestHandler) (*Server, error) {
	addr:= HOST + ":" + strconv.Itoa(PORT)
	server := &Server{
		servAddr: addr,
		handleApp: 	mainApplication,
	}
	return server, nil
}

func (s* Server) Start() error {
	log.Println("Starting server ...")
	listener, er := net.Listen("tcp", s.servAddr)
	if er != nil {
		log.Println("Error starting server", er)
		return er
	}
	s.listener = listener

	log.Println("Server started successfully on ", s.servAddr)

	s.Listen()

	return nil
}

func (s* Server) Listen(){
	log.Println("Listening for incoming connections")
	defer s.listener.Close()
	defer log.Println("Server stopped")

	for {
		clientConn, er := s.listener.Accept()
		if er != nil {
			log.Println("Error accepting connection", er)
			continue
		}

		if  isAllowedHostMiddleware(clientConn.RemoteAddr().String()) {
			log.Println("Connection accepted from: ", clientConn.RemoteAddr().String())
			go s.ConnProcessing(clientConn)
		} else {
			log.Println("Connection refused from: ", clientConn.RemoteAddr().String())
			clientConn.Close()
		}
	}
}


func (s* Server) ConnProcessing(clientConn net.Conn) {
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

		log.Println(clientConn.RemoteAddr().String(), strings.Split(string(buf[:bytesRead]), "\n")[0])

		request := new(HttpRequest)
		er = request.ParseRequest(buf[:bytesRead])
		if er != nil {
			log.Println("Error parsing request", er)
			clientConn.Write(HTTP400.ToBytes())
			continue
		}

		er = middleware(request, clientConn)
		if er != nil {
			log.Println("Error in middleware", er)
			return
		}

		response, er := s.handleApp(request)
		if er != nil {
			log.Println("Error handling request", er)
			clientConn.Write(HTTP500.ToBytes())
			continue
		}

		_, er = clientConn.Write(response)
		if er != nil {
			if netErr, ok := er.(net.Error); ok && netErr.Timeout() {
				log.Println("Write timeout", netErr)
				return
			}
			log.Println("Error writing response", er)
			continue
		} else {
			log.Println(clientConn.RemoteAddr().String(), strings.Split(string(response), "\n")[0])
		}
		
		er = keepAliveMiddleware(request, clientConn)
		if er != nil {
			log.Println("Error in keep-alive middleware", er)
			return
		}
	}
}


func isAllowedHostMiddleware(clientAddr string) bool {
	for _, allowedHost := range ALLOWED_HOSTS {
		if strings.Split(clientAddr, ":")[0] == allowedHost {
			return true
		}
	}
	return false
}

func middleware(request *HttpRequest, clientConn net.Conn ) (error) {

	methodFlag := false
	for _, allowedMethod := range ALLOWED_METHODS {
		if request.Method == allowedMethod {
			methodFlag = true
			break
		}
	}
	if !methodFlag {
		clientConn.Write(HTTP405.ToBytes())
		return errors.New("Method not allowed")
	}

	contentTypeFlag := false
	for _, supportedMediaType := range SUPPORTED_MEDIA_TYPES {
		if request.Body == "" {
			contentTypeFlag = true
			break
		}
		if request.Headers["Content-Type"] == supportedMediaType {
			contentTypeFlag = true
			break
		}
	}
	if !contentTypeFlag {
		clientConn.Write(HTTP415.ToBytes())
		return errors.New("Unsupported media type")
	}

	contentLengthFlag := false
	contentLengthStr := request.Headers["Content-Length"]
	if contentLengthStr != "" && len(request.Body) > 0 {
		contentLength, err := strconv.Atoi(contentLengthStr)
		if err != nil {
			log.Println("Invalid Content-Length:", contentLengthStr)
		} else {
			if contentLength == len([]byte(request.Body)) {
				contentLengthFlag = true
			}
		}
	}
	if contentLengthStr == "" && len(request.Body) == 0 {
		contentLengthFlag = true
	}
	if !contentLengthFlag {
		clientConn.Write(HTTP411.ToBytes())
		return errors.New("Content-Length required")
	}

	return nil
}

func keepAliveMiddleware(request *HttpRequest, clientConn net.Conn) (error) {
	for key := range request.Headers {
		if key == "Connection" {
			if request.Headers[key] == "close" {
				return errors.New("Connection: close")
			}
			if request.Headers[key] == "keep-alive" {
				clientConn.SetDeadline(time.Now().Add(CONN_TIMEOUT * time.Second))
			}
			if request.Headers[key] == "" {
				return errors.New("Connection: close")
			}
		}
	}
	return nil
}