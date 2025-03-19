package gcf

import (
	"net/http"

	"github.com/Befous/backend.befous.com/middleware"
	"github.com/Befous/backend.befous.com/routes"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("Befous", Befous)
}

func Befous(w http.ResponseWriter, r *http.Request) {
	app := http.NewServeMux()
	routes.Route(app)
	wrappedApp := middleware.Cors(app)
	wrappedApp.ServeHTTP(w, r)
}
