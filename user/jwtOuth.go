package user

import (
	"RestAPI/server"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


func GenerateSecretKey(length int) (string, error) {
    key := make([]byte, length)
    
    _, err := rand.Read(key)
    if err != nil {
        return "", err
    }
    
    return base64.URLEncoding.EncodeToString(key), nil
}


func GenerateAccessToken(mobile string, isActive bool) (string, error) {
	if mobile == "" {
		return "", errors.New("Mobile number is empty")
	}
	claims := jwt.MapClaims{
		"mobile": mobile,
		"isActive": isActive,
		"exp": time.Now().Add(server.JWT_ACCESS_EXPIRATION_TIME).Unix(),
		"token_type": "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(server.JWT_ACCESS_SECRET_KEY))
}

func GenerateRefreshToken(mobile string, isActive bool) (string, error) {
	if mobile == "" {
		return "", errors.New("Mobile number is empty")
	}
	claims := jwt.MapClaims{
		"mobile": mobile,
		"isActive": isActive,
		"exp": time.Now().Add(server.JWT_REFRESH_EXPIRATION_TIME).Unix(),
		"token_type": "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(server.JWT_REFRESH_SECRET_KEY))
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(server.JWT_ACCESS_SECRET_KEY), nil
	})

	if err != nil {
		token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(server.JWT_REFRESH_SECRET_KEY), nil
		})
		if err != nil {
			return nil, err
		}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func RefreshTokens(refreshTokenString string) (string, string, error) {
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(server.JWT_REFRESH_SECRET_KEY), nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["token_type"] != "refresh" {
		return "", "", errors.New("Invalid token type")
	}
	mobile, isActive := claims["mobile"].(string), claims["isActive"].(bool)

	newAccessToken, err := GenerateAccessToken(mobile, isActive)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := GenerateRefreshToken(mobile, isActive)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}