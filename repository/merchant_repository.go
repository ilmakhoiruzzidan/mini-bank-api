package repository

import (
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"mini-bank-api/models"
	"os"
)

type MerchantRepositoryInterface interface {
	LoadAll() ([]models.Merchant, error)
	SaveMerchants(merchants []models.Merchant) error
	UpdateMerchantBalance(merchants []models.Merchant, merchantID string, amount float64) error
	FindMerchantByID(merchantID string) (*models.Merchant, error)
}

type JSONMerchantRepository struct {
	FilePath string
}

func NewJSONMerchantRepository() MerchantRepositoryInterface {
	return &JSONMerchantRepository{
		FilePath: "data/merchants.json",
	}
}

func (repo *JSONMerchantRepository) LoadAll() ([]models.Merchant, error) {
	var merchants []models.Merchant
	data, err := os.ReadFile(repo.FilePath)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return merchants, nil
	}
	err = json.Unmarshal(data, &merchants)
	if err != nil {
		return nil, err
	}

	return merchants, nil
}

func (repo *JSONMerchantRepository) FindMerchantByID(merchantID string) (*models.Merchant, error) {
	merchants, err := repo.LoadAll()
	if err != nil {
		return nil, err
	}

	for _, merchant := range merchants {
		if merchant.ID == merchantID {
			return &merchant, nil
		}
	}
	return nil, errors.New("merchant not found")
}

func (repo *JSONMerchantRepository) SaveMerchants(merchants []models.Merchant) error {
	file, err := os.Create(repo.FilePath)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)

	err = json.NewEncoder(file).Encode(merchants)
	if err != nil {
		return err
	}
	return nil
}

func (repo *JSONMerchantRepository) UpdateMerchantBalance(merchants []models.Merchant, merchantID string, amount float64) error {
	for i, merchant := range merchants {
		if merchant.ID == merchantID {
			merchants[i].Balance += amount
			break
		}
	}

	err := repo.SaveMerchants(merchants)

	if err != nil {
		return err
	}

	return err
}
