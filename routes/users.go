package routes

import (
	"net/http"

	"github.com/Befous/api.befous.com/controllers"
)

func Route(router *http.ServeMux) {
	router.HandleFunc("GET /", controllers.RootController)
}
