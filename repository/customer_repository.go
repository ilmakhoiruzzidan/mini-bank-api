package repository

import (
	"encoding/json"
	"mini-bank-api/models"
	"os"
)

type CustomerRepository interface {
	LoadAll() ([]models.Customer, error)
}

type JSONCustomerRepository struct {
	FilePath string
}

func NewJSONCustomerRepository() *JSONCustomerRepository {
	return &JSONCustomerRepository{}
}

func (repo *JSONCustomerRepository) LoadAll() ([]models.Customer, error) {
	data, err := os.ReadFile("data/customers.json")
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
