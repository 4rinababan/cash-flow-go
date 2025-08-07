package models

import (
	"time"

	"github.com/lib/pq"
)

// Transaction mewakili entitas transaksi keuangan
type Transaction struct {
	ID            uint           `json:"id" example:"1" gorm:"primaryKey"`
	Type          string         `json:"type" example:"pengeluaran"`
	Amount        float64        `json:"amount" example:"15000"`
	Description   string         `json:"description" example:"Beli Mie Gacoan"`
	Category      string         `json:"category" example:"makanan"`
	Categories    pq.StringArray `json:"categories" gorm:"type:text[]" swaggertype:"array,string" example:"[\"makanan\",\"jajan\"]"`
	TransactionAt time.Time      `json:"transaction_at" example:"2025-08-07T12:00:00Z"`
	CreatedAt     time.Time      `json:"created_at" example:"2025-08-07T12:00:00Z"`

	// View-only field for Swagger or API response
	CategoriesView []string `json:"categories_view" gorm:"-"`
}
