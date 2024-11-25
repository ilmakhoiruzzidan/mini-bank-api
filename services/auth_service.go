package services

import (
	"errors"
	"mini-bank-api/models"
	"mini-bank-api/repository"
	"mini-bank-api/utils"
)

type AuthServiceInterface interface {
	Login(username, password string) (string, error)
	Logout(token string) error
	GetCurrentUserInfo(token string) (*models.Customer, error)
}

type AuthService struct {
	customerRepo repository.CustomerRepositoryInterface
	historyRepo  repository.HistoryRepositoryInterface
	tokenRepo    repository.TokenRepositoryInterface
}

func NewAuthService(customerRepo repository.CustomerRepositoryInterface,
	historyRepo repository.HistoryRepositoryInterface,
	tokenRepo repository.TokenRepositoryInterface) AuthServiceInterface {
	return &AuthService{
		customerRepo: customerRepo,
		historyRepo:  historyRepo,
		tokenRepo:    tokenRepo,
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

	accessToken, err := utils.GenerateJWT(customer)

	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (auth *AuthService) Logout(token string) error {
	claims, err := utils.ParseToken(token)
	if err != nil {
		return err
	}

	if auth.tokenRepo.IsTokenRevoked(token) {
		return errors.New("token already revoked")
	}

	err = auth.tokenRepo.AddToRevocationList(token)
	if err != nil {
		return err
	}

	// Log action
	customerID, _ := claims["id"].(string)
	err = auth.historyRepo.LogAction(customerID, "logout")
	if err != nil {
		return err
	}

	return nil
}

func (auth *AuthService) GetCurrentUserInfo(token string) (*models.Customer, error) {
	claims, err := utils.ParseToken(token)
	if err != nil {
		return nil, err
	}

	CustomerID, ok := claims["id"].(string)
	if !ok {
		return nil, errors.New("failed to parse ID")
	}

	customer, err := auth.customerRepo.FindCustomerByID(CustomerID)
	if err != nil {
		return nil, err
	}

	err = auth.historyRepo.LogAction(CustomerID, "check profile")
	if err != nil {
		return nil, err
	}

	return customer, nil
}
