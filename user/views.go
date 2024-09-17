package user

import (
	"RestAPI/server"
	"encoding/json"
	"fmt"
	"strconv"
)

func HelloView(request server.HttpRequest) (server.HttpResponse) {
	if request.Method == "GET" {
		response := server.HTTP200
		response.Body = `{"Hello, world!"}`
		response.Headers["Content-Length"] = strconv.Itoa(len(response.Body))
		return response
	} else {
		response := server.HTTP405
		return response
	}
}

func GoodbyeView(request server.HttpRequest) (server.HttpResponse) {
	if request.Method == "GET" {
		response := server.HTTP200
		response.Body = `{"Goodbye, world!"}`
		response.Headers["Content-Length"] = strconv.Itoa(len(response.Body))
		return response
	} else {
		response := server.HTTP405
		return response
	}
}

func AddView(request server.HttpRequest) (server.HttpResponse) {
	if request.Method == "POST" {
		response := server.HTTP201
		response.Body = `{"Addition"}`
		response.Headers["Content-Length"] = strconv.Itoa(len(response.Body))
		return response
	} else {
		response := server.HTTP405
		return response
	}
}

func CreateUserView(request server.HttpRequest) (server.HttpResponse) {
	if request.Method != "POST" {
		return server.HTTP405
	}

	jsonString := request.Body

	fmt.Println("jsonString:", jsonString)

	user := new(User)

	err := json.Unmarshal([]byte(jsonString), user)
	if err != nil {
		response := server.HTTP400
		response.Body = `{"Message": "Invalid JSON"}`
		response.Headers["Content-Length"] = strconv.Itoa(len(response.Body))
		return response
	}

	otp := genareteOTP()
	user.Otp = otp

	err = user.SendSMS(fmt.Sprintf("Your OTP is %s", otp))
	if err != nil {
		fmt.Println("Error sending SMS:", err)
		response := server.HTTP500
		response.Body = `{"Message": "Internal server error"}`
		response.Headers["Content-Length"] = strconv.Itoa(len(response.Body))
		return response
	}

	userStore[user.Mobile] = user

	response := server.HTTP201
	response.Body = `{"Message": "User created, please verify"}`
	response.Headers["Content-Length"] = strconv.Itoa(len(response.Body))
	return response
}

func VerifyUserView(request server.HttpRequest) (server.HttpResponse) {
	if request.Method != "POST" {
		return server.HTTP405
	}

	jsonString := request.Body

	tmpUser := new(User)

	err := json.Unmarshal([]byte(jsonString), tmpUser)
	if err != nil {
		response := server.HTTP400
		response.Body = `{"Message": "Invalid JSON"}`
		response.Headers["Content-Length"] = strconv.Itoa(len(response.Body))
		return response
	}

	if userStore[tmpUser.Mobile].Otp == tmpUser.Otp {
		userStore[tmpUser.Mobile].isActive = true
		response := server.HTTP200
		response.Body = `{"Message": "User verified"}`
		response.Headers["Content-Length"] = strconv.Itoa(len(response.Body))
		return response
	} else {
		response := server.HTTP401
		response.Body = `{"Message": "User not verified"}`
		response.Headers["Content-Length"] = strconv.Itoa(len(response.Body))
		return response
	}
}