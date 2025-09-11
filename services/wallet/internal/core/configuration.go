package core

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Configuration struct {
	JwtSecret        string
	Version          string
	Port             int
	DatabaseUrl      string
	MigrationUrl     string
	OTELCollectorUrl string
}

func InitAppConfigurations() *Configuration {
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

	return &Configuration{
		JwtSecret:        os.Getenv("JWT_SECRET"),
		Version:          os.Getenv("VERSION"),
		Port:             port,
		DatabaseUrl:      os.Getenv("POSTGRESQL_URL"),
		MigrationUrl:     os.Getenv("MIGRATIONS_URL"),
		OTELCollectorUrl: os.Getenv("OTEL_COLLECTOR_URL"),
	}
}
