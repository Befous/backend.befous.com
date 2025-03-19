package routes

import (
	"net/http"

	"github.com/Befous/backend.befous.com/controllers"
	"github.com/Befous/backend.befous.com/middleware"
)

func Route(router *http.ServeMux) {
	router.HandleFunc("GET /", controllers.Root)
	router.HandleFunc("POST /login", controllers.Login)
	router.Handle("GET /session", middleware.IsAuthenticated(http.HandlerFunc(controllers.Session)))
}
