package controllers

import (
	"net/http"

	"github.com/Befous/backend.befous.com/models"
	"github.com/Befous/backend.befous.com/utils"
)

func Root(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSONResponse(w, http.StatusOK, models.Pesan{
		Message: "Hello World",
	})
}

func Session(w http.ResponseWriter, r *http.Request) {
	var session models.Users
	session = utils.DecodeJWT(r)

	utils.WriteJSONResponse(w, http.StatusOK, models.Pesan{
		Message: "Berikut data session anda",
		Data:    session,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	mconn := utils.SetConnection()
	var user models.Users

	utils.ParseBody(w, r, &user)

	if !utils.UsernameExists(mconn, user) {
		utils.WriteJSONResponse(w, http.StatusBadRequest, models.Pesan{
			Message: "Akun tidak ditemukan",
		})
		return
	}

	if !utils.IsPasswordValid(mconn, user) {
		utils.WriteJSONResponse(w, http.StatusBadRequest, models.Pesan{
			Message: "Password salah",
		})
		return
	}

	userAgent := r.UserAgent()
	token, err := utils.SignedJWT(mconn, user, userAgent)

	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, models.Pesan{
			Message: err.Error(),
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, models.Pesan{
		Message: "Berhasil login",
		Token:   token,
	})
}
