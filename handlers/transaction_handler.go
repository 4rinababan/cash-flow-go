package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	db "cash-flow-go/database"
	"cash-flow-go/models"

	"github.com/go-chi/chi"
)

// @Summary Tambah transaksi baru
// @Description Menambahkan data transaksi
// @Tags Transactions
// @Accept json
// @Produce json
// @Param transaction body models.Transaction true "Transaksi baru"
// @Success 201 {object} models.Transaction
// @Router /api/transactions [post]
func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var tx models.Transaction

	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(tx.Categories) > 3 {
		http.Error(w, "Max 3 kategori", http.StatusBadRequest)
		return
	}

	// Gunakan waktu sekarang jika CreatedAt tidak dikirim dari frontend
	if tx.CreatedAt.IsZero() {
		tx.CreatedAt = time.Now()
	}

	db.DB.Create(&tx)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tx)
}

// @Summary Ambil daftar transaksi
// @Description Menampilkan semua transaksi dengan filter dan pagination
// @Tags Transactions
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Limit per page (default 10)"
// @Param type query string false "Filter by type (pemasukan/pengeluaran)"
// @Success 200 {array} models.Transaction
// @Router /api/transactions [get]
// GetTransactions handles fetching transactions with optional filters and pagination
func GetTransactions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Pagination
	page := 1
	limit := 10
	if val := query.Get("page"); val != "" {
		if p, err := strconv.Atoi(val); err == nil && p > 0 {
			page = p
		}
	}
	if val := query.Get("limit"); val != "" {
		if l, err := strconv.Atoi(val); err == nil && l > 0 {
			limit = l
		}
	}

	// Filters
	txType := query.Get("type")
	category := query.Get("category")
	startDate := query.Get("start_date")
	endDate := query.Get("end_date")
	note := query.Get("note")
	minAmount := query.Get("min_amount")
	maxAmount := query.Get("max_amount")

	// Base builder
	queryBuilder := db.DB.Model(&models.Transaction{})
	if txType == "pemasukan" || txType == "pengeluaran" {
		queryBuilder = queryBuilder.Where("type = ?", txType)
	}
	if category != "" {
		queryBuilder = queryBuilder.Where("category = ?", category)
	}
	if startDate != "" {
		queryBuilder = queryBuilder.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		queryBuilder = queryBuilder.Where("created_at <= ?", endDate)
	}
	if note != "" {
		queryBuilder = queryBuilder.Where("note LIKE ?", "%"+note+"%")
	}
	if minAmount != "" {
		if min, err := strconv.ParseFloat(minAmount, 64); err == nil {
			queryBuilder = queryBuilder.Where("amount >= ?", min)
		}
	}
	if maxAmount != "" {
		if max, err := strconv.ParseFloat(maxAmount, 64); err == nil {
			queryBuilder = queryBuilder.Where("amount <= ?", max)
		}
	}

	// Hitung total count dan total amount dulu
	var totalCount int64
	var totalAmount float64
	queryBuilder.Count(&totalCount)
	queryBuilder.Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount)

	// Ambil data transaksi dengan limit dan offset
	var txs []models.Transaction
	offset := (page - 1) * limit
	if err := queryBuilder.Order("created_at desc").Offset(offset).Limit(limit).Find(&txs).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response dengan data + meta
	response := map[string]interface{}{
		"data":         txs,
		"total_count":  totalCount,
		"total_amount": totalAmount,
		"page":         page,
		"limit":        limit,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Hapus transaksi
// @Description Menghapus transaksi berdasarkan ID
// @Tags Transactions
// @Param id path int true "Transaction ID"
// @Success 200 {object} map[string]string
// @Router /api/transactions/{id} [delete]
func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL
	id := chi.URLParam(r, "id")

	var tx models.Transaction
	if err := db.DB.First(&tx, id).Error; err != nil {
		http.Error(w, "Transaksi tidak ditemukan", http.StatusNotFound)
		return
	}

	if err := db.DB.Delete(&tx).Error; err != nil {
		http.Error(w, "Gagal menghapus transaksi", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaksi berhasil dihapus"})
}
