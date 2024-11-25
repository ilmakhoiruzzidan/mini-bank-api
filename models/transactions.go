package models

import "time"

type Transaction struct {
	ID         string    `json:"transaction_id"`
	SenderID   string    `json:"sender_id"`
	MerchantID string    `json:"merchant_id"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}
