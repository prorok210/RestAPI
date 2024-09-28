package app

import (
	"RestAPI/core"
)

type HandlerFunc func(core.HttpRequest) core.HttpResponse

type funcInfo struct {
	HandlerFunc
	name string
}

var HandlersList = make(map[string]funcInfo)

func registerHandler(url string, f HandlerFunc, name ...string) {
	handlerName := ""
	if len(name) > 0 {
		handlerName = name[0]
	}
	HandlersList[url] = funcInfo{f, handlerName}
}

func router(url string) HandlerFunc {
	return HandlersList[url].HandlerFunc
}
