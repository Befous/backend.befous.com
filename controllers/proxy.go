package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const imageURL = "https://mangadex.org/covers/a08cd14f-68d4-4942-9305-2fe9d7d0339e/d59ecfc6-2391-44f6-9d45-361565cc343e.jpg.512.jpg"

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
