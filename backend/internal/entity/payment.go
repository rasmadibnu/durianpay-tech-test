package entity

import "time"

type Payment struct {
	ID           string    `json:"id"`
	MerchantID   int       `json:"merchant_id"`
	MerchantName string    `json:"merchant_name,omitempty"`
	Amount       string    `json:"amount"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}
