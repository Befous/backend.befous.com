package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Befous/api.befous.com/models"
	"github.com/Befous/api.befous.com/utils"
	"github.com/golang-jwt/jwt/v5"
)

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			utils.WriteJSONResponse(w, http.StatusUnauthorized, models.Pesan{
				Message: "Authorization token required",
			})
			return
		}
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.WriteJSONResponse(w, http.StatusUnauthorized, models.Pesan{
				Message: "Invalid authorization format",
			})
			return
		}
		tokenString = parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return utils.ReadPublicKeyFromEnv("publickey")
		})
		if err != nil || !token.Valid {
			utils.WriteJSONResponse(w, http.StatusUnauthorized, models.Pesan{
				Message: "Invalid or expired token",
			})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["username"] == nil {
			utils.WriteJSONResponse(w, http.StatusUnauthorized, models.Pesan{
				Message: "Invalid token claims",
			})
			return
		}
		userID := claims["username"].(string)
		session := utils.GetSession(utils.SetConnection(), tokenString, userID)
		if time.Now().After(session.Expire_At) {
			utils.WriteJSONResponse(w, http.StatusUnauthorized, models.Pesan{
				Message: "Session invalid",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}
