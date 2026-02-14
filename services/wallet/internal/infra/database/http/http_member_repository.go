package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/lopesgabriel/tellawl/packages/tracing"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/errx"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// HTTPMemberRepository implements member repository by calling member-service HTTP API.
// baseURL is private and holds the address of the member-service (e.g. http://member-service:8080).
type HTTPMemberRepository struct {
	baseURL *url.URL
	client  *http.Client
	tracer  trace.Tracer
}

// NewHTTPMemberRepository creates a new HTTP member repository. base must be a valid URL and will be normalized.
func NewHTTPMemberRepository(base string, client *http.Client) (*HTTPMemberRepository, error) {
	if client == nil {
		client = http.DefaultClient
	}

	u, err := url.Parse(strings.TrimRight(base, "/"))
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}

	return &HTTPMemberRepository{
		baseURL: u,
		client:  client,
		tracer:  tracing.GetTracer("github.com/lopesgabriel/tellawl/services/wallet/internal/infra/database/http/HTTPMemberRepository"),
	}, nil
}

// memberAPIResponse mirrors the API response from member-service internal endpoints.
type memberAPIResponse struct {
	Id        string     `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// membersListResponse models the /internal/members list response format.
type membersListResponse struct {
	Data []memberAPIResponse `json:"data"`
	Meta struct {
		Total int `json:"total"`
	} `json:"meta"`
}

// toModel converts API response to domain model.
func toModel(a memberAPIResponse) *models.Member {
	return &models.Member{
		Id:        a.Id,
		FirstName: a.FirstName,
		LastName:  a.LastName,
		Email:     a.Email,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func (r *HTTPMemberRepository) FindByID(ctx context.Context, id string) (*models.Member, error) {
	ctx, span := r.tracer.Start(ctx, "FindByID")
	defer span.End()

	urlStr := r.baseURL.String() + path.Join("/internal/members/", id)
	span.SetAttributes(attribute.String("http.url", urlStr), attribute.String("member.id", id))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		span.SetStatus(codes.Error, "member not found")
		return nil, errx.ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	var m memberAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "member found")
	return toModel(m), nil
}

func (r *HTTPMemberRepository) FindByEmail(ctx context.Context, email string) (*models.Member, error) {
	ctx, span := r.tracer.Start(ctx, "FindByEmail")
	defer span.End()

	q := url.Values{}
	q.Set("email", email)

	urlStr := fmt.Sprintf("%s%s?%s", r.baseURL.String(), "/internal/members", q.Encode())
	span.SetAttributes(attribute.String("http.url", urlStr), attribute.String("member.email", email))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		span.SetStatus(codes.Error, "member not found")
		return nil, errx.ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	var list membersListResponse
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if len(list.Data) == 0 {
		span.SetStatus(codes.Error, "member not found")
		return nil, errx.ErrNotFound
	}

	span.SetStatus(codes.Ok, "member found")
	return toModel(list.Data[0]), nil
}

func (r *HTTPMemberRepository) ValidateToken(ctx context.Context, token string) (*models.Member, error) {
	ctx, span := r.tracer.Start(ctx, "ValidateToken")
	defer span.End()

	urlStr := r.baseURL.String() + "/public/me"
	span.SetAttributes(attribute.String("http.url", urlStr))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := r.client.Do(req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		span.SetStatus(codes.Error, "invalid credentials")
		return nil, errx.ErrInvalidCredentials
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	var m memberAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "valid token")
	return toModel(m), nil
}
