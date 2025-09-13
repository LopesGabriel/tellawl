package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
)

type InitLoggerArgs struct {
	CollectorURL     string
	ServiceName      string
	ServiceNamespace string
	ServiceVersion   string
}

var logger *slog.Logger

func Init(ctx context.Context, args InitLoggerArgs) (*log.LoggerProvider, error) {
	if args.ServiceNamespace == "" {
		args.ServiceNamespace = "tellawl"
	}

	res, err := newResource(args)
	if err != nil {
		return nil, err
	}

	logProvider, err := newProvider(ctx, res, args)
	if err != nil {
		return nil, err
	}

	global.SetLoggerProvider(logProvider)

	otelslog.NewLogger(
		fmt.Sprintf("%s/%s", args.ServiceNamespace, args.ServiceName),
		otelslog.WithLoggerProvider(logProvider),
	)

	logger = otelslog.NewLogger(
		fmt.Sprintf("%s/%s", args.ServiceNamespace, args.ServiceName),
		otelslog.WithLoggerProvider(logProvider),
	)

	return logProvider, nil
}

func Info(ctx context.Context, message string, args ...any) {
	logger.InfoContext(ctx, message, args...)
	fmt.Printf("\033[34m[INFO]\033[0m %s %s: %v\n", time.Now().UTC().Format(time.RFC3339), message, args)
}

func Error(ctx context.Context, message string, args ...any) {
	args = append(args, slog.Bool("fatal", false))
	logger.ErrorContext(ctx, message, args...)
	fmt.Printf("\033[31m[ERROR]\033[0m %s %s: %v\n", time.Now().UTC().Format(time.RFC3339), message, args)
}

func Debug(ctx context.Context, message string, args ...any) {
	logger.DebugContext(ctx, message, args...)
	fmt.Printf("\033[36m[DEBUG]\033[0m %s %s: %v\n", time.Now().UTC().Format(time.RFC3339), message, args)
}

func Warn(ctx context.Context, message string, args ...any) {
	logger.WarnContext(ctx, message, args...)
	fmt.Printf("\033[33m[WARN]\033[0m %s %s: %v\n", time.Now().UTC().Format(time.RFC3339), message, args)
}

/**
 * Fatal logs a message and exits the program
 */
func Fatal(ctx context.Context, message string, args ...any) {
	args = append(args, slog.Bool("fatal", true))
	logger.ErrorContext(ctx, message, args...)
	fmt.Printf("\033[31m[ERROR]\033[0m %s %s: %v\n", time.Now().UTC().Format(time.RFC3339), message, args)
	os.Exit(1)
}
