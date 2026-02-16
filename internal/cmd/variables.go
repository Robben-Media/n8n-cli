package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/builtbyrobben/n8n-cli/internal/outfmt"
)

type VariablesCmd struct {
	List   VariablesListCmd   `cmd:"" help:"List variables"`
	Create VariablesCreateCmd `cmd:"" help:"Create a variable"`
	Delete VariablesDeleteCmd `cmd:"" help:"Delete a variable"`
}

type VariablesListCmd struct{}

func (cmd *VariablesListCmd) Run(ctx context.Context) error {
	client, err := getN8NClient()
	if err != nil {
		return err
	}

	variables, err := client.ListVariables(ctx)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, variables)
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"ID", "KEY", "VALUE"}
		var rows [][]string

		for _, v := range variables {
			rows = append(rows, []string{v.ID, v.Key, v.Value})
		}

		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	if len(variables) == 0 {
		fmt.Fprintln(os.Stderr, "No variables found")
		return nil
	}

	fmt.Fprintf(os.Stderr, "Found %d variables\n\n", len(variables))

	for _, v := range variables {
		fmt.Printf("%s  %s = %s\n", v.ID, v.Key, v.Value)
	}

	return nil
}

type VariablesCreateCmd struct {
	Key   string `required:"" help:"Variable key"`
	Value string `required:"" help:"Variable value"`
}

func (cmd *VariablesCreateCmd) Run(ctx context.Context) error {
	client, err := getN8NClient()
	if err != nil {
		return err
	}

	variable, err := client.CreateVariable(ctx, cmd.Key, cmd.Value)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, variable)
	}

	if outfmt.IsPlain(ctx) {
		return outfmt.WritePlain(os.Stdout, []string{"ID", "KEY", "VALUE"}, [][]string{{variable.ID, variable.Key, variable.Value}})
	}

	fmt.Fprintf(os.Stderr, "Variable created: %s = %s (%s)\n", variable.Key, variable.Value, variable.ID)

	return nil
}

type VariablesDeleteCmd struct {
	ID string `arg:"" required:"" help:"Variable ID"`
}

func (cmd *VariablesDeleteCmd) Run(ctx context.Context) error {
	client, err := getN8NClient()
	if err != nil {
		return err
	}

	if err := client.DeleteVariable(ctx, cmd.ID); err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]string{
			"status":  "success",
			"message": "Variable deleted",
			"id":      cmd.ID,
		})
	}

	if outfmt.IsPlain(ctx) {
		return outfmt.WritePlain(os.Stdout, []string{"STATUS", "ID"}, [][]string{{"success", cmd.ID}})
	}

	fmt.Fprintf(os.Stderr, "Variable %s deleted\n", cmd.ID)

	return nil
}
