package model

import "time"

type History struct {
	TransactionId    string    `json:"transaction_id"`
	CustomerUsername string    `json:"customer_username"`
	MerchantCode     string    `json:"merchant_code"`
	Amount           float64   `json:"amount"`
	Date             time.Time `json:"date"`
}
