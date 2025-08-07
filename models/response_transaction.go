package models

type TransactionResponse struct {
	ID            uint    `json:"id"`
	Type          string  `json:"type"`
	Category      string  `json:"category"`
	Description   string  `json:"description"`
	Amount        float64 `json:"amount"`
	TransactionAt string  `json:"transaction_at"`
	CreatedAt     string  `json:"created_at"` // string untuk tampil WIB
}
