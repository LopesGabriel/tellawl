package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/pubsub/v2"
	"go.opentelemetry.io/otel"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"github.com/lopesgabriel/tellawl/packages/broker"
	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/packages/tracing"
	"github.com/lopesgabriel/tellawl/services/notifier/internal/config"
	"github.com/lopesgabriel/tellawl/services/notifier/internal/domain/inbox"
	"github.com/lopesgabriel/tellawl/services/notifier/internal/infra/database"
	"github.com/lopesgabriel/tellawl/services/notifier/internal/infra/listener"
	"github.com/lopesgabriel/tellawl/services/notifier/internal/infra/publisher"
	"github.com/lopesgabriel/tellawl/services/notifier/internal/infra/telegram"
)

func main() {
	ctx := context.Background()

	// Carrega as configurações da aplicação
	config := config.InitAppConfigurations()

	// Inicializa a telemetria (OpenTelemetry)
	shutdown, err := initTelemetry(ctx, config)
	if err != nil {
		log.Fatalf("Erro ao inicializar telemetria: %s", err)
	}
	defer shutdown()

	applogger, err := logger.GetLogger()
	if err != nil {
		log.Fatalf("Erro ao obter logger: %s", err)
	}

	applogger.Info(ctx, "Iniciando o Notifier Service...")

	var kafkaBroker broker.Broker
	if len(config.KafkaBrokers) > 0 && config.KafkaTopic != "" {
		kafkaBroker, err = broker.NewKafkaBroker(broker.NewKafkaBrokerArgs{
			BootstrapServers: config.KafkaBrokers,
			Service:          config.ServiceName,
			Topic:            config.KafkaTopic,
			Logger:           applogger,
		})
	}

	// Inicializa o publisher de eventos
	eventPublisher := publisher.InitEventPublisher(ctx, config, applogger, kafkaBroker)

	// Inicializa o banco de dados (PostgreSQL)
	db := initDatabase(ctx, config, applogger)
	defer db.Close()

	dbTracer := tracing.GetTracer("github.com/lopesgabriel/tellawl/services/notifier/internal/infra/database")
	processedMessagesRepository := database.NewPostgreSQLProcessedMessagesRepository(
		db,
		eventPublisher,
		dbTracer,
	)
	telegramRepository := database.NewPostgreSQLTelegramNotificationTargetRepository(
		db,
		dbTracer,
	)

	// Telegram Bot API Client
	telegramClient := telegram.NewClient(config.TelegramBotToken)

	// Inicia o Kafka consumer
	kafkaListener := listener.NewKafkaListener(listener.NewKafkaListenerParams{
		Topic:                       config.KafkaTopic,
		Broker:                      kafkaBroker,
		ProcessedMessagesRepository: processedMessagesRepository,
		Tracer:                      tracing.GetTracer("github.com/lopesgabriel/tellawl/services/notifier/internal/infra/listener"),
		AppLogger:                   applogger,
		TelegramClient:              telegramClient,
		TelegramRepo:                telegramRepository,
	})
	kafkaListener.Start()

	// Cliente OAuth2 para Gmail
	b, err := os.ReadFile(config.CredentialsFile)
	if err != nil {
		applogger.Fatal(ctx, "Erro ao ler credenciais OAuth2", slog.String("file", config.CredentialsFile), slog.Any("error", err))
	}

	oauthConfig, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope, pubsub.ScopePubSub)
	if err != nil {
		applogger.Fatal(ctx, "Erro ao parsear config OAuth2", slog.Any("error", err))
	}

	httpClient, err := getClient(oauthConfig, config.TokenFile)
	if err != nil {
		applogger.Fatal(ctx, "Erro ao obter cliente HTTP", slog.Any("error", err))
	}

	// Serviço Gmail
	gmailService, err := gmail.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		applogger.Fatal(ctx, "Erro ao criar serviço Gmail", slog.Any("error", err))
	}

	// Cliente PubSub
	psClient, err := pubsub.NewClient(ctx, config.GoogleProjectId, option.WithAuthCredentialsFile(option.ServiceAccount, config.ServiceCredentialsFile))
	if err != nil {
		applogger.Fatal(ctx, "Erro ao criar cliente PubSub", slog.Any("error", err))
	}
	defer psClient.Close()

	// Inicializa o listener
	listener := inbox.NewPubSubListener(ctx, inbox.NewPubSubListenerParams{
		ProjectId:                   config.GoogleProjectId,
		Topic:                       config.PubSubTopic,
		PsClient:                    psClient,
		GmailService:                gmailService,
		Publisher:                   eventPublisher,
		ProcessedMessagesRepository: processedMessagesRepository,
	})
	//config.GoogleProjectId, config.PubSubTopic, psClient, gmailService, eventPublisher

	// Inicia o listener em uma goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- listener.Start()
	}()

	// Escuta sinais do sistema para graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		applogger.Info(ctx, "Sinal recebido. Encerrando...", slog.String("signal", sig.String()))
		if err := listener.Stop(); err != nil {
			applogger.Fatal(ctx, "Erro ao parar listener", slog.Any("error", err))
		}
	case err := <-errCh:
		if err != nil {
			applogger.Fatal(ctx, "Listener encerrado com erro", slog.Any("error", err))
		}
	}

	applogger.Info(ctx, "Aplicação encerrada.")
}

func getClient(config *oauth2.Config, tokenFile string) (*http.Client, error) {
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		return nil, fmt.Errorf("Token não encontrado em %s. Execute o CLI para autenticar primeiro.", tokenFile)
	}
	return config.Client(context.Background(), tok), nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func initTelemetry(ctx context.Context, appConfig *config.AppConfiguration) (func() error, error) {
	appLogger, err := logger.Init(ctx, logger.InitLoggerArgs{
		CollectorURL:     appConfig.OTELCollectorEndpoint,
		ServiceName:      appConfig.ServiceName,
		ServiceNamespace: appConfig.ServiceNamespace,
		ServiceVersion:   appConfig.ServiceVersion,
		Level:            appConfig.LoggerLevel,
	})
	if err != nil {
		return nil, err
	}

	traceProvider, err := tracing.Init(ctx, tracing.NewTraceProviderArgs{
		CollectorURL:     appConfig.OTELCollectorEndpoint,
		ServiceName:      appConfig.ServiceName,
		ServiceNamespace: appConfig.ServiceNamespace,
		ServiceVersion:   appConfig.ServiceVersion,
	})
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(traceProvider)

	return func() error {
		if err := appLogger.Shutdown(ctx); err != nil {
			return err
		}
		if err := traceProvider.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	}, nil
}

func initDatabase(ctx context.Context, appConfig *config.AppConfiguration, appLogger *logger.AppLogger) *sql.DB {
	db, err := database.NewPostgresClient(context.Background(), appConfig.PostgreSQLURL)
	if err != nil {
		appLogger.Fatal(ctx, "failed to create the postgres client", slog.String("error", err.Error()))
	}

	err = db.Ping()
	if err != nil {
		appLogger.Fatal(ctx, "failed to ping database", slog.String("error", err.Error()))
	}

	appLogger.Info(ctx, "Connected to the database successfully")

	err = database.MigrateUp(appConfig.MigrationsURL, appConfig.PostgreSQLURL)
	if err != nil {
		appLogger.Fatal(ctx, "failed to apply database migration", slog.String("error", err.Error()))
	}

	return db
}
