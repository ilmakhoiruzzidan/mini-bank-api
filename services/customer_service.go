package services

import (
	"mini-bank-api/models"
	"mini-bank-api/repository"
)

type CustomerServiceInterface interface {
	GetAllCustomers() ([]models.Customer, error)
}

type CustomerService struct {
	repo repository.CustomerRepositoryInterface
}

func NewCustomerService(repo repository.CustomerRepositoryInterface) CustomerServiceInterface {
	return &CustomerService{repo: repo}
}

func (si *CustomerService) GetAllCustomers() ([]models.Customer, error) {
	return si.repo.LoadAll()
}
