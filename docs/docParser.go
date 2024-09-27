package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type HandlerInfo struct {
	tag          string
	name         string
	path         string
	method       string
	summary      string
	description  string
	isAuth       string
	requestbody  string
	responsebody string
}

var apps = []string{
	"user",
}

var handlers = make(map[string]*HandlerInfo)

func parseDocs() error {
	for _, app := range apps {
		err := filepath.Walk("../"+app, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".go" {
				parseFile(path)
			}
			return nil
		})
		return err
	}
	return nil
}

func parseFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	content := string(data)

	re := regexp.MustCompile(`(?s)docs\((.*?)\)docs`)

	matches := re.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			handlerInfo := new(HandlerInfo)
			for _, m := range strings.Split(match[1], ";") {
				parts := strings.SplitN(m, ":", 2)
				if len(parts) != 2 {
					continue
				}
				fmt.Println("Parts:", parts)
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				switch key {
				case "name":
					handlerInfo.name = value
				case "tag":
					handlerInfo.tag = value
				case "path":
					handlerInfo.path = value
				case "method":
					handlerInfo.method = value
				case "summary":
					handlerInfo.summary = value
				case "description":
					handlerInfo.description = value
				case "isAuth":
					handlerInfo.isAuth = value
				case "requestbody":
					handlerInfo.requestbody = value
				case "responsebody":
					handlerInfo.responsebody = value
				}

				if handlerInfo.name != "" && handlerInfo.path != "" && handlerInfo.method != "" {
					if _, ok := handlers[handlerInfo.name]; ok {
						fmt.Println("Handler already exists", handlerInfo.name)
						continue
					}
					handlers[handlerInfo.name] = handlerInfo
				}
			}
		}
	}
	return
}

func main() {
	err := parseDocs()
	if err != nil {
		fmt.Println("Error parsing docs", err)
		return
	}
	fmt.Println("Handlers:", handlers)
}
