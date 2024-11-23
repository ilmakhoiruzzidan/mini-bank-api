package services

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"mini-bank-api/models"
	"mini-bank-api/repository"
	"time"
)

type AuthService interface {
	Login(username, password string) (string, error)
}

type AuthServiceImpl struct {
	customerRepo repository.CustomerRepository
}

func NewAuthService(customerRepo repository.CustomerRepository) *AuthServiceImpl {
	return &AuthServiceImpl{customerRepo: customerRepo}
}

func (auth *AuthServiceImpl) Login(username, password string) (string, error) {
	// Here, you will validate the user's credentials
	// For simplicity, let's assume we have a mock validation.
	customers, err := auth.customerRepo.LoadAll()
	if err != nil {
		return "", err
	}

	var user *models.Customer
	for _, customer := range customers {
		if customer.Username == username && customer.Password == password {
			user = &customer
			break
		}
	}

	if user == nil {
		return "", errors.New("invalid credentials")
	}

	// Generate a token (JWT, for example)
	token, err := generateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func generateJWT(user *models.Customer) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
	})

	tokenString, err := tokenClaims.SignedString([]byte("secret_key")) // Replace with your secret key
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
