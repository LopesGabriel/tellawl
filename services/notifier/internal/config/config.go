package config

import (
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type AppConfiguration struct {
	PostgreSQLURL          string
	MigrationsURL          string
	GoogleProjectId        string
	PubSubTopic            string
	CredentialsFile        string
	ServiceCredentialsFile string
	TokenFile              string
	OTELCollectorEndpoint  string
	ServiceName            string
	ServiceNamespace       string
	ServiceVersion         string
	LoggerLevel            slog.Level

	KafkaBrokers []string
	KafkaTopic   string

	SMTPHost     string
	SMTPPort     string
	SMTPFrom     string
	SMTPPassword string
}

func InitAppConfigurations() *AppConfiguration {
	_ = godotenv.Load()

	projectId := os.Getenv("GOOGLE_PROJECT_ID")
	if projectId == "" {
		log.Fatal("GOOGLE_PROJECT_ID is required")
	}

	topic := os.Getenv("PUBSUB_TOPIC")
	if topic == "" {
		log.Fatal("PUBSUB_TOPIC is required")
	}

	postgresqlURL := os.Getenv("POSTGRESQL_URL")
	if postgresqlURL == "" {
		log.Fatal("POSTGRESQL_URL is required")
	}

	logLevel := getEnv("LOG_LEVEL", "DEBUG")

	kafkaBrokers := []string{}
	if rawBrokers := getEnv("KAFKA_BROKERS", ""); rawBrokers != "" {
		kafkaBrokers = strings.Split(rawBrokers, ",")
	}

	return &AppConfiguration{
		GoogleProjectId:        projectId,
		PubSubTopic:            topic,
		PostgreSQLURL:          postgresqlURL,
		MigrationsURL:          getEnv("MIGRATIONS_URL", "file://db/migrations"),
		CredentialsFile:        getEnv("GOOGLE_OAUTH_CREDENTIALS_FILE", "credentials.json"),
		ServiceCredentialsFile: getEnv("GOOGLE_SERVICE_CREDENTIALS_FILE", "service-credentials.json"),
		TokenFile:              getEnv("GOOGLE_TOKEN_FILE", "token.json"),
		OTELCollectorEndpoint:  getEnv("OTEL_COLLECTOR_ENDPOINT", "localhost:4317"),
		ServiceName:            getEnv("SERVICE_NAME", "notifier"),
		ServiceNamespace:       getEnv("SERVICE_NAMESPACE", "tellawl"),
		ServiceVersion:         getEnv("SERVICE_VERSION", "1.0.0"),
		LoggerLevel:            parseLogLevel(logLevel),
		KafkaBrokers:           kafkaBrokers,
		KafkaTopic:             getEnv("KAFKA_TOPIC", ""),
		SMTPHost:               getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:               getEnv("SMTP_PORT", "587"),
		SMTPFrom:               getEnv("SMTP_FROM", ""),
		SMTPPassword:           getEnv("SMTP_PASSWORD", ""),
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
