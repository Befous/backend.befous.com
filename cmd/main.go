package main

import (
	"log"
	"net/http"

	"github.com/Befous/api.befous.com/middleware"
	"github.com/Befous/api.befous.com/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}
	app := http.NewServeMux()
	routes.Route(app)
	server := http.Server{
		Addr:    ":3000",
		Handler: middleware.Cors(app),
	}
	server.ListenAndServe()
}
