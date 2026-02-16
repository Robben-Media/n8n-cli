package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhooksTriggerRun_UsesConfiguredMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/webhook/test-path" {
			t.Fatalf("expected path /webhook/test-path, got %s", r.URL.Path)
		}

		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	t.Setenv("N8N_URL", server.URL)
	t.Setenv("N8N_API_KEY", "test-api-key")

	cmd := WebhooksTriggerCmd{
		Path:   "test-path",
		Method: http.MethodGet,
	}

	if err := cmd.Run(context.Background()); err != nil {
		t.Fatalf("expected webhook trigger to succeed with GET method, got error: %v", err)
	}
}
