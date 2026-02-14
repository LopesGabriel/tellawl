package controllers

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
)

type contextKey string

const (
	memberContextKey = contextKey("member")
)

func (h *APIHandler) jwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			h.logger.Error(r.Context(), "missing authorization header", slog.String("error", "Authorization Header not available"))
			WriteError(w, http.StatusUnauthorized, map[string]any{
				"message": "Missing Authorization header",
			})
			return
		}

		// Expect header in form: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			h.logger.Error(r.Context(), "Invalid Authorization header value", slog.String("error", "expected token in 'Bearer <token>' format"))
			WriteError(w, http.StatusUnauthorized, map[string]any{
				"message": "Invalid Authorization header format",
			})
			return
		}

		tokenString := parts[1]
		member, err := h.usecases.AuthenticateUser(r.Context(), usecases.AuthenticateUserUseCaseInput{
			Token: tokenString,
		})
		if err != nil {
			h.logger.Error(r.Context(), "Invalid token", slog.String("error", err.Error()))
			WriteError(w, http.StatusUnauthorized, map[string]any{
				"message": "Invalid token",
			})
			return
		}

		if member != nil {
			ctx := context.WithValue(r.Context(), memberContextKey, member)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		h.logger.Error(r.Context(), "Invalid token claims")
		WriteError(w, http.StatusUnauthorized, map[string]any{
			"message": "Invalid token claims",
		})
	})
}
