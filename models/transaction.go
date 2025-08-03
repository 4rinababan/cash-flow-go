package models

import (
	"time"

	"github.com/lib/pq"
)

type Transaction struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Type           string         `json:"type"` // pemasukan / pengeluaran
	Amount         float64        `json:"amount"`
	Note           string         `json:"note"`
	Category       string         `json:"category"`
	Categories     pq.StringArray `gorm:"type:text[]" json:"-"` // hide dari Swagger
	CategoriesView []string       `json:"categories"`           // untuk Swagger
	CreatedAt      time.Time      `json:"created_at"`
}
