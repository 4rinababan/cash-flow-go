package models

import (
	"time"

	"github.com/lib/pq"
)

type Transaction struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Type        string         `json:"type"`
	Categories  pq.StringArray `json:"categories" gorm:"type:text[]"`
	Description string         `json:"description"`
	Amount      int            `json:"amount"`
	CreatedAt   time.Time      `json:"created_at"`
}
