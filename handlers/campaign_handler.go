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

type Campaign struct {
	ID       int    `json:"id"`
	ImageURL string `json:"image_url"`
	IsActive bool   `json:"-"`
}

var (
	campaigns  = []Campaign{}
	idCounter  = 1
	campaignMu sync.Mutex
)

// CreateCampaign godoc
// @Summary Upload campaign baru (dengan gambar)
// @Description Mengunggah campaign dan menjadikannya aktif (hanya 1 yang aktif).
// @Tags Campaign
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Gambar campaign (jpg/png)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/campaigns [post]
func CreateCampaign(w http.ResponseWriter, r *http.Request) {
	campaignMu.Lock()
	defer campaignMu.Unlock()

	err := r.ParseMultipartForm(10 << 20) // max 10MB
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

	// Nonaktifkan campaign yang lain
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
	}
	idCounter++
	campaigns = append(campaigns, newCampaign)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Campaign created successfully",
	})
}

// GetActiveCampaign godoc
// @Summary Ambil campaign yang aktif saat ini
// @Description Mendapatkan campaign yang sedang aktif untuk ditampilkan sebagai popup.
// @Tags Campaign
// @Produce json
// @Success 200 {object} handlers.Campaign
// @Failure 404 {object} map[string]string
// @Router /api/campaigns/active [get]
func GetActiveCampaign(w http.ResponseWriter, r *http.Request) {
	campaignMu.Lock()
	defer campaignMu.Unlock()

	for _, c := range campaigns {
		if c.IsActive {
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
