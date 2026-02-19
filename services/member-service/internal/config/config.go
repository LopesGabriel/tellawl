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
	JwtSecret        string
	Version          string
	Port             int
	LogLevel         slog.Level
	DatabaseUrl      string
	MigrationUrl     string
	OTELCollectorUrl string
	ServiceName      string
	ServiceNamespace string
	WalletTopic      string
	KafkaBrokers     []string
}

func InitAppConfigurations() *AppConfiguration {
	_ = godotenv.Load()

	rawPort := getEnv("PORT", "8080")
	port, err := strconv.Atoi(rawPort)
	if err != nil {
		port = 8080
	}

	brokers := []string{}
	if os.Getenv("KAFKA_BROKERS") != "" {
		brokers = strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	}

	return &AppConfiguration{
		JwtSecret:        getEnv("JWT_SECRET", ""),
		Version:          getEnv("VERSION", "1.0.0"),
		Port:             port,
		LogLevel:         parseLogLevel(getEnv("LOG_LEVEL", "INFO")),
		DatabaseUrl:      getEnv("POSTGRESQL_URL", ""),
		MigrationUrl:     getEnv("MIGRATIONS_URL", "file://db/migrations"),
		OTELCollectorUrl: getEnv("OTEL_COLLECTOR_URL", ""),
		ServiceName:      getEnv("SERVICE_NAME", "member-service"),
		ServiceNamespace: getEnv("SERVICE_NAMESPACE", "tellawl"),
		WalletTopic:      getEnv("WALLET_TOPIC", ""),
		KafkaBrokers:     brokers,
	}
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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
