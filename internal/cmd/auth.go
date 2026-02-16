package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/builtbyrobben/n8n-cli/internal/outfmt"
	"github.com/builtbyrobben/n8n-cli/internal/secrets"
)

type AuthCmd struct {
	SetKey AuthSetKeyCmd `cmd:"" help:"Set API key and n8n URL"`
	Status AuthStatusCmd `cmd:"" help:"Show authentication status"`
	Remove AuthRemoveCmd `cmd:"" help:"Remove stored credentials"`
}

type AuthSetKeyCmd struct {
	URL string `required:"" help:"n8n instance URL (e.g. https://n8n.example.com)"`
	Key string `arg:"" optional:"" help:"API key (discouraged; exposes in shell history)"`
}

func (cmd *AuthSetKeyCmd) Run(ctx context.Context) error {
	var apiKey string

	switch {
	case cmd.Key != "":
		fmt.Fprintln(os.Stderr, "Warning: passing keys as arguments exposes them in shell history.")
		apiKey = strings.TrimSpace(cmd.Key)
	case term.IsTerminal(int(os.Stdin.Fd())):
		fmt.Fprint(os.Stderr, "Enter API key: ")

		byteKey, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(os.Stderr)

		if err != nil {
			return fmt.Errorf("read API key: %w", err)
		}

		apiKey = strings.TrimSpace(string(byteKey))
	default:
		byteKey, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("read API key from stdin: %w", err)
		}

		apiKey = strings.TrimSpace(string(byteKey))
	}

	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	apiURL := strings.TrimRight(strings.TrimSpace(cmd.URL), "/")
	if apiURL == "" {
		return fmt.Errorf("API URL cannot be empty")
	}

	store, err := secrets.OpenDefault()
	if err != nil {
		return fmt.Errorf("open credential store: %w", err)
	}

	if err := store.SetAPIKey(apiKey); err != nil {
		return fmt.Errorf("store API key: %w", err)
	}

	if err := store.SetAPIURL(apiURL); err != nil {
		return fmt.Errorf("store API URL: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]string{
			"status":  "success",
			"message": "Credentials stored in keyring",
			"url":     apiURL,
		})
	}

	if outfmt.IsPlain(ctx) {
		return outfmt.WritePlain(os.Stdout, []string{"STATUS", "MESSAGE", "URL"}, [][]string{{"success", "Credentials stored in keyring", apiURL}})
	}

	fmt.Fprintln(os.Stderr, "Credentials stored in keyring")

	return nil
}

type AuthStatusCmd struct{}

func (cmd *AuthStatusCmd) Run(ctx context.Context) error {
	store, err := secrets.OpenDefault()
	if err != nil {
		return fmt.Errorf("open credential store: %w", err)
	}

	hasKey, err := store.HasKey()
	if err != nil {
		return fmt.Errorf("check API key: %w", err)
	}

	envKey := os.Getenv("N8N_API_KEY")
	envURL := os.Getenv("N8N_URL")
	envOverride := envKey != "" && envURL != ""

	status := map[string]any{
		"has_key":         hasKey,
		"env_override":    envOverride,
		"storage_backend": "keyring",
	}

	if hasKey && !envOverride {
		key, err := store.GetAPIKey()
		if err == nil && len(key) > 8 {
			status["key_redacted"] = key[:4] + "..." + key[len(key)-4:]
		}

		url, err := store.GetAPIURL()
		if err == nil && url != "" {
			status["url"] = url
		}
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, status)
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"HAS_KEY", "ENV_OVERRIDE", "STORAGE"}
		rows := [][]string{{
			fmt.Sprintf("%t", hasKey),
			fmt.Sprintf("%t", envOverride),
			"keyring",
		}}

		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	fmt.Fprintf(os.Stdout, "Storage: %s\n", status["storage_backend"])

	switch {
	case envOverride:
		fmt.Fprintln(os.Stdout, "API Key: Using N8N_API_KEY + N8N_URL environment variables")
	case hasKey:
		fmt.Fprintln(os.Stdout, "API Key: Authenticated")

		if redacted, ok := status["key_redacted"].(string); ok {
			fmt.Fprintf(os.Stdout, "Key: %s\n", redacted)
		}

		if url, ok := status["url"].(string); ok {
			fmt.Fprintf(os.Stdout, "URL: %s\n", url)
		}
	default:
		fmt.Fprintln(os.Stdout, "API Key: Not authenticated")
		fmt.Fprintln(os.Stderr, "Run: n8n-cli auth set-key --url <url>")
	}

	return nil
}

type AuthRemoveCmd struct{}

func (cmd *AuthRemoveCmd) Run(ctx context.Context) error {
	store, err := secrets.OpenDefault()
	if err != nil {
		return fmt.Errorf("open credential store: %w", err)
	}

	if err := store.DeleteAll(); err != nil {
		return fmt.Errorf("remove credentials: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]string{
			"status":  "success",
			"message": "Credentials removed",
		})
	}

	if outfmt.IsPlain(ctx) {
		return outfmt.WritePlain(os.Stdout, []string{"STATUS", "MESSAGE"}, [][]string{{"success", "Credentials removed"}})
	}

	fmt.Fprintln(os.Stderr, "Credentials removed")

	return nil
}
