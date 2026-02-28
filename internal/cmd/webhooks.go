package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/builtbyrobben/n8n-cli/internal/outfmt"
)

type WebhooksCmd struct {
	Trigger WebhooksTriggerCmd `cmd:"" help:"Trigger a webhook"`
}

type WebhooksTriggerCmd struct {
	Path   string `required:"" help:"Webhook path"`
	Method string `help:"HTTP method" default:"POST"`
	Data   string `help:"JSON data to send"`
}

func (cmd *WebhooksTriggerCmd) Run(ctx context.Context) error {
	client, err := getN8NClient()
	if err != nil {
		return err
	}

	var data any

	if cmd.Data != "" {
		if err := json.Unmarshal([]byte(cmd.Data), &data); err != nil {
			return fmt.Errorf("parse data JSON: %w", err)
		}
	}

	result, err := client.TriggerWebhook(ctx, cmd.Path, cmd.Method, data)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	if outfmt.IsPlain(ctx) {
		return outfmt.WritePlain(os.Stdout, []string{"STATUS", "PATH"}, [][]string{{"success", cmd.Path}})
	}

	if result != nil {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	fmt.Fprintf(os.Stderr, "Webhook %s triggered\n", cmd.Path)

	return nil
}
