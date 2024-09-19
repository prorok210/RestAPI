package user

import (
	"RestAPI/server"
	"encoding/json"
	"fmt"
)


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
		return response
	}

	otp := genareteOTP()
	user.Otp = otp

	err = user.SendSMS(fmt.Sprintf("Your OTP is %s", otp))
	if err != nil {
		fmt.Println("Error sending SMS:", err)
		response := server.HTTP500
		response.Body = `{"Message": "Internal server error"}`
		return response
	}

	userStore[user.Mobile] = user

	response := server.HTTP201
	response.Body = `{"Message": "User created, please verify"}`
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