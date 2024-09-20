package user

import (
	"RestAPI/server"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

/*
	Test mobile.go
*/
func TestSendSms(t *testing.T) {
	testCases := []struct {
		message string
		Mobile string
		err error
	}{
		{
			message: "Hello",
			Mobile: "79936977511",
			err: nil,
		},
		{
			message: "",
			Mobile: "79936977511",
			err: errors.New("Message is empty"),
		},
		{
			message: "Hello",
			Mobile: "",
			err: errors.New("Mobile number is empty"),
		},
	}
	
	err := godotenv.Load("../.env")
	if err != nil {
		t.Errorf("Error env load %v", err)
	}
	server.MTS_API_KEY = os.Getenv("MTS_API_KEY")
	server.MTS_API_NUMBER = os.Getenv("MTS_API_NUMBER")
	if (server.MTS_API_KEY == "" || server.MTS_API_NUMBER == "") {
		t.Errorf("Error env load %v", err)
	}

	for _, tc := range testCases {
		u := &User{
			Mobile: tc.Mobile,
		}
		err := u.SendSMS(tc.message)
		if err != nil {
            if tc.err == nil {
                t.Errorf("Unexpected error: %v", err)
            } else if err.Error() != tc.err.Error() {
                t.Errorf("Expected error: %v, got: %v", tc.err, err)
            }
        } else if tc.err != nil {
            t.Errorf("Expected error: %v, but got none", tc.err)
        }
	}
}

/*
	Test jwtOuth.go
*/
func initSecrets() {
	server.JWT_ACCESS_SECRET_KEY = "testAccessSecretKey"
	server.JWT_REFRESH_SECRET_KEY = "testRefreshSecretKey"
	server.JWT_ACCESS_EXPIRATION_TIME = time.Minute * 5
	server.JWT_REFRESH_EXPIRATION_TIME = time.Minute * 10
}

func TestJWTFunctions(t *testing.T) {
	testCases := []struct {
		name        string
		mobile      string
		isActive    bool
		expectedErr error
		testFunc    func() error
	}{
		{
			name:     "Generate Secret Key",
			mobile:   "",
			isActive: false,
			expectedErr: nil,
			testFunc: func() error {
				_, err := GenerateSecretKey(32)
				return err
			},
		},
		{
			name:     "Generate Access Token",
			mobile:   "1234567890",
			isActive: true,
			expectedErr: nil,
			testFunc: func() error {
				_, err := GenerateAccessToken("1234567890", true)
				return err
			},
		},
		{
			name:     "Generate Refresh Token",
			mobile:   "1234567890",
			isActive: true,
			expectedErr: nil,
			testFunc: func() error {
				_, err := GenerateRefreshToken("1234567890", true)
				return err
			},
		},
		{
			name:     "Validate Access Token",
			mobile:   "1234567890",
			isActive: true,
			expectedErr: nil,
			testFunc: func() error {
				token, err := GenerateAccessToken("1234567890", true)
				if err != nil {
					return err
				}
				_, err = ValidateToken(token)
				return err
			},
		},
		{
			name:     "Invalid Access Token",
			mobile:   "",
			isActive: false,
			expectedErr: errors.New("error expected"),
			testFunc: func() error {
				_, err := ValidateToken("invalidToken")
				return err
			},
		},
		{
			name:     "Refresh Tokens",
			mobile:   "1234567890",
			isActive: true,
			expectedErr: nil,
			testFunc: func() error {
				refreshToken, err := GenerateRefreshToken("1234567890", true)
				if err != nil {
					return err
				}
				_, _, err = RefreshTokens(refreshToken)
				return err
			},
		},
		{
			name:     "Invalid Token Type for Refresh",
			mobile:   "1234567890",
			isActive: true,
			expectedErr: errors.New("error expected"),
			testFunc: func() error {
				accessToken, err := GenerateAccessToken("1234567890", true)
				if err != nil {
					return err
				}
				_, _, err = RefreshTokens(accessToken)
				return err
			},
		},
		{
			name:     "Expired Access Token",
			mobile:   "1234567890",
			isActive: true,
			expectedErr: errors.New("error expected"),
			testFunc: func() error {
				server.JWT_ACCESS_EXPIRATION_TIME = time.Millisecond * 100
				token, err := GenerateAccessToken("1234567890", true)
				if err != nil {
					return err
				}
				time.Sleep(time.Millisecond * 200)
				_, err = ValidateToken(token)
				server.JWT_ACCESS_EXPIRATION_TIME = time.Minute * 5
				return err
			},
		},
		{
			name:        "Generate Access Token with Empty Mobile",
			mobile:      "",
			isActive:    true,
			expectedErr: errors.New("Mobile number is empty"),
			testFunc: func() error {
				_, err := GenerateAccessToken("", true)
				return err
			},
		},
		{
			name:        "Generate Refresh Token with Empty Mobile",
			mobile:      "",
			isActive:    true,
			expectedErr: errors.New("Mobile number is empty"),
			testFunc: func() error {
				_, err := GenerateRefreshToken("", true)
				return err
			},
		},
		{
			name:        "Validate Invalid Token",
			mobile:      "",
			isActive:    false,
			expectedErr: errors.New("error expected"),
			testFunc: func() error {
				_, err := ValidateToken("invalidToken")
				return err
			},
		},
		{
			name:        "Refresh with Invalid Token Type",
			mobile:      "1234567890",
			isActive:    true,
			expectedErr: errors.New("Invalid token type"),
			testFunc: func() error {
				accessToken, err := GenerateAccessToken("1234567890", true)
				if err != nil {
					return err
				}
				_, _, err = RefreshTokens(accessToken)
				return err
			},
		},
	}


	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.testFunc()
			if tc.expectedErr != nil {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}