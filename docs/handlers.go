package docs

import (
	"RestAPI/core"
	"log"
	"os"
)

func GetDocs(request core.HttpRequest) core.HttpResponse {
	doc, err := os.ReadFile("docs/docs.html")
	if err != nil {
		log.Println("Error reading docs file:", err)
		return core.HTTP500
	}

	response := core.HTTP200
	response.Body = string(doc)
	response.Headers = map[string]string{
		"Content-Type": "text/html",
	}
	return response
}

func GetDocsCSS(request core.HttpRequest) core.HttpResponse {
	css, err := os.ReadFile("docs/templates/css/styles.css")
	if err != nil {
		log.Println("Error reading docs css file:", err)
		return core.HTTP500
	}

	response := core.HTTP200
	response.Body = string(css)
	response.Headers = map[string]string{
		"Content-Type": "text/css",
	}
	return response
}

func GetDocsJS(request core.HttpRequest) core.HttpResponse {
	js, err := os.ReadFile("docs/templates/js/script.js")
	if err != nil {
		log.Println("Error reading docs js file:", err)
		return core.HTTP500
	}

	response := core.HTTP200
	response.Body = string(js)
	response.Headers = map[string]string{
		"Content-Type": "application/javascript",
	}
	return response
}
