package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	usecases "github.com/lopesgabriel/tellawl/services/member-service/internal/use_cases"
	"go.opentelemetry.io/otel/codes"
)

func (h *apiHandler) HandleMe(w http.ResponseWriter, r *http.Request) error {
	ctx, span := h.tracer.Start(r.Context(), "HandleMe")
	defer span.End()
	w.Header().Add("Trace-ID", span.SpanContext().TraceID().String())

	token := r.Header.Get("Authorization")
	if token == "" {
		span.SetStatus(codes.Error, "Missing token")
		return NewUnauthorizedError(ctx, "Missing token")
	}
	token = strings.Replace(token, "Bearer ", "", 1)

	member, err := h.usecases.GetMemberFromToken(ctx, usecases.GetMemberFromTokenUseCaseInput{
		Token: token,
	})
	if err != nil {
		if errors.Is(err, usecases.ErrInvalidCredentials) {
			span.SetStatus(codes.Error, "Invalid credentials")
			return NewUnauthorizedError(ctx, "Invalid credentials")
		}
		if errors.Is(err, usecases.ErrMemberNotFound) {
			span.SetStatus(codes.Error, "Member not found")
			return NewNotFoundError(ctx, "Member not found")
		}
		span.SetStatus(codes.Error, "Could not retrieve member")
		return NewInternalError(ctx, "Could not retrieve member", err)
	}

	span.SetStatus(codes.Ok, "OK")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(toMemberAPIResponse(member))
	return nil
}
