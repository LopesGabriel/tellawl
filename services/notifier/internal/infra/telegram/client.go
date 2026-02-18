package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const telegramAPIBaseURL = "https://api.telegram.org"

type ParseMode string

const (
	ParseModeMarkdown ParseMode = "Markdown"
	ParseModeHTML     ParseMode = "HTML"
)

type ChatAction string

const (
	ChatActionTyping          ChatAction = "typing"
	ChatActionUploadPhoto     ChatAction = "upload_photo"
	ChatActionRecordVideo     ChatAction = "record_video"
	ChatActionUploadVideo     ChatAction = "upload_video"
	ChatActionRecordVoice     ChatAction = "record_voice"
	ChatActionUploadVoice     ChatAction = "upload_voice"
	ChatActionUploadDocument  ChatAction = "upload_document"
	ChatActionFindLocation    ChatAction = "find_location"
	ChatActionRecordVideoNote ChatAction = "record_video_note"
	ChatActionUploadVideoNote ChatAction = "upload_video_note"
)

type SendChatActionRequest struct {
	ChatID int        `json:"chat_id"`
	Action ChatAction `json:"action"`
}

type SendMessageRequest struct {
	ChatID    int       `json:"chat_id"`
	Text      string    `json:"text"`
	ParseMode ParseMode `json:"parse_mode,omitempty"`
}

type apiResponse[T any] struct {
	OK          bool   `json:"ok"`
	Result      T      `json:"result"`
	Description string `json:"description,omitempty"`
}

// Client is an HTTP client for the Telegram Bot API.
type Client struct {
	botToken   string
	httpClient *http.Client
}

// NewClient creates a new Telegram Client with an OpenTelemetry-instrumented HTTP client.
func NewClient(botToken string) *Client {
	return &Client{
		botToken:   botToken,
		httpClient: &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)},
	}
}

func (c *Client) botURL(method string) string {
	return fmt.Sprintf("%s/bot%s/%s", telegramAPIBaseURL, c.botToken, method)
}

func (c *Client) post(ctx context.Context, method string, body any) (*http.Response, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("telegram: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.botURL(method), bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("telegram: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(req)
}

// SetWebhook registers a webhook URL with the Telegram Bot API.
func (c *Client) SetWebhook(ctx context.Context, url string) error {
	resp, err := c.post(ctx, "setWebhook", map[string]string{"url": url})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telegram: setWebhook failed: %s", resp.Status)
	}
	return nil
}

// SendChatAction sends a chat action (e.g. "typing") to the given chat.
func (c *Client) SendChatAction(ctx context.Context, req SendChatActionRequest) error {
	resp, err := c.post(ctx, "sendChatAction", req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telegram: sendChatAction failed: %s", resp.Status)
	}
	return nil
}

// SendMessage sends a text message to the given chat and returns the sent Message.
func (c *Client) SendMessage(ctx context.Context, req SendMessageRequest) (*Message, error) {
	resp, err := c.post(ctx, "sendMessage", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result apiResponse[Message]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("telegram: decode sendMessage response: %w", err)
	}

	if !result.OK {
		return nil, fmt.Errorf("telegram: sendMessage failed: %s", result.Description)
	}

	return &result.Result, nil
}

// DeleteWebhook removes the webhook integration.
func (c *Client) DeleteWebhook(ctx context.Context) error {
	resp, err := c.post(ctx, "deleteWebhook", struct{}{})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telegram: deleteWebhook failed: %s", resp.Status)
	}
	return nil
}
