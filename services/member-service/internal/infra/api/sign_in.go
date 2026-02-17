package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	usecases "github.com/lopesgabriel/tellawl/services/member-service/internal/use_cases"
	"go.opentelemetry.io/otel/codes"
)

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *apiHandler) HandleSignIn(w http.ResponseWriter, r *http.Request) error {
	ctx, span := h.tracer.Start(r.Context(), "HandleSignIn")
	defer span.End()
	w.Header().Add("Trace-ID", span.SpanContext().TraceID().String())

	var requestData signInRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	defer r.Body.Close()
	if err != nil {
		span.SetStatus(codes.Error, "Could not parse the request body")
		return NewBadRequestError(ctx, "Could not parse the request body, are you sending a JSON?", err)
	}

	h.logger.Debug(
		ctx,
		"Attempting to Sign In member",
		slog.String("email", fmt.Sprintf("%s%s", requestData.Email[:4], strings.Repeat("*", len(requestData.Email)-4))),
	)
	output, err := h.usecases.SignIn(ctx, usecases.SignInUseCaseInput{
		Email:    requestData.Email,
		Password: requestData.Password,
	})
	if err != nil {
		if errors.Is(err, usecases.ErrInvalidCredentials) {
			span.SetStatus(codes.Error, "Invalid credentials")
			return NewUnauthorizedError(ctx, "Invalid credentials")
		}

		span.SetStatus(codes.Error, "Could not sign in the member")
		return NewInternalError(ctx, "Could not sign in the member", err)
	}

	h.logger.Debug(
		ctx,
		"Member signed in successfully",
		slog.String("email", fmt.Sprintf("%s%s", requestData.Email[:4], strings.Repeat("*", len(requestData.Email)-4))),
	)
	span.SetStatus(codes.Ok, "OK")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
	return nil
}
