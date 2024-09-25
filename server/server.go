package server

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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
	httpAddr  		string
	httpsAddr 		string
	httpListener  	net.Listener
	httpsListener	net.Listener
	handleApp 		RequestHandler
	certFile  		string
	keyFile   		string
}

func CreateServer(mainApplication RequestHandler) (*Server, error) {
	regex := regexp.MustCompile(`^([a-zA-Z0-9.-]+):([0-9]{1,5})$`)
	httpAddr := HOST + ":" + strconv.Itoa(HTTP_PORT)
	httpsAddr := HOST + ":" + strconv.Itoa(HTTPS_PORT)
	if !regex.MatchString(httpAddr) || !regex.MatchString(httpsAddr) {
		log.Println("Invalid address format")
		return nil, errors.New("Invalid address format")
	}
	if mainApplication == nil {
		log.Println("Main application handler is not set")
		return nil, errors.New("Main application handler is not set")
	}

	server := &Server{
		httpAddr:  httpAddr,
		httpsAddr: httpsAddr,
		handleApp: mainApplication,
		certFile: CERT_FILE,
		keyFile:  KEY_FILE,
	}
	return server, nil
}

func (s *Server) Start() error {
	if _, err := os.Stat(s.certFile); os.IsNotExist(err) {
        return fmt.Errorf("certificate file not found: %s", s.certFile)
    }
    if _, err := os.Stat(s.keyFile); os.IsNotExist(err) {
        return fmt.Errorf("key file not found: %s", s.keyFile)
    }

	log.Println("Starting server ...")
	if s == nil {
		log.Println("Server is not created")
		return errors.New("Server is not created")
	}

	if s.httpAddr == "" || s.handleApp == nil {
		log.Println("Server address or handle func is not set")
		return errors.New("Server address or handle func is not set")
	}

	listener, er := net.Listen("tcp", s.httpAddr)
	if er != nil {
		log.Println("Error starting server", er)
		return er
	}
	s.httpListener = listener
	go s.ListenHTTP()

	log.Println("Http server started successfully on ", s.httpAddr)

	cert, err := tls.LoadX509KeyPair(s.certFile, s.keyFile)
	if err != nil {
		log.Printf("Error loading SSL certificates: %v", err)
		return err
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	httpsListener, err := tls.Listen("tcp", s.httpsAddr, config)
	if err != nil {
		log.Printf("Error starting HTTPS server: %v", err)
		return err
	}
	s.httpsListener = httpsListener
	go s.ListenHTTPS()

	log.Println("Https server started successfully on ", s.httpsAddr)

	return nil
}

func (s *Server) ListenHTTP() {
	defer s.httpListener.Close()
	defer log.Println("Http server stopped")

	for {
		clientConn, er := s.httpListener.Accept()
		if er != nil {
			log.Println("Error accepting connection", er)
			continue
		}

		go func (clientConn Conn) {
			defer clientConn.Close()
			defer log.Println("Connection closed with: ", clientConn.RemoteAddr().String())

			reader := bufio.NewReader(clientConn)
			reqLine, _, err := reader.ReadLine()
			if err != nil {
				log.Println("Error reading HTTP request: ", err)
				return
			}

			parts := strings.Split(string(reqLine), " ")
			if len(parts) < 2 {
				log.Println("Invalid HTTP request")
				return
			}

			path := parts[1]
			httpsHost := strings.Split(s.httpsAddr, ":")[0]
			httpsURL := fmt.Sprintf("https://%s%s", httpsHost, path)

			response := fmt.Sprintf("HTTP/1.1 301 Moved Permanently\r\nLocation: %s\r\nConnection: close\r\n\r\n", httpsURL)
			clientConn.Write([]byte(response))
		}(clientConn)
	}
}

func (s *Server) ListenHTTPS() {
	defer s.httpsListener.Close()
	defer log.Println("Https server stopped")

	for {
		clientConn, err := s.httpsListener.Accept()
		if err != nil {
			log.Println("Error accepting HTTPS connection", err)
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

	tlsConn, ok := clientConn.(*tls.Conn)
	if !ok {
		log.Println("Error type assertion")
		return
	}

	er := tlsConn.Handshake()
	if er != nil {
		log.Println("Error TLS handshake", er)
		return
	}

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

func (s *Server) Stop() {
	if s == nil {
		log.Println("Server is not created")
		return
	}

	if s.httpListener == nil || s.httpsListener == nil {
		log.Println("Server is not started")
		return
	}

	log.Println("Stopping server ...")
	s.httpListener.Close()
	s.httpsListener.Close()
}
