package controllers

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/lopesgabriel/tellawl/packages/logger"
)

type contextKey string

const (
	userContextKey = contextKey("user")
)

func jwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Error(r.Context(), "missing authorization header", slog.String("error", "Authorization Header not available"))
			WriteError(w, http.StatusUnauthorized, map[string]any{
				"message": "Missing Authorization header",
			})
			return
		}

		// Expect header in form: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			logger.Error(r.Context(), "Invalid Authorization header value", slog.String("error", "expected token in 'Bearer <token>' format"))
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
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			logger.Error(r.Context(), "Invalid token", slog.String("error", err.Error()))
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

		logger.Error(r.Context(), "Invalid token claims")
		WriteError(w, http.StatusUnauthorized, map[string]any{
			"message": "Invalid token claims",
		})
	})
}
