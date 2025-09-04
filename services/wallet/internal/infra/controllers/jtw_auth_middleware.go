package controllers

import (
	"context"
	"net/http"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	userContextKey = contextKey("user")
	jwtSecret      = "your-secret-key" // TODO: move to env var
)

func jwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			WriteError(w, http.StatusUnauthorized, map[string]any{
				"message": "Missing Authorization header",
			})
			return
		}

		// Expect header in form: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			WriteError(w, http.StatusUnauthorized, map[string]any{
				"message": "Invalid Authorization header format",
			})
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			WriteError(w, http.StatusUnauthorized, map[string]any{
				"message": "Invalid token",
			})
			return
		}

		// Optional: Store claims in context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx := context.WithValue(r.Context(), userContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		WriteError(w, http.StatusUnauthorized, map[string]any{
			"message": "Invalid token claims",
		})
	})
}
