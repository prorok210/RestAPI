package user

import (
	"RestAPI/core"
	"RestAPI/db"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

type User struct {
	db.User
}

/*
docs(

	tag: user;
	name: CreateUser;
	path: /user/create;
	method: DELETE;
	summary: Create a new user;
	description: Create a new user with the given data;
	isAuth: true;

)docs
*/
func CreateUserHandler(request core.HttpRequest) core.HttpResponse {
	if request.Method != "POST" {
		return core.HTTP405
	}

	jsonString := request.Body

	user := new(User)

	err := json.Unmarshal([]byte(jsonString), user)
	if err != nil {
		response := core.HTTP400
		response.Body = `{"Message": "Invalid JSON"}`
		return response
	}

	otp := genareteOTP()
	user.Otp = otp

	err = user.SendSMS(fmt.Sprintf("Your OTP is %s", otp))
	if err != nil {
		log.Println("Error sending SMS:", err)
		response := core.HTTP500
		response.Body = `{"Message": "Internal server error"}`
		return response
	}

	response := core.HTTP201
	response.Body = `{"Message": "User created, please verify"}`
	return response
}

/*
docs(

	name: VerifyUserHandler;
	tag: user;
	path: /user/verify;
	method: GET;
	сontent_type: application/json;
	summary: Verify user;
	description: Verify user with the given data;
	isAuth: false;
	req_content_types: application/json;
	requestbody: {
		"mobile": "string",
		"otp": "string"
	};
	resp_content_type: application/json;
	responsebody: {
		"Message": "User verified"
	};

)docs
*/
func VerifyUserHandler(request core.HttpRequest) core.HttpResponse {
	if request.Method != "POST" {
		return core.HTTP405
	}

	if request.Body == "" {
		response := core.HTTP400
		response.Body = `{"Message": "Invalid JSON"}`
		return response
	}

	jsonString := request.Body

	tmpUser := new(User)

	err := json.Unmarshal([]byte(jsonString), tmpUser)
	if err != nil {
		response := core.HTTP400
		response.Body = `{"Message": "Invalid JSON"}`
		return response
	}

	return core.HTTP200
}

/*
docs(

	name: CreateUserFormdataHandler;
	tag: user;
	path: /user/createformdata;
	method: POST;
	summary: Create a new user with form data;
	description: Create a new user with the given data and images;
	isAuth: false;
	req_content_type: multipart/form-data;
	requestbody: {
		"name": "string",
		"surname": "string",
		"mobile": "string",
		"email": "string",
		"age": "int"
	};
	resp_content_type: application/json;
	responsebody: {
		"name": "string",
		"surname": "string",
		"mobile": "string",
		"email": "string",
		"age": "int"
	};

)docs
*/
func CreateUserFormdataHandler(request core.HttpRequest) core.HttpResponse {
	if request.Method != "POST" {
		return core.HTTP405
	}

	user := new(User)

	user.Mobile = request.FormData.Fields["mobile"]
	user.Name = request.FormData.Fields["name"]
	user.Surname = request.FormData.Fields["surname"]
	user.Email = request.FormData.Fields["email"]

	age, err := strconv.Atoi(request.FormData.Fields["age"])
	if err != nil {
		fmt.Println("Ошибка при конвертации возраста:", err)
		return core.HTTP400
	}
	user.Age = age

	saveFile := func(filename string, fileData []byte) error {
		currentDir, er := os.Getwd()
		filePath := currentDir + core.IMAGES_DIR + "/" + filename
		if er != nil {
			return er
		}
		if _, err := os.Stat(currentDir + core.IMAGES_DIR); os.IsNotExist(err) {
			err := os.MkdirAll(currentDir+core.IMAGES_DIR, 0755)
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
			return core.HTTP500
		}
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error marshaling user:", err)
		return core.HTTP500
	}
	response := core.HTTP201
	response.Body = string(jsonData)
	return response
}

func ImageHandler(request core.HttpRequest) core.HttpResponse {
	if request.Method != "GET" {
		return core.HTTP405
	}

	currentDir, er := os.Getwd()
	if er != nil {
		fmt.Println("Error getting current directory:", er)
		return core.HTTP500
	}

	filename := request.Query["filename"]

	filePath := currentDir + core.IMAGES_DIR + "/" + filename
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return core.HTTP404
	}

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return core.HTTP500
	}

	response := core.HTTP200
	response.Body = string(fileData)
	response.SetHeader("Content-Type", "image/jpeg")
	return response
}
