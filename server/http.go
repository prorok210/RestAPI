package server

import (
	"errors"
	"strconv"
	"strings"
)


type HttpRequest struct {
	Method string
	Url string
	Version string
	Headers map[string]string
	Body string
}


type HttpResponse struct {
	Version string
	Status int
	Reason string
	Headers map[string]string
	Body string
}


func (rqst *HttpRequest) ParseRequest(buffer []byte) (error) {
	if len(buffer) == 0 {
		return errors.New("Empty request")
	}

	reqStr := string(buffer)
	lines := strings.Split(reqStr, "\r\n")
	if len(lines) < 1 {
		return errors.New("Invalid request: no lines found")
	}

	requestLine := strings.Split(lines[0], " ")
	if len(requestLine) < 3 {
		return errors.New("Invalid request line")
	}
	rqst.Method = requestLine[0]
	rqst.Url = requestLine[1]
	rqst.Version = requestLine[2]

	rqst.Headers = make(map[string]string)

	i := 1
	for i < len(lines) && lines[i] != "" {
		headerParts := strings.SplitN(lines[i], ": ", 2)
		if len(headerParts) == 2 {
			rqst.Headers[headerParts[0]] = headerParts[1]
		}
		i++
	}

	if i+1 < len(lines) {
		rqst.Body = strings.Join(lines[i+1:], "\r\n")
	} else {
		rqst.Body = ""
	}

	return nil
}


func (rqst *HttpRequest) ToString() string {
	reqStr := rqst.Method + " " + rqst.Url + " " + rqst.Version + "\r\n"
	for key, value := range rqst.Headers {
		reqStr += key + ": " + value + "\r\n"
	}
	reqStr += "\r\n" + rqst.Body

	return reqStr
}


func (resp *HttpResponse) ToString() string {
	respStr := resp.Version + " " + strconv.Itoa(resp.Status) + " " + resp.Reason + "\r\n"
	for key, value := range resp.Headers {
		respStr += key + ": " + value + "\r\n"
	}
	respStr += "\r\n" + resp.Body

	return respStr
}


func (resp *HttpResponse) ToBytes() []byte {
	return []byte(resp.ToString())
}