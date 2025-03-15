package main

import (
	"net/http"

	"github.com/Befous/api.befous.com/middleware"
	"github.com/Befous/api.befous.com/routes"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")
	app := http.NewServeMux()
	routes.Route(app)
	server := http.Server{
		Addr:    ":3000",
		Handler: middleware.Cors(app),
	}
	server.ListenAndServe()
}
