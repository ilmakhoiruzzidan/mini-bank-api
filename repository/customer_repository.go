package repository

import (
	"encoding/json"
	"errors"
	"mini-bank-api/models"
	"os"
)

type CustomerRepositoryInterface interface {
	LoadAll() ([]models.Customer, error)
	LoadByUsername(username string) (*models.Customer, error)
}

type JSONCustomerRepository struct {
	FilePath string
}

func NewJSONCustomerRepository() CustomerRepositoryInterface {
	return &JSONCustomerRepository{
		FilePath: "data/customers.json",
	}
}

func (repo *JSONCustomerRepository) LoadAll() ([]models.Customer, error) {
	data, err := os.ReadFile(repo.FilePath)
	if err != nil {
		return nil, err
	}

	var customers []models.Customer
	err = json.Unmarshal(data, &customers)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (repo *JSONCustomerRepository) LoadByUsername(username string) (*models.Customer, error) {
	var customers []models.Customer
	data, err := os.ReadFile(repo.FilePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &customers)
	if err != nil {
		return nil, err
	}

	for _, customer := range customers {
		if customer.Username == username {
			return &customer, nil
		}
	}
	return nil, errors.New("customer not found")
}
