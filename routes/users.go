package routes

import (
	"net/http"

	"github.com/Befous/backend.befous.com/controllers"
	"github.com/Befous/backend.befous.com/middleware"
)

func Route(router *http.ServeMux) {
	router.HandleFunc("GET /email/{email}", controllers.KirimGmail)
	router.HandleFunc("POST /image", controllers.UploadImage)
	router.HandleFunc("DELETE /image", controllers.DeleteImage)
	router.HandleFunc("GET /", controllers.IpapiProxyHandler)
	router.HandleFunc("POST /login", controllers.Login)
	router.Handle("GET /session", middleware.IsAuthenticated(http.HandlerFunc(controllers.Session)))
	router.Handle("GET /mangadex/cover/{id_manga}", middleware.IsAuthenticated(http.HandlerFunc(controllers.CoverMangadex)))
}
