package user

import (
	"RestAPI/server"
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)


func genareteOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(999999)
	return fmt.Sprintf("%06d", otp) 
}


func (u *User) SendSMS(message string) error {
	url := "https://api.exolve.ru/messaging/v1/SendSMS"

	smsJson := fmt.Sprintf(`{"number": "%s", "destination": "%s", "text": "%s"}`, server.MTS_API_NUMBER, u.Mobile, message)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(smsJson)))
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + server.MTS_API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return err
	}
	log.Println(resp)
	defer resp.Body.Close()
	log.Println("Response Status:", resp.Status)

	if resp.Status == "200 OK" {
		return nil
	} else {
		return errors.New("Error sending SMS")
	}
}