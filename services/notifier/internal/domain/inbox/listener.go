package inbox

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"log/slog"
	"time"

	"cloud.google.com/go/iam/apiv1/iampb"
	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/packages/tracing"
	"github.com/lopesgabriel/tellawl/services/notifier/internal/domain/events"
	"github.com/lopesgabriel/tellawl/services/notifier/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/notifier/internal/domain/repositories"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Listener define a interface para escutar notificações de email.
type Listener interface {
	Start() error
	Stop() error
	HandleMessage(ctx context.Context, message []byte) error
}

type pubSubListener struct {
	projectId                   string
	topic                       string
	subscription                string
	psClient                    *pubsub.Client
	gmailService                *gmail.Service
	publisher                   events.EventPublisher
	ctx                         context.Context
	cancel                      context.CancelFunc
	done                        chan struct{}
	logger                      *logger.AppLogger
	tracer                      trace.Tracer
	processedMessagesRepository repositories.ProcessedMessagesRepository
}

type NewPubSubListenerParams struct {
	ProjectId                   string
	Topic                       string
	PsClient                    *pubsub.Client
	GmailService                *gmail.Service
	Publisher                   events.EventPublisher
	ProcessedMessagesRepository repositories.ProcessedMessagesRepository
}

// NewPubSubListener cria um novo Listener baseado em Google PubSub e Gmail API.
func NewPubSubListener(ctx context.Context, params NewPubSubListenerParams) Listener {
	ctx, cancel := context.WithCancel(ctx)
	appLogger, err := logger.GetLogger()
	if err != nil {
		log.Fatalf("Erro ao obter logger: %v", err)
	}

	return &pubSubListener{
		projectId:                   params.ProjectId,
		topic:                       params.Topic,
		subscription:                params.Topic + "-sub",
		psClient:                    params.PsClient,
		gmailService:                params.GmailService,
		publisher:                   params.Publisher,
		processedMessagesRepository: params.ProcessedMessagesRepository,
		logger:                      appLogger,
		tracer:                      tracing.GetTracer("github.com/lopesgabriel/tellawl/services/notifier/internal/domain/inbox"),
		ctx:                         ctx,
		cancel:                      cancel,
		done:                        make(chan struct{}),
	}
}

func (l *pubSubListener) topicKey() string {
	return fmt.Sprintf("projects/%s/topics/%s", l.projectId, l.topic)
}

func (l *pubSubListener) subscriptionKey() string {
	return fmt.Sprintf("projects/%s/subscriptions/%s", l.projectId, l.subscription)
}

// ensureTopic verifica se o tópico existe. Se não existir, cria.
func (l *pubSubListener) ensureTopic() error {
	_, err := l.psClient.TopicAdminClient.GetTopic(l.ctx, &pubsubpb.GetTopicRequest{
		Topic: l.topicKey(),
	})
	if err == nil {
		l.logger.Info(l.ctx, "Tópico já existe: %s", l.topicKey())
		return nil
	}

	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.NotFound {
		return fmt.Errorf("erro ao verificar tópico: %w", err)
	}

	_, err = l.psClient.TopicAdminClient.CreateTopic(l.ctx, &pubsubpb.Topic{
		Name: l.topicKey(),
	})
	if err != nil {
		return fmt.Errorf("erro ao criar tópico: %w", err)
	}

	l.logger.Info(l.ctx, "Tópico criado: %s", l.topicKey())
	return nil
}

// ensureSubscription verifica se a subscription existe. Se não existir, cria.
func (l *pubSubListener) ensureSubscription() error {
	_, err := l.psClient.SubscriptionAdminClient.GetSubscription(l.ctx, &pubsubpb.GetSubscriptionRequest{
		Subscription: l.subscriptionKey(),
	})
	if err == nil {
		l.logger.Info(l.ctx, "Subscription já existe: %s", l.subscriptionKey())
		return nil
	}

	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.NotFound {
		return fmt.Errorf("erro ao verificar subscription: %w", err)
	}

	_, err = l.psClient.SubscriptionAdminClient.CreateSubscription(l.ctx, &pubsubpb.Subscription{
		Name:  l.subscriptionKey(),
		Topic: l.topicKey(),
	})
	if err != nil {
		return fmt.Errorf("erro ao criar subscription: %w", err)
	}

	l.logger.Info(l.ctx, "Subscription criada: %s", l.subscriptionKey())
	return nil
}

// ensureTopicPermissions garante que a conta de serviço do Gmail tenha permissão de publicação no tópico.
func (l *pubSubListener) ensureTopicPermissions() error {
	policy, err := l.psClient.TopicAdminClient.GetIamPolicy(l.ctx, &iampb.GetIamPolicyRequest{
		Resource: l.topicKey(),
	})
	if err != nil {
		return fmt.Errorf("erro ao obter IAM do tópico: %w", err)
	}

	principal := "serviceAccount:gmail-api-push@system.gserviceaccount.com"
	role := "roles/pubsub.publisher"

	foundMember := false
	var targetBinding *iampb.Binding

	for _, binding := range policy.Bindings {
		if binding.Role == role {
			targetBinding = binding
			for _, member := range binding.Members {
				if member == principal {
					foundMember = true
					break
				}
			}
			break
		}
	}

	if foundMember {
		l.logger.Info(l.ctx, "Permissão 'Pub/Sub Publisher' já configurada para %s", principal)
		return nil
	}

	if targetBinding != nil {
		targetBinding.Members = append(targetBinding.Members, principal)
	} else {
		policy.Bindings = append(policy.Bindings, &iampb.Binding{
			Role:    role,
			Members: []string{principal},
		})
	}

	_, err = l.psClient.TopicAdminClient.SetIamPolicy(l.ctx, &iampb.SetIamPolicyRequest{
		Resource: l.topicKey(),
		Policy:   policy,
	})
	if err != nil {
		return fmt.Errorf("erro ao definir IAM do tópico: %w", err)
	}

	l.logger.Info(l.ctx, "Permissão 'Pub/Sub Publisher' concedida a %s", principal)
	return nil
}

