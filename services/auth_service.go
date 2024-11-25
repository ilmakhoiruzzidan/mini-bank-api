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

	if customer.IsLoggedOut {
		return nil, errors.New("you are already logged out")
	}
	return customer, nil
}
