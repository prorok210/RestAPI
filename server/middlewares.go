package server

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

func isAllowedHostMiddleware(clientAddr string) bool {
	for _, allowedHost := range ALLOWED_HOSTS {
		if strings.Split(clientAddr, ":")[0] == allowedHost {
			return true
		}
	}
	return false
}

func reqMiddleware(request *HttpRequest, clientConn Conn) (error) {

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

	contentLengthStr, hasContentLength := request.Headers["Content-Length"]

	if hasContentLength {
		contentLength, err := strconv.Atoi(contentLengthStr)
		if err != nil {
			log.Println("Invalid Content-Length:", contentLengthStr)
			clientConn.Write(HTTP411.ToBytes())
			return errors.New("invalid Content-Length header")
		}
	
		if contentLength != len(request.Body) {
			clientConn.Write(HTTP411.ToBytes())
			return errors.New("Content-Length does not match body length")
		}
	
	} else if len(request.Body) > 0 {
		clientConn.Write(HTTP411.ToBytes())
		return errors.New("Content-Length required")
	}

	return nil
}

func keepAliveMiddleware(request *HttpRequest, clientConn Conn) (error) {
	if request.Headers == nil {
		return errors.New("Connection: close")
	}
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
