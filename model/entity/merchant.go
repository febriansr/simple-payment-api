package model

type Merchant struct {
	Uuid         string `json:"uuid"`
	MerchantCode string `json:"merchant_code"`
	Name         string `json:"name"`
}
