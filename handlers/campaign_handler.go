package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	campaigns  = []Campaign{}
	idCounter  = 1
	campaignMu sync.Mutex
)

type Campaign struct {
	ID       int       `json:"id"`
	ImageURL string    `json:"image_url"`
	IsActive bool      `json:"-"`
	StartAt  time.Time `json:"start_at"`
	EndAt    time.Time `json:"end_at"`
}

// CreateCampaign godoc
// @Summary Upload campaign baru dengan waktu aktif
// @Description Mengunggah campaign (dengan waktu mulai & akhir) dan menjadikannya aktif.
// @Tags Campaign
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Gambar campaign"
// @Param start_at formData string true "Waktu mulai campaign (format: 2006-01-02T15:04:05)"
// @Param end_at formData string true "Waktu akhir campaign (format: 2006-01-02T15:04:05)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/campaigns [post]
func CreateCampaign(w http.ResponseWriter, r *http.Request) {
	campaignMu.Lock()
	defer campaignMu.Unlock()

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Image is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	startStr := r.FormValue("start_at")
	endStr := r.FormValue("end_at")

	if startStr == "" || endStr == "" {
		http.Error(w, "Start and end time are required", http.StatusBadRequest)
		return
	}

	startAt, err := time.Parse("2006-01-02T15:04:05", startStr)
	if err != nil {
		http.Error(w, "Invalid start_at format (use YYYY-MM-DDTHH:MM:SS)", http.StatusBadRequest)
		return
	}

	endAt, err := time.Parse("2006-01-02T15:04:05", endStr)
	if err != nil {
		http.Error(w, "Invalid end_at format (use YYYY-MM-DDTHH:MM:SS)", http.StatusBadRequest)
		return
	}

	if endAt.Before(startAt) {
		http.Error(w, "end_at harus setelah start_at", http.StatusBadRequest)
		return
	}

	os.MkdirAll("./uploads", os.ModePerm)
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(handler.Filename))
	filepath := "./uploads/" + filename

	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = dst.ReadFrom(file)
	if err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	// Nonaktifkan campaign lain
	for i := range campaigns {
		campaigns[i].IsActive = false
	}

	host := r.Host
	if r.TLS != nil {
		host = "https://" + host
	} else {
		host = "http://" + host
	}

	newCampaign := Campaign{
		ID:       idCounter,
		ImageURL: host + "/uploads/" + filename,
		IsActive: true,
		StartAt:  startAt,
		EndAt:    endAt,
	}
	idCounter++
	campaigns = append(campaigns, newCampaign)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Campaign created successfully",
	})
}

// GetActiveCampaign godoc
// @Summary Ambil campaign yang aktif dan dalam rentang waktu
// @Description Mendapatkan campaign yang sedang aktif berdasarkan waktu saat ini
// @Tags Campaign
// @Produce json
// @Success 200 {object} handlers.Campaign
// @Failure 404 {object} map[string]string
// @Router /api/campaigns/active [get]
func GetActiveCampaign(w http.ResponseWriter, r *http.Request) {
	campaignMu.Lock()
	defer campaignMu.Unlock()

	now := time.Now()
	for _, c := range campaigns {
		if c.IsActive && now.After(c.StartAt) && now.Before(c.EndAt) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "No active campaign",
	})
}
