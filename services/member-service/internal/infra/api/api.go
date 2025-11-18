package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/lopesgabriel/tellawl/packages/logger"
	usecases "github.com/lopesgabriel/tellawl/services/member-service/internal/use_cases"
	"go.opentelemetry.io/otel/trace"
)

type apiHandler struct {
	usecases *usecases.UseCases
	version  string
	tracer   trace.Tracer
	logger   *logger.AppLogger
}

func NewApiHandler(usecases *usecases.UseCases, version string, tracer trace.Tracer, logger *logger.AppLogger) *apiHandler {
	return &apiHandler{
		usecases: usecases,
		version:  version,
		tracer:   tracer,
		logger:   logger,
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
