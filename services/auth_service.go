package services

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"log"
	"mini-bank-api/models"
	"mini-bank-api/repository"
	"os"
	"time"
)

type AuthServiceInterface interface {
	Login(username, password string) (string, error)
	Logout(token string) error
}

type AuthService struct {
	customerRepo repository.CustomerRepositoryInterface
	historyRepo  repository.HistoryRepositoryInterface
}

func NewAuthService(customerRepo repository.CustomerRepositoryInterface, historyRepo repository.HistoryRepositoryInterface) AuthServiceInterface {
	return &AuthService{
		customerRepo: customerRepo,
		historyRepo:  historyRepo,
	}
}

func (auth *AuthService) Login(username, password string) (string, error) {
	customer, err := auth.customerRepo.LoadByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if customer.Username != username || customer.Password != password {
		return "", errors.New("invalid credentials")
	}

	// log
	err = auth.historyRepo.LogAction(customer.ID, "login")
	if err != nil {
		return "", err
	}

	token, err := generateJWT(customer)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (auth *AuthService) Logout(token string) error {
	err := auth.historyRepo.LogAction(token, "logout")
	if err != nil {
		return err
	}
	return nil
}

func generateJWT(user *models.Customer) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	secretKey := []byte(getSecretKey())
	accessToken, err := tokenClaims.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func getSecretKey() string {
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
