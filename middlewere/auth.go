package middlewere

import (
	"context"
	"net/http"
	"strings"

	"go-jwt/helpers"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helpers.Response(w, 401, "Unauthorized", nil)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			helpers.Response(w, 401, "Invalid token format", nil)
			return
		}

		accessToken := parts[1] 
		user, err := helpers.VerifyToken(accessToken)
		if err != nil {
			helpers.Response(w, 401, err.Error(), nil)
			return
		}

		ctx := context.WithValue(r.Context(), "userinfo", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
