package services

import (
	"errors"
	"fmt"
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

	if !customer.IsLoggedOut {
		return "", errors.New("you already logged in")
	}

	if customer.Username != username || customer.Password != password {
		return "", errors.New("invalid credentials")
	}

	err = auth.customerRepo.MarkCustomerAsLoggedIn(username)
	if err != nil {
		return "", err
	}

	// log
	err = auth.historyRepo.LogAction(customer.ID, "login")
	if err != nil {
		return "", err
	}

	accessToken, err := generateJWT(customer)

	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (auth *AuthService) Logout(token string) error {
	claims, err := parseToken(token)
	if err != nil {
		return err
	}

	CustomerID, ok := claims["id"].(string)
	if !ok {
		return errors.New("failed to parse ID")
	}

	customer, err := auth.customerRepo.FindCustomerByID(CustomerID)
	if err != nil {
		return err
	}

	if customer.IsLoggedOut {
		return errors.New("you are already logged out")
	}

	err = auth.historyRepo.LogAction(CustomerID, "logout")
	if err != nil {
		return err
	}

	err = auth.customerRepo.MarkCustomerAsLoggedOut(CustomerID)
	if err != nil {
		return err
	}

	return nil
}

func generateJWT(user *models.Customer) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(10 * time.Minute).Unix(),
	})

	secretKey := []byte(getSecretKey())
	accessToken, err := tokenClaims.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(getSecretKey()), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
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
