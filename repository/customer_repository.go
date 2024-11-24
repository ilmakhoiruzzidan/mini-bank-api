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
	SaveCustomers([]models.Customer) error
	FindCustomerByID(customerID string) (*models.Customer, error)
	MarkCustomerAsLoggedOut(customerID string) error
	MarkCustomerAsLoggedIn(customerID string) error
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

func (repo *JSONCustomerRepository) FindCustomerByID(customerID string) (*models.Customer, error) {
	customers, err := repo.LoadAll()
	if err != nil {
		return nil, err
	}

	for _, customer := range customers {
		if customer.ID == customerID {
			return &customer, nil
		}
	}

	return nil, errors.New("customer not found")
}

func (repo *JSONCustomerRepository) SaveCustomers(customers []models.Customer) error {
	file, err := os.Create(repo.FilePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	encoder := json.NewEncoder(file)
	err = encoder.Encode(customers)
	if err != nil {
		return err
	}

	return nil
}

func (repo *JSONCustomerRepository) MarkCustomerAsLoggedOut(customerID string) error {
	customers, err := repo.LoadAll()
	if err != nil {
		return err
	}

	for i, customer := range customers {
		if customer.ID == customerID {
			if customer.IsLoggedOut {
				return errors.New("customer has already logged out")
			}
			customers[i].IsLoggedOut = true
			err := repo.SaveCustomers(customers)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("customer not found")
}

func (repo *JSONCustomerRepository) MarkCustomerAsLoggedIn(username string) error {
	customers, err := repo.LoadAll()
	if err != nil {
		return err
	}

	for i, customer := range customers {
		if customer.Username == username {
			customers[i].IsLoggedOut = false
			err := repo.SaveCustomers(customers)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("customer not found")
}
