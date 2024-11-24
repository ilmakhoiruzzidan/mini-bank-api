package models

import "time"

type History struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	Action     string    `json:"action"`
	Amount     int       `json:"amount"`
	Timestamp  time.Time `json:"timestamp"`
}
