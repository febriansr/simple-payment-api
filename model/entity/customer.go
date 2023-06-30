package model

type Customer struct {
	Uuid     string  `json:"uuid"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Balance  float64 `json:"balance"`
}
