package server

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"
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
	Mutex     *sync.Mutex 
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
		Mutex:    &sync.Mutex{},
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

	for {
		clientConn.SetDeadline(time.Now().Add(CONN_TIMEOUT * time.Second))

		receivedData, er := func(clientConn Conn) ([]byte, error) {
			reader := bufio.NewReader(clientConn)

			startLine, err := reader.ReadString('\n')
			if err != nil {
				return nil, err
			}
			startLine = strings.TrimSpace(startLine)
		
			headers := make([]byte, 0, 4096)
			for {
				line, err := reader.ReadBytes('\n')
				headers = append(headers, line...)
				if err != nil || len(bytes.TrimSpace(line)) == 0 {
					break
				}
			}
		
			contentLength := 0
			headerLines := bytes.Split(headers, []byte("\r\n"))
			for _, line := range headerLines {
				if bytes.HasPrefix(bytes.ToLower(line), []byte("content-length:")) {
					parts := bytes.SplitN(line, []byte(":"), 2)
					if len(parts) == 2 {
						contentLength, _ = strconv.Atoi(string(bytes.TrimSpace(parts[1])))
						break
					}
				}
			}

			body := make([]byte, contentLength)
			_, err = io.ReadFull(reader, body)
			if err != nil {
				return nil, err
			}
		
			fullRequest := append([]byte(startLine+"\r\n"), headers...)
			fullRequest = append(fullRequest, body...)
		
			return fullRequest, nil
		}(clientConn)

		if er != nil {
			if netErr, ok := er.(net.Error); ok && netErr.Timeout() {
				log.Println("Read timeout", netErr)
				clientConn.Write(HTTP408.ToBytes())
			} else if er == io.EOF {
				log.Println("Connection closed by client")
			} else {
				log.Println("Error reading request", er)
				clientConn.Write(HTTP400.ToBytes())
			}
			return
		}

		log.Println(clientConn.RemoteAddr().String(), strings.Split(string(receivedData), "\n")[0])

		request := &HttpRequest{}
		err := request.ParseRequest(receivedData)
		if err != nil {
			log.Println("Error parsing request", err)
			clientConn.Write(HTTP400.ToBytes())
			return
		}

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

		clientConn.SetWriteDeadline(time.Now().Add(WRITE_TIMEOUT * time.Second))
		_, err = clientConn.Write(response)
		if err != nil {
			log.Println("Error writing response", err)
			return
		}
		log.Println(clientConn.RemoteAddr().String(), strings.Split(string(response), "\n")[0])

		er = keepAliveMiddleware(request, clientConn)
		if er != nil {
			log.Println("Error in keep-alive middleware", er)
			return
		}
	}
}