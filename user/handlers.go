package user

import (
	"RestAPI/server"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

func CreateUserHandler(request server.HttpRequest) server.HttpResponse {
	if request.Method != "POST" {
		return server.HTTP405
	}

	jsonString := request.Body

	user := new(User)

	err := json.Unmarshal([]byte(jsonString), user)
	if err != nil {
		response := server.HTTP400
		response.Body = `{"Message": "Invalid JSON"}`
		return response
	}

	otp := genareteOTP()
	user.Otp = otp

	err = user.SendSMS(fmt.Sprintf("Your OTP is %s", otp))
	if err != nil {
		log.Println("Error sending SMS:", err)
		response := server.HTTP500
		response.Body = `{"Message": "Internal server error"}`
		return response
	}

	userStore[user.Mobile] = user

	response := server.HTTP201
	response.Body = `{"Message": "User created, please verify"}`
	return response
}

func VerifyUserHandler(request server.HttpRequest) server.HttpResponse {
	if request.Method != "POST" {
		return server.HTTP405
	}

	if request.Body == "" {
		response := server.HTTP400
		response.Body = `{"Message": "Invalid JSON"}`
		return response
	}

	jsonString := request.Body

	tmpUser := new(User)

	err := json.Unmarshal([]byte(jsonString), tmpUser)
	if err != nil {
		response := server.HTTP400
		response.Body = `{"Message": "Invalid JSON"}`
		return response
	}

	if len(userStore) == 0 {
		response := server.HTTP400
		response.Body = `{"Message": "User not found"}`
		return response
	}

	if userStore[tmpUser.Mobile].Otp == tmpUser.Otp {
		userStore[tmpUser.Mobile].isActive = true
		response := server.HTTP200
		response.Body = `{"Message": "User verified"}`
		return response
	} else {
		response := server.HTTP401
		response.Body = `{"Message": "User not verified"}`
		return response
	}
}


func CreateUserFormdataHandler(request server.HttpRequest) server.HttpResponse {
	if request.Method != "POST" {
		return server.HTTP405
	}

	user := new(User)
	
	user.Mobile		= request.FormData.Fields["mobile"]
	user.Name 		= request.FormData.Fields["name"]
	user.Surname 	= request.FormData.Fields["surname"]
	user.Email		= request.FormData.Fields["email"]

	age, err := strconv.Atoi(request.FormData.Fields["age"])
	if err != nil {
		fmt.Println("Ошибка при конвертации возраста:", err)
		return server.HTTP400
	}
	user.Age = age

	saveFile := func(filename string, fileData []byte) error {
		currentDir , er := os.Getwd()
		filePath := currentDir + server.IMAGES_DIR + "/" + filename
		if er != nil {
			return er
		}
		if _, err := os.Stat(currentDir + server.IMAGES_DIR); os.IsNotExist(err) {
			err := os.MkdirAll(currentDir + server.IMAGES_DIR, 0755)
			if err != nil {
				return err
			}
		}

		file, er := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
		if er != nil {
			return er
		}
		defer file.Close()

		_, er = file.Write(fileData)
		if er != nil {
			return er
		}
		return nil
	}

	for _, fileData := range request.FormData.Files["images"] {
		er := saveFile(fileData.FileName, fileData.FileData)
		if er != nil {
			fmt.Println("Error saving file:", er)
			return server.HTTP500
		}
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error marshaling user:", err)
		return server.HTTP500
	}
	response := server.HTTP201
	response.Body = string(jsonData)
	return response
}


func ImageHandler(request server.HttpRequest) server.HttpResponse {
	if request.Method != "GET" {
		return server.HTTP405
	}

	currentDir , er := os.Getwd()
	if er != nil {
		fmt.Println("Error getting current directory:", er)
		return server.HTTP500
	}

	filename := request.Query["filename"]

	filePath := currentDir + server.IMAGES_DIR + "/" + filename
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return server.HTTP404
	}

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return server.HTTP500
	}

	response := server.HTTP200
	response.Body = string(fileData)
	response.SetHeader("Content-Type", "image/jpeg")
	return response
}