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
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
	}

	rawPort := os.Getenv("PORT")
	if rawPort == "" {
		rawPort = "8080"
	}

	logLevel := slog.LevelInfo
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		logLevel = parseLogLevel(level)
	}

	port, err := strconv.Atoi(rawPort)
	if err != nil {
		port = 8080
	}

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "member-service"
	}

	serviceNamespace := os.Getenv("SERVICE_NAMESPACE")
	if serviceNamespace == "" {
		serviceNamespace = "tellawl"
	}

	brokers := []string{}
	if os.Getenv("KAFKA_BROKERS") != "" {
		brokers = strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	}

	return &AppConfiguration{
		JwtSecret:        os.Getenv("JWT_SECRET"),
		Version:          os.Getenv("VERSION"),
		Port:             port,
		LogLevel:         logLevel,
		DatabaseUrl:      os.Getenv("POSTGRESQL_URL"),
		MigrationUrl:     os.Getenv("MIGRATIONS_URL"),
		OTELCollectorUrl: os.Getenv("OTEL_COLLECTOR_URL"),
		ServiceName:      serviceName,
		ServiceNamespace: serviceNamespace,
		WalletTopic:      os.Getenv("WALLET_TOPIC"),
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
