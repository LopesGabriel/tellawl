package database

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/errx"
)

func TestFindByID_Success(t *testing.T) {
	member := map[string]interface{}{
		"id":         "123",
		"first_name": "John",
		"last_name":  "Doe",
		"email":      "john@example.com",
		"created_at": time.Now().Format(time.RFC3339),
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/internal/members/123" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(member)
	}))
	defer s.Close()

	r, err := NewHTTPMemberRepository(s.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	m, err := r.FindByID(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Id != "123" || m.Email != "john@example.com" {
		t.Fatalf("unexpected member: %+v", m)
	}
}

func TestFindByID_NotFound(t *testing.T) {
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	r, err := NewHTTPMemberRepository(s.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = r.FindByID(context.Background(), "nope")
	if !errors.Is(err, errx.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestFindByEmail_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"id":         "abc",
				"first_name": "Jane",
				"last_name":  "Roe",
				"email":      "jane@example.com",
				"created_at": time.Now().Format(time.RFC3339),
			},
		},
		"meta": map[string]interface{}{"total": 1},
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/internal/members" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(payload)
	}))
	defer s.Close()

	r, err := NewHTTPMemberRepository(s.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	m, err := r.FindByEmail(context.Background(), "jane@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Id != "abc" || m.Email != "jane@example.com" {
		t.Fatalf("unexpected member: %+v", m)
	}
}

func TestValidateToken(t *testing.T) {
	// valid token returns 200
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/public/me" && r.Header.Get("Authorization") == "Bearer valid" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"id": "x"})
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer s.Close()

	r, err := NewHTTPMemberRepository(s.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	member, err := r.ValidateToken(context.Background(), "valid")

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if member == nil {
		t.Fatalf("expected member, got nil")
	}

	member, err = r.ValidateToken(context.Background(), "invalid")

	if err != errx.ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
