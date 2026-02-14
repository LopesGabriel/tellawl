package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/repository"
	usecases "github.com/lopesgabriel/tellawl/services/member-service/internal/use_cases"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type apiHandler struct {
	repositories *repository.Repositories
	usecases     *usecases.UseCases
	version      string
	tracer       trace.Tracer
	logger       *logger.AppLogger
}

type NewAPiArgs struct {
	Repositories *repository.Repositories
	Usecases     *usecases.UseCases
	Version      string
	Tracer       trace.Tracer
	Logger       *logger.AppLogger
}

func NewApiHandler(args NewAPiArgs) *apiHandler {
	return &apiHandler{
		repositories: args.Repositories,
		usecases:     args.Usecases,
		version:      args.Version,
		tracer:       args.Tracer,
		logger:       args.Logger,
	}
}

func (handler *apiHandler) Listen(ctx context.Context, port int) error {
	router := handler.setupRoutes()
	handler.logger.Info(ctx, "Starting api server", slog.Int("port", port))
	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

func (handler *apiHandler) requestInterceptorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Server-Version", handler.version)
		handler.logger.Debug(r.Context(), "processing new request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)
		next.ServeHTTP(w, r)
	})
}

func (handler *apiHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	_, span := handler.tracer.Start(r.Context(), "HandleHealthCheck")
	defer span.End()

	result, _ := json.Marshal(map[string]string{
		"api": "OK",
	})

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
	span.SetStatus(codes.Ok, "OK")
}
