package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/builtbyrobben/n8n-cli/internal/outfmt"
)

type CredentialsCmd struct {
	List CredentialsListCmd `cmd:"" help:"List credentials"`
}

type CredentialsListCmd struct{}

func (cmd *CredentialsListCmd) Run(ctx context.Context) error {
	client, err := getN8NClient()
	if err != nil {
		return err
	}

	credentials, err := client.ListCredentials(ctx)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, credentials)
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"ID", "NAME", "TYPE", "CREATED_AT"}
		var rows [][]string

		for _, c := range credentials {
			rows = append(rows, []string{c.ID, c.Name, c.Type, c.CreatedAt.Format("2006-01-02 15:04:05")})
		}

		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	if len(credentials) == 0 {
		fmt.Fprintln(os.Stderr, "No credentials found")
		return nil
	}

	fmt.Fprintf(os.Stderr, "Found %d credentials\n\n", len(credentials))

	for _, c := range credentials {
		fmt.Printf("%s  %s\n", c.ID, c.Name)
		fmt.Printf("  Type: %s\n", c.Type)
		fmt.Printf("  Created: %s\n\n", c.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	return nil
}
