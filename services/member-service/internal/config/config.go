package config

import (
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

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")

	return &AppConfiguration{
		JwtSecret:        os.Getenv("JWT_SECRET"),
		Version:          os.Getenv("VERSION"),
		Port:             port,
		DatabaseUrl:      os.Getenv("POSTGRESQL_URL"),
		MigrationUrl:     os.Getenv("MIGRATIONS_URL"),
		OTELCollectorUrl: os.Getenv("OTEL_COLLECTOR_URL"),
		ServiceName:      serviceName,
		ServiceNamespace: serviceNamespace,
		WalletTopic:      os.Getenv("WALLET_TOPIC"),
		KafkaBrokers:     brokers,
	}
}
