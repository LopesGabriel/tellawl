package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type CredentialsFile struct {
	Installed struct {
		ProjectId string `json:"project_id"`
	} `json:"installed"`
}

func main() {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Erro ao ler credenciais: %v", err)
	}

	var credFile CredentialsFile
	err = json.Unmarshal(b, &credFile)
	if err != nil {
		log.Fatalf("Erro ao realizar o parse do arquivo de credenciais: %v", err)
	}

	// Define os escopos necessários
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope, pubsub.ScopePubSub)
	if err != nil {
		log.Fatalf("Erro ao parsear config: %v", err)
	}
	client := getClient(config)

	// Inicializa o cliente Pub/Sub
	psClient, err := pubsub.NewClient(ctx, credFile.Installed.ProjectId, option.WithAuthCredentialsFile(option.ServiceAccount, "service-credentials.json"))
	if err != nil {
		log.Fatalf("Erro PubSub: %v", err)
	}
	defer psClient.Close()

	// Definindo flags
	topicFlag := flag.String("topic", "", "Tópico da operação")
	createFlag := flag.Bool("create", false, "Cria um novo topico")
	deleteFlag := flag.Bool("delete", false, "Deleta um topico")
	showFlag := flag.Bool("show", false, "Mostra o status de um topico")
	watchFlag := flag.Bool("watch", false, "Iniciar o Watch")
	subscribeFlag := flag.Bool("sub", false, "Subscribe to a project")
	listenFlag := flag.Bool("listen", false, "Listen to a Subscription")
	readFlag := flag.Bool("read", false, "Ler um email")
	helpFlag := flag.Bool("help", false, "Mostra ajuda")

	// Fazendo parse dos argumentos
	flag.Parse()

	// Se flag de ajuda foi passado, mostra instruções
	if *helpFlag {
		printHelp()
		return
	}

	if *topicFlag == "" {
		log.Fatal("Erro: Informe o nome do tópico")
	}

	topicKey := fmt.Sprintf("projects/%s/topics/%s", credFile.Installed.ProjectId, *topicFlag)
	subID := *topicFlag + "-sub"
	subscriptionKey := fmt.Sprintf("projects/%s/subscriptions/%s", credFile.Installed.ProjectId, subID)

	// Validando argumentos obrigatórios
	if !*createFlag && !*deleteFlag && !*showFlag && !*watchFlag && !*subscribeFlag && !*listenFlag && !*readFlag {
		log.Fatal("Erro: informe ao menos um dos argumentos --create, --delete, --show, --watch, --sub ou --listen")
	}

	// Usando argumentos posicionais (não flag)
	positionalArgs := flag.Args()

	// Executando logica
	fmt.Printf("🔔 Topic: %s\n", *topicFlag)
	if len(positionalArgs) > 0 {
		fmt.Printf("📌 Argumentos adicionais: %v\n", positionalArgs)
	}

	if *createFlag {
		t, err := psClient.TopicAdminClient.CreateTopic(ctx, &pubsubpb.Topic{
			Name: topicKey,
		})
		if err != nil {
			log.Fatalf("Erro ao criar: %v", err)
		}
		fmt.Printf("Tópico criado: %v\n", t.Name)
	}

	if *deleteFlag {
		deleteArgs := &pubsubpb.DeleteTopicRequest{Topic: topicKey}
		if err := psClient.TopicAdminClient.DeleteTopic(ctx, deleteArgs); err != nil {
			log.Fatalf("Erro ao deletar: %v", err)
		}
		fmt.Println("Tópico deletado com sucesso.")
	}

	if *showFlag {
		topic, err := psClient.TopicAdminClient.GetTopic(ctx, &pubsubpb.GetTopicRequest{
			Topic: topicKey,
		})
		if err != nil {
			log.Fatalf("Erro ao recuperar tópico: %v", err)
		}

		fmt.Printf("Tópico: %s\nStatus: Ativo\nRetenção: %v\n", topic.Name, topic.GetMessageRetentionDuration())

		// Listar assinaturas vinculadas
		it := psClient.TopicAdminClient.ListTopicSubscriptions(ctx, &pubsubpb.ListTopicSubscriptionsRequest{
			Topic: topic.Name,
		})
		fmt.Println("Assinaturas:")
		for {
			sub, err := it.Next()
			if err != nil {
				break
			}
			fmt.Printf("- %s\n", sub)
		}
	}

	if *watchFlag {
		gmailService, err := gmail.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			log.Fatalf("Erro ao criar serviço Gmail: %v", err)
		}

		req := &gmail.WatchRequest{
			TopicName: topicKey,
			LabelIds:  []string{"INBOX"},
		}
		res, err := gmailService.Users.Watch("me", req).Do()
		if err != nil {
			log.Fatalf("Erro no Watch: %v", err)
		}
		fmt.Printf("Watch ativo! ID da História: %v, Expira em: %v\n", res.HistoryId, res.Expiration)
	}

	if *readFlag {
		gmailService, err := gmail.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			log.Fatalf("Erro ao criar serviço Gmail: %v", err)
		}

		list, err := gmailService.Users.Messages.List("me").MaxResults(1).Do()
		if err != nil {
			log.Fatalf("Erro ao listar histórico: %v\n", err)
		}

		msgId := list.Messages[0].Id

		msg, err := gmailService.Users.Messages.Get("me", msgId).Do(
			googleapi.QueryParameter("format", "full"),
		)

		var subject, sender string
		if msg.Payload != nil {
			for _, header := range msg.Payload.Headers {
				if header.Name == "Subject" {
					subject = header.Value
				}
				if header.Name == "From" {
					sender = header.Value
				}
			}
		}

		corpo := getBody(msg.Payload)

		fmt.Printf("\nMensagem %s:\n Subject: %s\n Sender: %s\n Payload: %s\n", msg.Id, subject, sender, corpo)
	}

	if *subscribeFlag {
		sub, err := psClient.SubscriptionAdminClient.CreateSubscription(ctx, &pubsubpb.Subscription{
			Name:  subscriptionKey,
			Topic: topicKey,
		})
		if err != nil {
			log.Fatalf("Erro ao criar subscription: %v\n", err)
		}
		fmt.Printf("Subscription criada: %v\n", sub.Name)
	}

	if *listenFlag {
		gmailService, err := gmail.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			log.Fatalf("Erro ao criar serviço Gmail: %v", err)
		}

		sub := psClient.Subscriber(subscriptionKey)
		err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
			defer m.Ack()
			fmt.Println("🔔 Notificação recebida! Buscando e-mail...")

			list, err := gmailService.Users.Messages.List("me").MaxResults(1).Do()
			if err != nil {
				log.Fatalf("Erro ao listar histórico: %v\n", err)
			}

			msgId := list.Messages[0].Id

			msg, err := gmailService.Users.Messages.Get("me", msgId).Do(
				googleapi.QueryParameter("format", "full"),
			)

			var subject, sender string
			if msg.Payload != nil {
				for _, header := range msg.Payload.Headers {
					if header.Name == "Subject" {
						subject = header.Value
					}
					if header.Name == "From" {
						sender = header.Value
					}
				}
			}

			corpo := getBody(msg.Payload)

			fmt.Printf("\nMensagem %s:\n Subject: %s\n Sender: %s\n Payload: %s\n", msg.Id, subject, sender, corpo)
			fmt.Println("--------------")
		})
	}
}

