package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Befous/backend.befous.com/models"
	"github.com/Befous/backend.befous.com/utils"
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
			return utils.ReadPublicKeyFromEnv("public_key")
		})
		if err != nil || !token.Valid {
			utils.WriteJSONResponse(w, http.StatusUnauthorized, models.Pesan{
				Message: "Invalid token",
			})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["sub"] == nil {
			utils.WriteJSONResponse(w, http.StatusUnauthorized, models.Pesan{
				Message: "Invalid token claims",
			})
			return
		}
		exp, ok := claims["exp"].(float64)
		if !ok {
			utils.WriteJSONResponse(w, http.StatusUnauthorized, models.Pesan{
				Message: "Invalid token expiration",
			})
			return
		}
		if time.Now().Unix() > int64(exp) {
			utils.WriteJSONResponse(w, http.StatusUnauthorized, models.Pesan{
				Message: "Token has expired",
			})
			return
		}
		userID := claims["sub"].(string)
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
