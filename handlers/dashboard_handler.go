package handlers

import (
	db "cash-flow-go/database"
	"cash-flow-go/models"
	"encoding/json"
	"net/http"
	"time"
)

type MonthlyBalance struct {
	Month   string `json:"month"`
	Year    int    `json:"year"`
	Income  int64  `json:"income"`
	Expense int64  `json:"expense"`
	Saldo   int64  `json:"saldo"`
}

// @Summary Dashboard utama
// @Description Menampilkan ringkasan transaksi
// @Tags Dashboard
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/dashboard [get]
func GetDashboard(w http.ResponseWriter, r *http.Request) {
	var pemasukan int64
	var pengeluaran int64

	// Pastikan NULL dari SUM dihindari dengan COALESCE
	db.DB.Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("type = ?", "pemasukan").
		Scan(&pemasukan)

	db.DB.Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("type = ?", "pengeluaran").
		Scan(&pengeluaran)

	// Ambil saldo 3 bulan terakhir
	var monthly []MonthlyBalance
	now := time.Now()

	for i := 0; i < 3; i++ {
		target := now.AddDate(0, -i, 0)
		var income int64
		var expense int64

		db.DB.Model(&models.Transaction{}).
			Select("COALESCE(SUM(amount), 0)").
			Where("type = ? AND EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", "pemasukan", target.Month(), target.Year()).
			Scan(&income)

		db.DB.Model(&models.Transaction{}).
			Select("COALESCE(SUM(amount), 0)").
			Where("type = ? AND EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", "pengeluaran", target.Month(), target.Year()).
			Scan(&expense)

		monthly = append(monthly, MonthlyBalance{
			Month:   target.Month().String(),
			Year:    target.Year(),
			Income:  income,
			Expense: expense,
			Saldo:   income - expense,
		})

	}

	response := map[string]interface{}{
		"total_balance":   pemasukan - pengeluaran,
		"total_income":    pemasukan,
		"total_expense":   pengeluaran,
		"monthly_balance": monthly,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Grafik batang pengeluaran
// @Description Menampilkan grafik batang pengeluaran per kategori
// @Tags Dashboard
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/dashboard/bar [get]
func GetBarChart(w http.ResponseWriter, r *http.Request) {
	type Result struct {
		Category string
		Total    int
	}

	var results []Result
	db.DB.Raw(`
		SELECT unnest(categories) AS category, SUM(amount) AS total
		FROM transactions
		WHERE type = 'pengeluaran'
		GROUP BY category
	`).Scan(&results)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// @Summary Grafik donat pemasukan
// @Description Menampilkan grafik donat pemasukan per kategori
// @Tags Dashboard
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/dashboard/donut [get]
func GetDonutChart(w http.ResponseWriter, r *http.Request) {
	type Result struct {
		Category string
		Total    int
	}

	var results []Result
	db.DB.Raw(`
        SELECT json_each.value AS category, SUM(amount) AS total 
        FROM transactions, json_each(transactions.categories)
        WHERE type = 'pemasukan'
        GROUP BY category
    `).Scan(&results)

	json.NewEncoder(w).Encode(results)
}
