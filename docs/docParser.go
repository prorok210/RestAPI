package docs

import (
	"RestAPI/core"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type HandlerInfo struct {
	Tag             string
	Name            string
	Path            string
	Method          string
	ReqContentTypes []string
	RespContentType string
	Summary         string
	Description     string
	IsAuth          string
	RequestBody     string
	FormDataBody    map[string]string
	ResponseBody    string
}

type PageData struct {
	CSSPath  string
	JSPath   string
	Handlers map[string][]HandlerInfo
	Lower    func(string) string
}

var handlers = make(map[string]*HandlerInfo)

var groupedHandlers = make(map[string][]HandlerInfo)

func parseDocs() error {
	for _, app := range core.APPS {
		err := filepath.Walk(app, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".go" {
				parseFile(path)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	for _, handler := range handlers {
		groupedHandlers[handler.Tag] = append(groupedHandlers[handler.Tag], *handler)
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
			jsonIn := false
			formDataIn := false
			for _, m := range strings.Split(match[1], ";") {
				parts := strings.SplitN(m, ":", 2)
				if len(parts) != 2 {
					continue
				}

				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				key = strings.ToLower(key)

				switch key {
				case "name":
					handlerInfo.Name = value
				case "tag":
					handlerInfo.Tag = value
				case "path":
					handlerInfo.Path = value
				case "method":
					for _, method := range core.ALLOWED_METHODS {
						if value == method {
							handlerInfo.Method = value
							break
						}
					}
				case "summary":
					handlerInfo.Summary = value
				case "description":
					handlerInfo.Description = value
				case "isauth":
					if value == "true" || value == "false" {
						handlerInfo.IsAuth = value
					}
				case "req_content_type":
					for _, ct := range strings.Split(value, ",") {
						ct = strings.TrimSpace(ct)
						for _, contentType := range core.SUPPORTED_MEDIA_TYPES {
							if ct == contentType {
								handlerInfo.ReqContentTypes = append(handlerInfo.ReqContentTypes, ct)
								if ct == "application/json" {
									jsonIn = true
								}
								if ct == "multipart/form-data" {
									formDataIn = true
								}
								if ct == "application/x-www-form-urlencoded" {
									formDataIn = true
								}
							}
						}
					}
				case "req_content_types":
					for _, ct := range strings.Split(value, ",") {
						ct = strings.TrimSpace(ct)
						for _, contentType := range core.SUPPORTED_MEDIA_TYPES {
							if ct == contentType {
								handlerInfo.ReqContentTypes = append(handlerInfo.ReqContentTypes, ct)
								if ct == "application/json" {
									jsonIn = true
								}
								if ct == "multipart/form-data" {
									formDataIn = true
								}
								if ct == "application/x-www-form-urlencoded" {
									formDataIn = true
								}
							}
						}
					}
				case "requestbody":
					if jsonIn {
						lines := strings.Split(value, "\n")
						for _, line := range lines {
							trimmedLine := strings.TrimPrefix(line, "	")
							handlerInfo.RequestBody += trimmedLine + "\n"
						}
					}
					if formDataIn {
						reg := regexp.MustCompile(`\{([^}]*)\}`)
						matches := reg.FindStringSubmatch(value)
						if len(matches) == 1 {
							handlerInfo.FormDataBody = make(map[string]string)
							lines := strings.Split(value, "\n")
							for _, line := range lines {
								trimmedLine := strings.TrimPrefix(line, "	")
								parts := strings.Split(trimmedLine, ":")
								if len(parts) == 2 {
									handlerInfo.FormDataBody[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
								}
							}
						}
					}
				case "resp_content_type":
					handlerInfo.RespContentType = value
				case "responsebody":
					lines := strings.Split(value, "\n")
					for _, line := range lines {
						trimmedLine := strings.TrimPrefix(line, "	")
						handlerInfo.ResponseBody += trimmedLine + "\n"
					}
				}

				if handlerInfo.Name != "" && handlerInfo.Path != "" && handlerInfo.Method != "" {
					if _, ok := handlers[handlerInfo.Name]; ok {
						continue
					}
					handlers[handlerInfo.Name] = handlerInfo
				}

				for _, handler := range handlers {
					if handler.Tag == "" {
						handler.Tag = "Others"
					}
				}

				jsonIn = false
				formDataIn = false
			}
		}
	}
	return
}

func GenerateDocs() error {
	err := parseDocs()
	if err != nil {
		return err
	}

	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}

	tmpl, err := template.New("template.html").Funcs(funcMap).ParseFiles("docs/templates/html/template.html")
	if err != nil {
		return err
	}

	data := PageData{
		CSSPath:  "docs/templates/css/styles.css",
		JSPath:   "docs/templates/js/script.js",
		Handlers: groupedHandlers,
	}

	f, err := os.Create("docs/docs.html")
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, data)
	if err != nil {
		return err
	}

	return nil
}
