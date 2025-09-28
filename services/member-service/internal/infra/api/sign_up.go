package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/lopesgabriel/tellawl/packages/logger"
	usecases "github.com/lopesgabriel/tellawl/services/member-service/internal/use_cases"
	"go.opentelemetry.io/otel/codes"
)

type signUpRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (handler *apiHandler) HandleSignUp(w http.ResponseWriter, r *http.Request) error {
	ctx, span := handler.tracer.Start(r.Context(), "HandleSignUp")
	defer span.End()
	w.Header().Add("Trace-ID", span.SpanContext().TraceID().String())

	var requestData signUpRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	defer r.Body.Close()
	if err != nil {
		span.SetStatus(codes.Error, "Could not parse the request body")
		return NewBadRequestError(ctx, "Could not parse the request body, are you sending a JSON?", err)
	}

	logger.Debug(
		ctx,
		"Attempting to Sign Up new member",
		slog.String("email", fmt.Sprintf("%s%s", requestData.Email[:4], strings.Repeat("*", len(requestData.Email)-4))),
	)
	member, err := handler.usecases.EmailPasswordSignUp(ctx, usecases.EmailPasswordSignUpUseCaseInput{
		Email:     requestData.Email,
		FirstName: requestData.FirstName,
		LastName:  requestData.LastName,
		Password:  requestData.Password,
	})

	if err != nil {
		span.SetStatus(codes.Error, "Could not sign up the member")
		return NewInternalError(ctx, "Could not sign up the member", err)
	}

	logger.Debug(
		ctx,
		"Member signed up successfully",
		slog.String("email", fmt.Sprintf("%s%s", requestData.Email[:4], strings.Repeat("*", len(requestData.Email)-4))),
		slog.String("memberId", member.Id),
	)
	span.SetStatus(codes.Ok, "OK")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(member)
	return nil
}
