package main

import (
	"context"
	"log/slog"

	"github.com/lopesgabriel/tellawl/packages/logger"
)

func main() {
	ctx := context.Background()
	shutdown, err := initTelemetry(ctx)
	if err != nil {
		panic(err)
	}
	defer shutdown()

	logger.Info(ctx, "Starting member service", slog.String("example", "Example Value"))
	logger.Debug(ctx, "Debugging member service", slog.String("example", "Example Value"))
	logger.Warn(ctx, "Warning member service", slog.String("example", "Example Value"))
	logger.Error(ctx, "Error member service", slog.String("example", "Example Value"))
	<-ctx.Done()
}

func initTelemetry(ctx context.Context) (func() error, error) {
	logProvider, err := logger.Init(ctx, logger.InitLoggerArgs{
		CollectorURL:     "localhost:4317",
		ServiceName:      "member-service",
		ServiceNamespace: "tellawl",
		ServiceVersion:   "v1.0.0",
	})
	if err != nil {
		return nil, err
	}

	return func() error {
		return logProvider.Shutdown(ctx)
	}, nil
}
