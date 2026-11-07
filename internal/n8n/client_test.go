package n8n

import (
	"context"
	"errors"
	"testing"
)

func TestTriggerWebhookUnsupportedMethod(t *testing.T) {
	client := NewClient("http://127.0.0.1:1", "")

	result, err := client.TriggerWebhook(context.Background(), "example", " trace ", nil)
	if err == nil {
		t.Fatal("expected unsupported method error")
	}

	if err.Error() != "unsupported webhook method: TRACE" {
		t.Fatalf("unexpected error: %v", err)
	}

	if !errors.Is(err, errUnsupportedWebhookMethod) {
		t.Fatalf("error does not wrap unsupported method sentinel: %v", err)
	}

	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
}
