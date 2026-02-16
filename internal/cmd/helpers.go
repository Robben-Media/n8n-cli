package cmd

import (
	"fmt"
	"os"

	"github.com/builtbyrobben/n8n-cli/internal/n8n"
	"github.com/builtbyrobben/n8n-cli/internal/secrets"
)

func getN8NClient() (*n8n.Client, error) {
	apiURL := os.Getenv("N8N_URL")
	apiKey := os.Getenv("N8N_API_KEY")

	if apiURL != "" && apiKey != "" {
		return n8n.NewClient(apiURL, apiKey), nil
	}

	store, err := secrets.OpenDefault()
	if err != nil {
		return nil, fmt.Errorf("open credential store: %w", err)
	}

	key, _ := store.GetAPIKey()
	url, _ := store.GetAPIURL()

	if key == "" || url == "" {
		return nil, fmt.Errorf("no credentials found; run: n8n-cli auth set-key --url <url>")
	}

	return n8n.NewClient(url, key), nil
}
