package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/builtbyrobben/n8n-cli/internal/outfmt"
)

type WorkflowsCmd struct {
	List       WorkflowsListCmd       `cmd:"" help:"List workflows"`
	Get        WorkflowsGetCmd        `cmd:"" help:"Get a workflow by ID"`
	Activate   WorkflowsActivateCmd   `cmd:"" help:"Activate a workflow"`
	Deactivate WorkflowsDeactivateCmd `cmd:"" help:"Deactivate a workflow"`
	Delete     WorkflowsDeleteCmd     `cmd:"" help:"Delete a workflow"`
}

type WorkflowsListCmd struct {
	Active *bool  `help:"Filter by active status"`
	Tags   string `help:"Filter by tag name"`
	Limit  int    `help:"Maximum number of workflows to return" default:"20"`
	Cursor string `help:"Pagination cursor"`
}

func (cmd *WorkflowsListCmd) Run(ctx context.Context) error {
	client, err := getN8NClient()
	if err != nil {
		return err
	}

	result, err := client.ListWorkflows(ctx, cmd.Active, cmd.Tags, cmd.Limit, cmd.Cursor)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"ID", "NAME", "ACTIVE", "UPDATED_AT"}
		var rows [][]string

		for _, w := range result.Data {
			rows = append(rows, []string{w.ID, w.Name, fmt.Sprintf("%t", w.Active), w.UpdatedAt.Format("2006-01-02 15:04:05")})
		}

		if err := outfmt.WritePlain(os.Stdout, headers, rows); err != nil {
			return err
		}

		if result.NextCursor != "" {
			fmt.Fprintf(os.Stderr, "Next cursor: %s\n", result.NextCursor)
		}

		return nil
	}

	if len(result.Data) == 0 {
		fmt.Fprintln(os.Stderr, "No workflows found")
		return nil
	}

	fmt.Fprintf(os.Stderr, "Found %d workflows\n\n", len(result.Data))

	for _, w := range result.Data {
		fmt.Printf("%s  %s\n", w.ID, w.Name)
		fmt.Printf("  Active: %t\n", w.Active)
		fmt.Printf("  Updated: %s\n\n", w.UpdatedAt.Format("2006-01-02 15:04:05"))
	}

	if result.NextCursor != "" {
		fmt.Fprintf(os.Stderr, "Next cursor: %s\n", result.NextCursor)
	}

	return nil
}

type WorkflowsGetCmd struct {
	ID string `arg:"" required:"" help:"Workflow ID"`
}

func (cmd *WorkflowsGetCmd) Run(ctx context.Context) error {
	client, err := getN8NClient()
	if err != nil {
		return err
	}

	workflow, err := client.GetWorkflow(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, workflow)
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"ID", "NAME", "ACTIVE", "CREATED_AT", "UPDATED_AT"}
		rows := [][]string{{
			workflow.ID, workflow.Name,
			fmt.Sprintf("%t", workflow.Active),
			workflow.CreatedAt.Format("2006-01-02 15:04:05"),
			workflow.UpdatedAt.Format("2006-01-02 15:04:05"),
		}}

		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	fmt.Printf("ID: %s\n", workflow.ID)
	fmt.Printf("Name: %s\n", workflow.Name)
	fmt.Printf("Active: %t\n", workflow.Active)
	fmt.Printf("Created: %s\n", workflow.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated: %s\n", workflow.UpdatedAt.Format("2006-01-02 15:04:05"))

	return nil
}

type WorkflowsActivateCmd struct {
	ID string `arg:"" required:"" help:"Workflow ID"`
}

func (cmd *WorkflowsActivateCmd) Run(ctx context.Context) error {
	client, err := getN8NClient()
	if err != nil {
		return err
	}

	workflow, err := client.ActivateWorkflow(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]any{
			"status":  "success",
			"message": "Workflow activated",
			"id":      workflow.ID,
			"active":  workflow.Active,
		})
	}

	if outfmt.IsPlain(ctx) {
		return outfmt.WritePlain(os.Stdout, []string{"STATUS", "ID", "ACTIVE"}, [][]string{{"success", workflow.ID, fmt.Sprintf("%t", workflow.Active)}})
	}

	fmt.Fprintf(os.Stderr, "Workflow %s activated\n", workflow.ID)

	return nil
}

type WorkflowsDeactivateCmd struct {
	ID string `arg:"" required:"" help:"Workflow ID"`
}

func (cmd *WorkflowsDeactivateCmd) Run(ctx context.Context) error {
	client, err := getN8NClient()
	if err != nil {
		return err
	}

	workflow, err := client.DeactivateWorkflow(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]any{
			"status":  "success",
			"message": "Workflow deactivated",
			"id":      workflow.ID,
			"active":  workflow.Active,
		})
	}

	if outfmt.IsPlain(ctx) {
		return outfmt.WritePlain(os.Stdout, []string{"STATUS", "ID", "ACTIVE"}, [][]string{{"success", workflow.ID, fmt.Sprintf("%t", workflow.Active)}})
	}

	fmt.Fprintf(os.Stderr, "Workflow %s deactivated\n", workflow.ID)

	return nil
}

type WorkflowsDeleteCmd struct {
	ID string `arg:"" required:"" help:"Workflow ID"`
}

func (cmd *WorkflowsDeleteCmd) Run(ctx context.Context) error {
	client, err := getN8NClient()
	if err != nil {
		return err
	}

	if err := client.DeleteWorkflow(ctx, cmd.ID); err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]string{
			"status":  "success",
			"message": "Workflow deleted",
			"id":      cmd.ID,
		})
	}

	if outfmt.IsPlain(ctx) {
		return outfmt.WritePlain(os.Stdout, []string{"STATUS", "ID"}, [][]string{{"success", cmd.ID}})
	}

	fmt.Fprintf(os.Stderr, "Workflow %s deleted\n", cmd.ID)

	return nil
}
