package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Befous/backend.befous.com/models"
	"github.com/Befous/backend.befous.com/utils"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadImage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, models.Pesan{
			Message: "Gagal parsing form",
		})
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, models.Pesan{
			Message: "Gambar tidak ditemukan",
		})
		return
	}
	defer file.Close()
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUD"),
		os.Getenv("KEY"),
		os.Getenv("SECRET"),
	)
	uploadResult, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{})
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, models.Pesan{
			Message: "Gagal upload gambar",
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, models.Pesan{
		Message: fmt.Sprintf(uploadResult.SecureURL),
	})
}

func DeleteImage(w http.ResponseWriter, r *http.Request) {
	// Pastikan metode yang dipakai POST atau DELETE
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, models.Pesan{
			Message: "Metode harus DELETE atau POST",
		})
		return
	}

	// Misal public_id dikirim via query param atau body JSON
	publicID := r.URL.Query().Get("public_id")
	if publicID == "" {
		utils.WriteJSONResponse(w, http.StatusBadRequest, models.Pesan{
			Message: "public_id harus disertakan",
		})
		return
	}

	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUD"),
		os.Getenv("KEY"),
		os.Getenv("SECRET"),
	)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, models.Pesan{
			Message: "Gagal inisialisasi Cloudinary",
		})
		return
	}

	// Delete image by public_id
	resp, err := cld.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, models.Pesan{
			Message: "Gagal menghapus gambar",
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, models.Pesan{
		Message: fmt.Sprintf("Gambar berhasil dihapus, result: %v", resp.Result),
	})

}
