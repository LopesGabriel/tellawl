package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	usecases "github.com/lopesgabriel/tellawl/services/member-service/internal/use_cases"
	"go.opentelemetry.io/otel/codes"
)

func (h *apiHandler) HandleRefreshToken(w http.ResponseWriter, r *http.Request) error {
	ctx, span := h.tracer.Start(r.Context(), "HandleRefreshToken")
	defer span.End()
	w.Header().Add("Trace-ID", span.SpanContext().TraceID().String())

	refreshToken := r.Header.Get("Authorization")
	if refreshToken == "" {
		span.SetStatus(codes.Error, "Authorization header is missing")
		return NewBadRequestError(ctx, "Authorization header is missing", nil)
	}
	refreshToken = strings.ReplaceAll(refreshToken, "Bearer ", "")

	output, err := h.usecases.RefreshToken(ctx, usecases.RefreshTokenUseCaseInput{
		RefreshToken: refreshToken,
	})
	if err != nil {
		if errors.Is(err, usecases.ErrInvalidCredentials) {
			span.SetStatus(codes.Error, "Invalid credentials")
			return NewBadRequestError(ctx, "Invalid credentials", err)
		}

		if errors.Is(err, usecases.ErrMemberNotFound) {
			span.SetStatus(codes.Error, "Member not found")
			return NewConflictError(ctx, "Member not found")
		}

		span.SetStatus(codes.Error, "Could not sign in the member")
		return NewInternalError(ctx, "Could not sign in the member", err)
	}

	span.SetStatus(codes.Ok, "OK")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
	return nil
}
