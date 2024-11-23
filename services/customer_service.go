package services

import (
	"mini-bank-api/models"
	"mini-bank-api/repository"
)

type CustomerService interface {
	GetAllCustomers() ([]models.Customer, error)
}

type CustomerServiceImpl struct {
	repo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) *CustomerServiceImpl {
	return &CustomerServiceImpl{repo: repo}
}

func (si *CustomerServiceImpl) GetAllCustomers() ([]models.Customer, error) {
	return si.repo.LoadAll()
}
