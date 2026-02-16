package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGet_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		if r.Header.Get("X-N8N-API-KEY") != "test-key" {
			t.Errorf("expected X-N8N-API-KEY header 'test-key', got %s", r.Header.Get("X-N8N-API-KEY"))
		}

		resp := map[string]any{
			"data": []map[string]any{
				{"id": "1", "name": "Test Workflow"},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, WithAuthFn(func(r *http.Request) {
		r.Header.Set("X-N8N-API-KEY", "test-key")
	}))

	var result struct {
		Data []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"data"`
	}

	err := client.Get(context.Background(), "/workflows", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Data) != 1 {
		t.Fatalf("expected 1 workflow, got %d", len(result.Data))
	}

	if result.Data[0].ID != "1" {
		t.Errorf("expected ID 1, got %s", result.Data[0].ID)
	}

	if result.Data[0].Name != "Test Workflow" {
		t.Errorf("expected Name Test Workflow, got %s", result.Data[0].Name)
	}
}

func TestPost_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)

		if body["name"] != "my-tag" {
			t.Errorf("expected name my-tag, got %s", body["name"])
		}

		resp := map[string]any{
			"id":   "tag-1",
			"name": "my-tag",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)

	var result struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	err := client.Post(context.Background(), "/tags", map[string]string{"name": "my-tag"}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != "tag-1" {
		t.Errorf("expected ID tag-1, got %s", result.ID)
	}
}

func TestGet_HTTPError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid API key"})
	}))
	defer server.Close()

	client := NewClient(server.URL)

	var result struct{}

	err := client.Get(context.Background(), "/workflows", nil, &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if apiErr.StatusCode != 401 {
		t.Errorf("expected status 401, got %d", apiErr.StatusCode)
	}

	if apiErr.Message != "Invalid API key" {
		t.Errorf("expected message 'Invalid API key', got %s", apiErr.Message)
	}
}

func TestDelete_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(server.URL)

	err := client.Delete(context.Background(), "/workflows/1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPatch_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)

		if body["active"] != true {
			t.Errorf("expected active=true, got %v", body["active"])
		}

		resp := map[string]any{
			"id":     "1",
			"name":   "Test",
			"active": true,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)

	var result struct {
		ID     string `json:"id"`
		Active bool   `json:"active"`
	}

	err := client.Patch(context.Background(), "/workflows/1", map[string]any{"active": true}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.Active {
		t.Error("expected active=true")
	}
}

func TestGet_WithQueryParams(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("active") != "true" {
			t.Errorf("expected query param active=true, got %s", r.URL.Query().Get("active"))
		}

		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("expected query param limit=10, got %s", r.URL.Query().Get("limit"))
		}

		resp := map[string]any{"data": []any{}}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)

	var result struct {
		Data []any `json:"data"`
	}

	q := make(map[string][]string)
	q["active"] = []string{"true"}
	q["limit"] = []string{"10"}

	err := client.Get(context.Background(), "/workflows", q, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
