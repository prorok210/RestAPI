package user

import (
	"RestAPI/core"
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
	if u.Mobile == "" {
		return errors.New("Mobile number is empty")
	}
	if message == "" {
		return errors.New("Message is empty")
	}

	url := "https://api.exolve.ru/messaging/v1/SendSMS"

	smsJson := fmt.Sprintf(`{"number": "%s", "destination": "%s", "text": "%s"}`, core.MTS_API_NUMBER, u.Mobile, message)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(smsJson)))
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+core.MTS_API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.Status == "200 OK" {
		return nil
	} else {
		return errors.New("Error sending SMS")
	}
}
