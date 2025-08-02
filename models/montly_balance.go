package models

type MonthlyBalance struct {
	Month string `json:"month"`
	Year  int    `json:"year"`
	Saldo int64  `json:"saldo"`
}
