# Notifier Service

The Notifier Service is responsible for processing notifications. It listens to both a Kafka broker and Google Pub/Sub for Gmail notifications, and publishes events to other services as needed.

## Features

- Listens to Google Pub/Sub for Gmail notifications
- Listens to Kafka broker for events
- Publishes processed events to other services
- Integrates with PostgreSQL for message tracking
- Supports distributed tracing and logging

## Dependencies

This service depends on the following external systems:

- **Google Cloud Platform**: Pub/Sub and Gmail API access
- **Kafka Broker**: For event streaming
- **PostgreSQL**: For storing processed messages
- **OpenTelemetry Collector**: For distributed tracing (optional, but recommended)

## Configuration

The service is configured via environment variables. You can use a `.env` file in the root of the service directory. The main configuration options are:

| Variable                        | Description                                      | Example / Default                |
|----------------------------------|--------------------------------------------------|----------------------------------|
| `GOOGLE_PROJECT_ID`              | Google Cloud project ID                          | `my-gcp-project`                 |
| `PUBSUB_TOPIC`                   | Pub/Sub topic to subscribe to                    | `gmail-notifications`            |
| `POSTGRESQL_URL`                 | PostgreSQL connection string                     | `postgres://user:pass@host/db`   |
| `MIGRATIONS_URL`                 | Path/URL to DB migrations                        | `file://db/migrations`           |
| `GOOGLE_OAUTH_CREDENTIALS_FILE`  | Path to Gmail OAuth2 credentials JSON            | `credentials.json`               |
| `GOOGLE_SERVICE_CREDENTIALS_FILE`| Path to Google service account JSON              | `service-credentials.json`       |
| `GOOGLE_TOKEN_FILE`              | Path to Gmail OAuth2 token JSON                  | `token.json`                     |
| `OTEL_COLLECTOR_ENDPOINT`        | OpenTelemetry collector endpoint                 | `localhost:4317`                 |
| `SERVICE_NAME`                   | Service name for tracing/logging                 | `notifier`                       |
| `SERVICE_NAMESPACE`              | Service namespace for tracing/logging            | `tellawl`                        |
| `SERVICE_VERSION`                | Service version                                  | `1.0.0`                          |
| `LOG_LEVEL`                      | Log level (`DEBUG`, `INFO`, `WARN`, `ERROR`)     | `DEBUG`                          |
| `KAFKA_BROKERS`                  | Comma-separated list of Kafka broker addresses   | `localhost:9092`                 |
| `KAFKA_TOPIC`                    | Kafka topic to subscribe to                      |                                  |

### K8s secrets creation

```sh
kubectl create secret generic notifier-credentials \
  --from-file=credentials.json=./services/notifier/credentials.json \
  --from-file=service-credentials.json=./services/notifier/service-credentials.json \
  --from-file=token.json=./services/notifier/token.json \
  -n <namespace>
```

## Running the Service

1. **Install dependencies**
	- Ensure you have Go installed (version 1.20+ recommended).
	- Install and configure PostgreSQL, Kafka, and Google Cloud credentials as needed.
    - `docker-compose.yml` is available at the repo root to help with some of the dependencies.

2. **Set up environment variables**
	- Copy `.env.example` to `.env` and fill in the required values, or set them directly in your environment.

3. **Run database migrations**
	- Migrations are applied automatically on startup, but you can also run them manually if needed.

4. **Start the service**
	- Run:
	  ```sh
	  go run cmd/listener/main.go
	  ```

## Entrypoint

The main entrypoint for the service is at:

- `cmd/listener/main.go`

## Authentication

To access Gmail and Google Pub/Sub, you must provide valid Google OAuth2 credentials and tokens. Use the CLI tool (see `cmd/cli/main.go`) to authenticate and generate the required token file before running the listener.

## License

MIT
