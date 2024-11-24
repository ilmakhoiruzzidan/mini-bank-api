package services

import (
	"errors"
	"mini-bank-api/repository"
)

type TransactionServiceInterface interface {
	ProcessTransaction(senderID, merchantID string, amount float64) (string, error)
}

type TransactionService struct {
	transactionRepo repository.TransactionRepositoryInterface
	merchantRepo    repository.MerchantRepositoryInterface
	historyRepo     repository.HistoryRepositoryInterface
}

func NewTransactionService(
	transactionRepo repository.TransactionRepositoryInterface,
	merchantRepo repository.MerchantRepositoryInterface,
	historyRepo repository.HistoryRepositoryInterface) TransactionServiceInterface {
	return &TransactionService{
		transactionRepo: transactionRepo,
		merchantRepo:    merchantRepo,
		historyRepo:     historyRepo,
	}
}

func (transaction *TransactionService) ProcessTransaction(senderID, merchantID string, amount float64) (string, error) {
	if amount <= 0 {
		return "", errors.New("amount must be greater than 0")
	}

	merchants, err := transaction.merchantRepo.LoadAll()
	if err != nil {
		return "", err
	}

	merchantExists := false
	for _, merchant := range merchants {
		if merchant.ID == merchantID {
			merchantExists = true
			break
		}
	}
	if !merchantExists {
		return "", errors.New("merchant id does not exist")
	}

	transactionId, err := transaction.transactionRepo.SaveTransaction(senderID, merchantID, amount)
	if err != nil {
		return "", err
	}

	err = transaction.historyRepo.LogAction(senderID, "payment")
	if err != nil {
		return "", err
	}

	return transactionId, nil
}
