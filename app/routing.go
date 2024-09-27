package app

import (
	"RestAPI/core"
)

type HandlerFunc func(core.HttpRequest) core.HttpResponse

type funcInfo struct {
	HandlerFunc
	name string
}

var viewsList = make(map[string]funcInfo)

func registerHandler(url string, f HandlerFunc, name ...string) {
	handlerName := ""
	if len(name) > 0 {
		handlerName = name[0]
	}
	viewsList[url] = funcInfo{f, handlerName}
}

func router(url string) HandlerFunc {
	return viewsList[url].HandlerFunc
}
