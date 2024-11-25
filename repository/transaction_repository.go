package repository

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"io"
	"mini-bank-api/models"
	"os"
	"time"
)

type TransactionRepositoryInterface interface {
	SaveTransaction(senderID, merchantID string, amount float64) (string, error)
	LoadTransactions() ([]models.Transaction, error)
}
type JSONTransactionRepository struct {
	FilePath string
}

func NewJSONTransactionRepository() TransactionRepositoryInterface {
	return &JSONTransactionRepository{
		FilePath: "data/transactions.json",
	}
}

func (repo *JSONTransactionRepository) SaveTransaction(senderId, merchantID string, amount float64) (string, error) {
	transactions, err := repo.LoadTransactions()
	if err != nil {
		return "", err
	}

	newTransaction := models.Transaction{
		ID:         uuid.New().String(),
		SenderID:   senderId,
		MerchantID: merchantID,
		Amount:     amount,
		Timestamp:  time.Now(),
	}
	transactions = append(transactions, newTransaction)

	file, err := os.Create(repo.FilePath)
	if err != nil {
		return "", err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Println("Error closing file:", cerr)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(transactions)
	if err != nil {
		return "", err
	}

	return newTransaction.ID, nil
}

func (repo *JSONTransactionRepository) LoadTransactions() ([]models.Transaction, error) {
	file, err := os.Open(repo.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Transaction{}, nil
		}
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	var transactions []models.Transaction
	err = json.NewDecoder(file).Decode(&transactions)
	if err == io.EOF {
		return []models.Transaction{}, nil
	}

	return transactions, nil
}
