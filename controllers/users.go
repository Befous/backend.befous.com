package controllers

import (
	"net/http"

	// "github.com/Befous/api.befous.com/helpers"
	"github.com/Befous/api.befous.com/models"
	"github.com/Befous/api.befous.com/utils"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

func RootController(w http.ResponseWriter, r *http.Request) {
	// mconn := utils.SetConnection()
	// users, err := helpers.GetAllDoc[[]models.Users](mconn, "users", primitive.M{})
	// if err != nil {
	// 	utils.WriteJSONResponse(w, http.StatusOK, models.Pesan{
	// 		Message: "Error:" + err.Error(),
	// 	})
	// 	return
	// }
	utils.WriteJSONResponse(w, http.StatusOK, models.Pesan{
		Message: "Hello World",
	})
}
