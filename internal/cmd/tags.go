package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/builtbyrobben/n8n-cli/internal/outfmt"
)

type TagsCmd struct {
	List   TagsListCmd   `cmd:"" help:"List tags"`
	Create TagsCreateCmd `cmd:"" help:"Create a tag"`
	Delete TagsDeleteCmd `cmd:"" help:"Delete a tag"`
}

type TagsListCmd struct{}

func (cmd *TagsListCmd) Run(ctx context.Context) error {
	client, err := getN8NClient()
	if err != nil {
		return err
	}

	tags, err := client.ListTags(ctx)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, tags)
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"ID", "NAME", "CREATED_AT"}
		var rows [][]string

		for _, t := range tags {
			rows = append(rows, []string{t.ID, t.Name, t.CreatedAt.Format("2006-01-02 15:04:05")})
		}

		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	if len(tags) == 0 {
		fmt.Fprintln(os.Stderr, "No tags found")
		return nil
	}

	fmt.Fprintf(os.Stderr, "Found %d tags\n\n", len(tags))

	for _, t := range tags {
		fmt.Printf("%s  %s\n", t.ID, t.Name)
	}

	return nil
}

type TagsCreateCmd struct {
	Name string `arg:"" required:"" help:"Tag name"`
}

func (cmd *TagsCreateCmd) Run(ctx context.Context) error {
	client, err := getN8NClient()
	if err != nil {
		return err
	}

	tag, err := client.CreateTag(ctx, cmd.Name)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, tag)
	}

	if outfmt.IsPlain(ctx) {
		return outfmt.WritePlain(os.Stdout, []string{"ID", "NAME"}, [][]string{{tag.ID, tag.Name}})
	}

	fmt.Fprintf(os.Stderr, "Tag created: %s (%s)\n", tag.Name, tag.ID)

	return nil
}

type TagsDeleteCmd struct {
	ID string `arg:"" required:"" help:"Tag ID"`
}

func (cmd *TagsDeleteCmd) Run(ctx context.Context) error {
	client, err := getN8NClient()
	if err != nil {
		return err
	}

	if err := client.DeleteTag(ctx, cmd.ID); err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]string{
			"status":  "success",
			"message": "Tag deleted",
			"id":      cmd.ID,
		})
	}

	if outfmt.IsPlain(ctx) {
		return outfmt.WritePlain(os.Stdout, []string{"STATUS", "ID"}, [][]string{{"success", cmd.ID}})
	}

	fmt.Fprintf(os.Stderr, "Tag %s deleted\n", cmd.ID)

	return nil
}