// ensureWatch configura o Gmail Watch na INBOX para o tópico PubSub.
// A API do Gmail permite chamar Watch novamente para renovar, então é sempre seguro executar.
func (l *pubSubListener) ensureWatch() error {
	req := &gmail.WatchRequest{
		TopicName: l.topicKey(),
		LabelIds:  []string{"INBOX"},
	}
	res, err := l.gmailService.Users.Watch("me", req).Do()
	if err != nil {
		return fmt.Errorf("erro ao configurar watch: %w", err)
	}

	l.logger.Info(l.ctx, "Watch configurado! HistoryId: %v, Expira em: %v", res.HistoryId, res.Expiration)
	return nil
}

// Start inicializa os recursos necessários (tópico, subscription, watch) e começa a
// escutar mensagens. Este método bloqueia até que Stop() seja chamado ou ocorra um erro.
// Deve ser chamado em uma goroutine separada.
func (l *pubSubListener) Start() error {
	defer close(l.done)

	if err := l.ensureTopic(); err != nil {
		return err
	}

	if err := l.ensureTopicPermissions(); err != nil {
		return err
	}

	if err := l.ensureSubscription(); err != nil {
		return err
	}
	if err := l.ensureWatch(); err != nil {
		return err
	}

	// Inicia a renovação periódica do watch
	go l.startWatchRenewal()

	l.logger.Info(l.ctx, "Listener iniciado. Aguardando mensagens no tópico %s...", l.topic)

	sub := l.psClient.Subscriber(l.subscriptionKey())
	err := sub.Receive(l.ctx, func(ctx context.Context, m *pubsub.Message) {
		ctx, span := l.tracer.Start(ctx, "NotificationReceived")
		defer span.End()

		if err := l.HandleMessage(ctx, m.Data); err != nil {
			l.logger.Error(ctx, "Erro ao processar mensagem: %v", err)
			span.SetStatus(otelcodes.Error, "Erro ao processar mensagem")
			span.RecordError(err)
			m.Nack()
			return
		}

		span.SetStatus(otelcodes.Ok, "success")
		m.Ack()
	})

	// Se o contexto foi cancelado (shutdown graceful), não é erro
	if err != nil && l.ctx.Err() != nil {
		return nil
	}
	return err
}

// Stop cancela o contexto do listener, aguardando que qualquer processamento em andamento
// seja finalizado antes de retornar.
func (l *pubSubListener) Stop() error {
	l.logger.Info(l.ctx, "Finalizando listener...")
	l.cancel()
	<-l.done
	l.logger.Info(l.ctx, "Listener finalizado.")
	return nil
}

func (l *pubSubListener) startWatchRenewal() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-l.ctx.Done():
			return
		case <-ticker.C:
			l.logger.Info(l.ctx, "Renovando watch do Gmail...")
			if err := l.ensureWatch(); err != nil {
				l.logger.Error(l.ctx, "Erro ao renovar watch: %v", err)
			}
		}
	}
}

// HandleMessage processa uma notificação recebida via PubSub buscando o email mais
// recente na inbox do Gmail e logando seus detalhes.
func (l *pubSubListener) HandleMessage(ctx context.Context, message []byte) error {
	l.logger.Info(ctx, "Notificação recebida! Buscando e-mail...")

	list, err := l.gmailService.Users.Messages.List("me").MaxResults(1).Do()
	if err != nil {
		return fmt.Errorf("erro ao listar mensagens: %w", err)
	}

	if len(list.Messages) == 0 {
		l.logger.Info(ctx, "Nenhuma mensagem encontrada.")
		return nil
	}

	msgId := list.Messages[0].Id

	l.logger.Info(ctx, "Mensagem encontrada! Verificando se já foi processada...", slog.String("message.id", msgId))

	exists, err := l.processedMessagesRepository.Exists(ctx, msgId)
	if err != nil {
		return fmt.Errorf("erro ao verificar se mensagem já foi processada: %w", err)
	}

	if exists {
		l.logger.Info(ctx, "Mensagem já foi processada. Ignorando.", slog.String("message.id", msgId))
		return nil
	}

	l.logger.Info(ctx, "Mensagem ainda não processada. Buscando detalhes...", slog.String("message.id", msgId))

	msg, err := l.gmailService.Users.Messages.Get("me", msgId).Do(
		googleapi.QueryParameter("format", "full"),
	)
	if err != nil {
		return fmt.Errorf("erro ao buscar mensagem %s: %w", msgId, err)
	}

	var subject, sender, recipient string
	if msg.Payload != nil {
		for _, header := range msg.Payload.Headers {
			switch header.Name {
			case "Subject":
				subject = header.Value
			case "From":
				sender = header.Value
			case "To":
				recipient = header.Value
			}
		}
	}

	body := getBody(msg.Payload)
	processedMessage := models.NewProcessedMessage(msgId, recipient, subject, sender, body)
	err = l.processedMessagesRepository.Save(ctx, processedMessage)

	return err
}

func getBody(payload *gmail.MessagePart) string {
	if payload == nil {
		return ""
	}

	if payload.Body != nil && payload.Body.Data != "" {
		data, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err != nil {
			return "Erro ao decodificar: " + err.Error()
		}
		return string(data)
	}

	for _, part := range payload.Parts {
		body := getBody(part)
		if body != "" {
			return body
		}
	}
	return ""
}
