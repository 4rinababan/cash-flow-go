package models

import (
	"time"

	"github.com/lib/pq"
)

type Transaction struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Type        string         `json:"type"`
	Amount      float64        `json:"amount"`
	Description string         `json:"description"`
	Category    string         `json:"category"`
	Categories  pq.StringArray `json:"categories" gorm:"type:text[]"`
	CreatedAt   time.Time      `json:"created_at"`

	// View-only field for Swagger or API response
	CategoriesView []string `json:"categories_view" gorm:"-"`
}
