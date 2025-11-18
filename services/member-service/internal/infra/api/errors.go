package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	usecases "github.com/lopesgabriel/tellawl/services/member-service/internal/use_cases"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// AppError represents application-level errors with HTTP status codes
type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Err     error  `json:"-"`
	Context context.Context
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// Common error constructors
func NewBadRequestError(ctx context.Context, message string, err error) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: message, Err: err, Context: ctx}
}

func NewNotFoundError(ctx context.Context, message string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: message, Context: ctx}
}

func NewConflictError(ctx context.Context, message string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: message, Context: ctx}
}

func NewInternalError(ctx context.Context, message string, err error) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: message, Err: err, Context: ctx}
}

// HandlerFunc wraps handlers that return errors
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// ErrorMiddleware converts HandlerFunc to http.HandlerFunc with error handling
func (h *apiHandler) ErrorMiddleware(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := handler(w, r); err != nil {
			var appErr *AppError
			if errors.As(err, &appErr) {
				// Handle known application errors
				span := trace.SpanFromContext(appErr.Context)

				if errors.Is(appErr.Err, usecases.ErrMemberAlreadyExists) {
					h.logger.Error(appErr.Context, "Member already exists", slog.String("error", err.Error()))
					writeError(w, http.StatusConflict, map[string]any{
						"message": "Member already exists",
					})
					span.SetStatus(codes.Error, "Member already exists")
					return
				}

				if errors.Is(appErr.Err, usecases.ErrMemberNotFound) {
					h.logger.Error(appErr.Context, "Member not found", slog.String("error", err.Error()))
					writeError(w, http.StatusNotFound, map[string]any{
						"message": "Member not found",
					})
					span.SetStatus(codes.Error, "Member not found")
					return
				}

				h.logger.Error(appErr.Context, appErr.Message, slog.String("error", appErr.Error()))
				writeError(w, appErr.Code, map[string]any{
					"message": appErr.Message,
				})
				span.SetStatus(codes.Error, appErr.Message)
				return
			}

			span := trace.SpanFromContext(ctx)
			// Handle unknown errors
			h.logger.Error(ctx, "Internal server error", slog.String("error", err.Error()))
			writeError(w, http.StatusInternalServerError, map[string]any{
				"message": "Internal server error",
			})
			span.SetStatus(codes.Error, "Internal Server Error")
		}
	}
}

func writeError(w http.ResponseWriter, statusCode int, data map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
