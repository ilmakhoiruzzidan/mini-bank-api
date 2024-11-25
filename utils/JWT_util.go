package utils

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"log"
	"mini-bank-api/models"
	"os"
	"time"
)

func GetSecretKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("Error loading JWT_SECRET_KEY file")
	}
	return secretKey
}

func GenerateJWT(user *models.Customer) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	secretKey := []byte(GetSecretKey())
	accessToken, err := tokenClaims.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(GetSecretKey()), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
