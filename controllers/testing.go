package controllers

import (
	"net/http"

	"github.com/Befous/backend.befous.com/models"
	"github.com/Befous/backend.befous.com/utils"
)

func TestInsert(w http.ResponseWriter, r *http.Request) {
	session := models.Session{
		Username: utils.RandomString(10),
	}
	utils.InsertTesting(utils.SetConnection(), session)

	utils.WriteJSONResponse(w, http.StatusOK, models.Pesan{
		Message: "Berikut data session anda",
		Data:    session,
	})
}
