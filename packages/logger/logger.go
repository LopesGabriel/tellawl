package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	sdkLog "go.opentelemetry.io/otel/sdk/log"
)

type InitLoggerArgs struct {
	CollectorURL     string
	ServiceName      string
	ServiceNamespace string
	ServiceVersion   string
	Level            slog.Level
	LoggerProvider   log.LoggerProvider
}

type AppLogger struct {
	logger         *slog.Logger
	loggerProvider log.LoggerProvider
	Level          slog.Level
}

var appLogger *AppLogger

func Init(ctx context.Context, args InitLoggerArgs) (*AppLogger, error) {
	if args.ServiceNamespace == "" {
		args.ServiceNamespace = "tellawl"
	}

	res, err := newResource(args)
	if err != nil {
		return nil, err
	}

	var logProvider log.LoggerProvider
	if args.LoggerProvider != nil {
		logProvider = args.LoggerProvider
	} else {
		logProvider, err = newProvider(ctx, res, args)
		if err != nil {
			return nil, err
		}
	}

	global.SetLoggerProvider(logProvider)

	logger := otelslog.NewLogger(
		fmt.Sprintf("%s/%s", args.ServiceNamespace, args.ServiceName),
		otelslog.WithLoggerProvider(logProvider),
	)

	appLogger = &AppLogger{
		logger:         logger,
		loggerProvider: logProvider,
		Level:          args.Level,
	}
	appLogger.SetLevel(ctx, args.Level)

	return appLogger, nil
}

func GetLogger() (*AppLogger, error) {
	if appLogger == nil {
		return nil, fmt.Errorf("logger not initialized")
	}

	return appLogger, nil
}

func (l *AppLogger) SetLevel(ctx context.Context, level slog.Level) {
	l.Level = level
	l.logger.Enabled(ctx, level)
}

func (l *AppLogger) Info(ctx context.Context, message string, args ...any) {
	l.logger.InfoContext(ctx, message, args...)
	fmt.Printf("\033[34m[INFO]\033[0m %s %s: %v\n", time.Now().UTC().Format(time.RFC3339), message, args)
}

func (l *AppLogger) Error(ctx context.Context, message string, args ...any) {
	args = append(args, slog.Bool("fatal", false))
	l.logger.ErrorContext(ctx, message, args...)
	fmt.Printf("\033[31m[ERROR]\033[0m %s %s: %v\n", time.Now().UTC().Format(time.RFC3339), message, args)
}

func (l *AppLogger) Debug(ctx context.Context, message string, args ...any) {
	l.logger.DebugContext(ctx, message, args...)
	fmt.Printf("\033[36m[DEBUG]\033[0m %s %s: %v\n", time.Now().UTC().Format(time.RFC3339), message, args)
}

func (l *AppLogger) Warn(ctx context.Context, message string, args ...any) {
	l.logger.WarnContext(ctx, message, args...)
	fmt.Printf("\033[33m[WARN]\033[0m %s %s: %v\n", time.Now().UTC().Format(time.RFC3339), message, args)
}

/**
 * Fatal logs a message and exits the program
 */
func (l *AppLogger) Fatal(ctx context.Context, message string, args ...any) {
	args = append(args, slog.Bool("fatal", true))
	l.logger.ErrorContext(ctx, message, args...)
	fmt.Printf("\033[31m[ERROR]\033[0m %s %s: %v\n", time.Now().UTC().Format(time.RFC3339), message, args)
	os.Exit(1)
}

func (l *AppLogger) Shutdown(ctx context.Context) error {
	return l.loggerProvider.(*sdkLog.LoggerProvider).Shutdown(ctx)
}
