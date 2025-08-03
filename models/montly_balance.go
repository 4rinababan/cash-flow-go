package models

type MonthlyBalance struct {
	Month   string `json:"month"`
	Year    int    `json:"year"`
	Income  int64  `json:"income"`
	Expense int64  `json:"expense"`
	Saldo   int64  `json:"saldo"`
}
