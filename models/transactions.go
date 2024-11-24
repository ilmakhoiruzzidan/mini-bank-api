package models

import "time"

type Transaction struct {
	ID         string    `json:"transaction_id"`
	SenderID   string    `json:"sender_id"`
	MerchantID string    `json:"merchant_id"`
	Amount     float64   `json:"amount"`
	Timestamp  time.Time `json:"timestamp"`
}
