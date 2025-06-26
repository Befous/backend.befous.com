package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/Befous/backend.befous.com/models"
	"github.com/Befous/backend.befous.com/utils"
)

type Relationship struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type MangaResponse struct {
	Data struct {
		Relationships []Relationship `json:"relationships"`
	} `json:"data"`
}

type CoverResponse struct {
	Data struct {
		Attributes struct {
			FileName string `json:"fileName"`
		} `json:"attributes"`
	} `json:"data"`
}

var (
	cacheData      []byte
	cacheTimestamp time.Time
	cacheMutex     sync.Mutex
	cacheDuration  = 10 * time.Minute // Durasi cache aktif
)

func getCoverFileName(mangaID string) (string, error) {
	mangaURL := fmt.Sprintf("https://api.mangadex.org/manga/%s", mangaID)
	resp, err := http.Get(mangaURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var mangaResp MangaResponse
	if err := json.NewDecoder(resp.Body).Decode(&mangaResp); err != nil {
		return "", err
	}

	var coverID string
	for _, rel := range mangaResp.Data.Relationships {
		if rel.Type == "cover_art" {
			coverID = rel.ID
			break
		}
	}

	if coverID == "" {
		return "", fmt.Errorf("cover ID not found")
	}

	coverURL := fmt.Sprintf("https://api.mangadex.org/cover/%s", coverID)
	resp, err = http.Get(coverURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var coverResp CoverResponse
	if err := json.NewDecoder(resp.Body).Decode(&coverResp); err != nil {
		return "", err
	}

	return coverResp.Data.Attributes.FileName, nil
}

func CoverMangadex(w http.ResponseWriter, r *http.Request) {
	mangaID := r.PathValue("id_manga")
	if mangaID == "" {
		http.Error(w, "Missing manga ID", http.StatusBadRequest)
		return
	}

	fileName, err := getCoverFileName(mangaID)
	if err != nil {
		http.Error(w, "Failed to get cover filename", http.StatusInternalServerError)
		return
	}

	imageURL := fmt.Sprintf("https://uploads.mangadex.org/covers/%s/%s.512.jpg", mangaID, fileName)
	resp, err := http.Get(imageURL)
	if err != nil {
		http.Error(w, "Failed to fetch image", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Failed to send image", http.StatusInternalServerError)
	}
}

func IpapiProxyHandler(w http.ResponseWriter, r *http.Request) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Cek apakah cache masih berlaku
	if time.Since(cacheTimestamp) < cacheDuration && cacheData != nil {
		w.Header().Set("Content-Type", "application/json")
		utils.WriteJSONResponse(w, http.StatusOK, models.Pesan{
			Data: json.RawMessage(cacheData),
		})
		return
	}

	// Ambil dari ipapi.co karena cache sudah kadaluarsa
	resp, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		http.Error(w, "Gagal mengambil data dari ipapi.co", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Gagal membaca data response", http.StatusInternalServerError)
		return
	}

	// Validasi apakah body adalah JSON
	var js json.RawMessage
	if err := json.Unmarshal(bodyBytes, &js); err != nil {
		utils.WriteJSONResponse(w, resp.StatusCode, models.Pesan{
			Message: string(bodyBytes), // kirim sebagai string biasa
		})
		return
	}

	// Simpan ke cache
	cacheData = bodyBytes
	cacheTimestamp = time.Now()

	// Kirim ke client
	w.Header().Set("Content-Type", "application/json")
	utils.WriteJSONResponse(w, resp.StatusCode, models.Pesan{
		Data: js,
	})
}
