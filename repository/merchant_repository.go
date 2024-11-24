package repository

import (
	"errors"
	"github.com/goccy/go-json"
	"mini-bank-api/models"
	"os"
)

type MerchantRepositoryInterface interface {
	LoadAll() ([]models.Merchant, error)
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
