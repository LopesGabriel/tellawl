package config

import (
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type AppConfiguration struct {
	Version          string
	Port             int
	DatabaseUrl      string
	MemberServiceUrl string
	MigrationUrl     string
	OTELCollectorUrl string
	ServiceName      string
	ServiceNamespace string
	KafkaTopic       string
	KafkaBrokers     []string
	LogLevel         slog.Level
}

func InitAppConfigurations() *AppConfiguration {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
	}

	rawPort := getEnv("PORT", "8080")
	port, err := strconv.Atoi(rawPort)
	if err != nil {
		port = 8080
	}

	brokers := strings.Split(getEnv("KAFKA_BROKERS", ""), ",")

	return &AppConfiguration{
		Version:          getEnv("VERSION", "1.0.0"),
		Port:             port,
		DatabaseUrl:      getEnv("POSTGRESQL_URL", ""),
		MemberServiceUrl: getEnv("MEMBER_SERVICE_URL", ""),
		MigrationUrl:     getEnv("MIGRATIONS_URL", "file://db/migrations"),
		OTELCollectorUrl: getEnv("OTEL_COLLECTOR_URL", "localhost:4317"),
		ServiceName:      getEnv("SERVICE_NAME", "wallet"),
		ServiceNamespace: getEnv("SERVICE_NAMESPACE", "tellawl"),
		KafkaTopic:       getEnv("KAFKA_TOPIC", ""),
		KafkaBrokers:     brokers,
		LogLevel:         parseLogLevel(getEnv("LOG_LEVEL", "INFO")),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func parseLogLevel(level string) slog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG", "TRACE":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		log.Printf("Invalid LOG_LEVEL '%s', defaulting to INFO", level)
		return slog.LevelInfo
	}
}
