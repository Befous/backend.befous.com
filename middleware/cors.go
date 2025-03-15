package middleware

import (
	"net/http"
	"strings"
)

var headers = []string{
	"Content-Type",
	"Content-Disposition",
	"Authorization",
	"Accept",
	"Origin",
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Authorization")
		w.Header().Set("Access-Control-Max-Age", "7200")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