func printHelp() {
	fmt.Fprintf(os.Stderr, `Notifier CLI - Sistema de notificações

Uso:
  notifier [flags] [argumentos adicionais]

Flags:
  -topic string
      Nome do tópico que deseja criar/deletar/exibir (obrigatório)
  -create
      Ação de criar
	-delete
			Ação de deletar
	-show
			Ação de exibir detalhes
	-watch
			Ação de criar uma watch
  -help
      Mostra esta ajuda

Exemplos:
  notifier -create -topic gmail-inbox
  notifier -delete -topic gmail-inbox
  notifier -show -topic gmail-inbox
  notifier -watch -topic gmail-inbox
`)
}

// getClient recupera um token, salva-o e retorna o cliente gerado.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Solicita o token pela web
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Abra o link no navegador e digite o código de autorização: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Erro ao ler código: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Erro ao resgatar token: %v", err)
	}
	return tok
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

func saveToken(path string, token *oauth2.Token) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Erro ao salvar token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getBody(payload *gmail.MessagePart) string {
	if payload.Body.Data != "" {
		// O Gmail usa Base64 URL-Safe
		data, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err != nil {
			return "Erro ao decodificar: " + err.Error()
		}
		return string(data)
	}

	// Se não estiver na raiz, procura nas sub-partes (comum em e-mails com anexos)
	for _, part := range payload.Parts {
		body := getBody(part)
		if body != "" {
			return body
		}
	}
	return ""
}
