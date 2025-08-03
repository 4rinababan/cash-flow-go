package handlers

import (
	db "cash-flow-go/database"
	"cash-flow-go/models"
	"encoding/json"
	"net/http"
	"time"
)

type MonthlyBalance struct {
	Month     string `json:"month"`
	Year      int    `json:"year"`
	Income    int64  `json:"income"`
	Expense   int64  `json:"expense"`
	PrevSaldo int64  `json:"prev_saldo"`
	Saldo     int64  `json:"saldo"`
}

// @Summary Dashboard utama
// @Description Menampilkan ringkasan transaksi
// @Tags Dashboard
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/dashboard [get]
func GetDashboard(w http.ResponseWriter, r *http.Request) {
	var pemasukan, pengeluaran int64

	// Total pemasukan dan pengeluaran
	db.DB.Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("type = ?", "pemasukan").
		Scan(&pemasukan)

	db.DB.Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("type = ?", "pengeluaran").
		Scan(&pengeluaran)

	// Ambil semua bulan dan tahun unik dari transaksi
	type MonthYear struct {
		Month int
		Year  int
	}
	var monthYears []MonthYear
	db.DB.Raw(`
		SELECT DISTINCT EXTRACT(MONTH FROM created_at) AS month, EXTRACT(YEAR FROM created_at) AS year
		FROM transactions
		ORDER BY year, month
	`).Scan(&monthYears)

	var monthly []MonthlyBalance
	var prevSaldo int64 = 0

	for _, my := range monthYears {
		var income, expense int64

		db.DB.Model(&models.Transaction{}).
			Select("COALESCE(SUM(amount), 0)").
			Where("type = ? AND EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", "pemasukan", my.Month, my.Year).
			Scan(&income)

		db.DB.Model(&models.Transaction{}).
			Select("COALESCE(SUM(amount), 0)").
			Where("type = ? AND EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", "pengeluaran", my.Month, my.Year).
			Scan(&expense)

		saldo := prevSaldo + income - expense

		monthly = append(monthly, MonthlyBalance{
			Month:     time.Month(my.Month).String(),
			Year:      my.Year,
			Income:    income,
			Expense:   expense,
			PrevSaldo: prevSaldo,
			Saldo:     saldo,
		})

		prevSaldo = saldo
	}

	// Ambil 3 bulan terakhir
	last3 := monthly
	if len(monthly) > 3 {
		last3 = monthly[len(monthly)-3:]
	}

	response := map[string]interface{}{
		"total_balance":   pemasukan - pengeluaran,
		"total_income":    pemasukan,
		"total_expense":   pengeluaran,
		"monthly_balance": last3,
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
