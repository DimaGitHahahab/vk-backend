package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := r.Cookie("Authorization")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenStr.Value, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if time.Now().Unix() > int64(claims["exp"].(float64)) {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				ctx := context.WithValue(r.Context(), "user_id", claims["sub"])
				ctx = context.WithValue(ctx, "user_role", claims["role"])
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	})
}
