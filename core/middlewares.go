package core

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func isAllowedHostMiddleware(clientAddr string) bool {
	if !IS_ALLOWED_HOSTS {
		return true
	}

	host, _, err := net.SplitHostPort(clientAddr)
	if err != nil {
		host = clientAddr
	}

	fmt.Println("Host:", host)
	for _, allowedHost := range ALLOWED_HOSTS {
		if host == allowedHost {
			return true
		}
		allowedIP := net.ParseIP(allowedHost)
		if allowedIP != nil {
			clientIP := net.ParseIP(host)
			if clientIP != nil && clientIP.Equal(allowedIP) {
				return true
			}
		}
	}

	return false
}

func reqMiddleware(request *HttpRequest, clientConn Conn) error {
	if !REQ_MIDDLEWARE {
		return nil
	}
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
		if strings.Split(request.Headers["Content-Type"], ";")[0] == supportedMediaType {
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
		contentType, hasContentType := request.Headers["Content-Type"]
		if hasContentType && contentType == "multipart/form-data" || contentType == "application/x-www-form-urlencoded" {
			if contentLength < 0 {
				clientConn.Write(HTTP411.ToBytes())
				return errors.New("Content-Length required")
			}
		} else {
			if contentLength != len(request.Body) {
				clientConn.Write(HTTP411.ToBytes())
				return errors.New("Content-Length does not match body length")
			}
		}

	} else if len(request.Body) > 0 {
		clientConn.Write(HTTP411.ToBytes())
		return errors.New("Content-Length required")
	}

	return nil
}

func keepAliveMiddleware(request *HttpRequest, clientConn Conn) error {
	if !KEEP_ALIVE {
		return nil
	}
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
