package repository

import (
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"mini-bank-api/models"
	"os"
	"time"
)

type HistoryRepositoryInterface interface {
	LogAction(customerID, action string) error
}

type JSONHistoryRepository struct {
	FilePath string
}

func NewJSONHistoryRepository() HistoryRepositoryInterface {
	return &JSONHistoryRepository{
		FilePath: "data/histories.json",
	}
}

func (repo *JSONHistoryRepository) LogAction(customerID, action string) error {
	var histories []models.History
	data, err := os.ReadFile(repo.FilePath)
	if err == nil {
		_ = json.Unmarshal(data, &histories)
	}

	newHistory := models.History{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		Action:     action,
		Timestamp:  time.Now(),
	}

	histories = append(histories, newHistory)
	historyData, err := json.MarshalIndent(histories, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(repo.FilePath, historyData, os.ModePerm)
}
