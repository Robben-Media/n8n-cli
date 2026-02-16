package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/builtbyrobben/n8n-cli/internal/outfmt"
)

type ExecutionsCmd struct {
	List   ExecutionsListCmd   `cmd:"" help:"List executions"`
	Get    ExecutionsGetCmd    `cmd:"" help:"Get an execution by ID"`
	Delete ExecutionsDeleteCmd `cmd:"" help:"Delete an execution"`
	Retry  ExecutionsRetryCmd  `cmd:"" help:"Retry a failed execution"`
}

type ExecutionsListCmd struct {
	WorkflowID string `help:"Filter by workflow ID"`
	Status     string `help:"Filter by status (success, error, waiting)"`
	Limit      int    `help:"Maximum number of executions to return" default:"20"`
	Cursor     string `help:"Pagination cursor"`
}

func (cmd *ExecutionsListCmd) Run(ctx context.Context) error {
	client, err := getN8NClient()
	if err != nil {
		return err
	}

	result, err := client.ListExecutions(ctx, cmd.WorkflowID, cmd.Status, cmd.Limit, cmd.Cursor)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"ID", "WORKFLOW_ID", "STATUS", "MODE", "STARTED_AT"}
		var rows [][]string

		for _, e := range result.Data {
			rows = append(rows, []string{e.ID, e.WorkflowID, e.Status, e.Mode, e.StartedAt.Format("2006-01-02 15:04:05")})
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
		fmt.Fprintln(os.Stderr, "No executions found")
		return nil
	}

	fmt.Fprintf(os.Stderr, "Found %d executions\n\n", len(result.Data))

	for _, e := range result.Data {
		fmt.Printf("%s  workflow=%s\n", e.ID, e.WorkflowID)
		fmt.Printf("  Status: %s  Mode: %s\n", e.Status, e.Mode)
		fmt.Printf("  Started: %s\n\n", e.StartedAt.Format("2006-01-02 15:04:05"))
	}

	if result.NextCursor != "" {
		fmt.Fprintf(os.Stderr, "Next cursor: %s\n", result.NextCursor)
	}

	return nil
}

type ExecutionsGetCmd struct {
	ID string `arg:"" required:"" help:"Execution ID"`
}

func (cmd *ExecutionsGetCmd) Run(ctx context.Context) error {
	client, err := getN8NClient()
	if err != nil {
		return err
	}

	execution, err := client.GetExecution(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, execution)
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"ID", "WORKFLOW_ID", "STATUS", "MODE", "FINISHED", "STARTED_AT", "STOPPED_AT"}
		rows := [][]string{{
			execution.ID, execution.WorkflowID, execution.Status, execution.Mode,
			fmt.Sprintf("%t", execution.Finished),
			execution.StartedAt.Format("2006-01-02 15:04:05"),
			execution.StoppedAt.Format("2006-01-02 15:04:05"),
		}}

		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	fmt.Printf("ID: %s\n", execution.ID)
	fmt.Printf("Workflow ID: %s\n", execution.WorkflowID)
	fmt.Printf("Status: %s\n", execution.Status)
	fmt.Printf("Mode: %s\n", execution.Mode)
	fmt.Printf("Finished: %t\n", execution.Finished)
	fmt.Printf("Started: %s\n", execution.StartedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Stopped: %s\n", execution.StoppedAt.Format("2006-01-02 15:04:05"))

	return nil
}

type ExecutionsDeleteCmd struct {
	ID string `arg:"" required:"" help:"Execution ID"`
}

func (cmd *ExecutionsDeleteCmd) Run(ctx context.Context) error {
	client, err := getN8NClient()
	if err != nil {
		return err
	}

	if err := client.DeleteExecution(ctx, cmd.ID); err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]string{
			"status":  "success",
			"message": "Execution deleted",
			"id":      cmd.ID,
		})
	}

	if outfmt.IsPlain(ctx) {
		return outfmt.WritePlain(os.Stdout, []string{"STATUS", "ID"}, [][]string{{"success", cmd.ID}})
	}

	fmt.Fprintf(os.Stderr, "Execution %s deleted\n", cmd.ID)

	return nil
}

type ExecutionsRetryCmd struct {
	ID string `arg:"" required:"" help:"Execution ID"`
}

func (cmd *ExecutionsRetryCmd) Run(ctx context.Context) error {
	client, err := getN8NClient()
	if err != nil {
		return err
	}

	if err := client.RetryExecution(ctx, cmd.ID); err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]string{
			"status":  "success",
			"message": "Execution retry initiated",
			"id":      cmd.ID,
		})
	}

	if outfmt.IsPlain(ctx) {
		return outfmt.WritePlain(os.Stdout, []string{"STATUS", "ID"}, [][]string{{"success", cmd.ID}})
	}

	fmt.Fprintf(os.Stderr, "Execution %s retry initiated\n", cmd.ID)

	return nil
}
