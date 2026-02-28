package n8n

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/builtbyrobben/n8n-cli/internal/api"
)

type Client struct {
	api        *api.Client
	rawBaseURL string
}

func NewClient(baseURL, apiKey string, opts ...api.ClientOption) *Client {
	allOpts := append([]api.ClientOption{
		api.WithAuthFn(func(r *http.Request) {
			r.Header.Set("X-N8N-API-KEY", apiKey)
		}),
		api.WithUserAgent("n8n-cli/1.0"),
	}, opts...)

	return &Client{
		api:        api.NewClient(baseURL+"/api/v1", allOpts...),
		rawBaseURL: baseURL,
	}
}

// --- Workflows ---

func (c *Client) ListWorkflows(ctx context.Context, active *bool, tags string, limit int, cursor string) (*WorkflowList, error) {
	q := url.Values{}

	if active != nil {
		q.Set("active", strconv.FormatBool(*active))
	}

	if tags != "" {
		q.Set("tags", tags)
	}

	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}

	if cursor != "" {
		q.Set("cursor", cursor)
	}

	var result WorkflowList

	if err := c.api.Get(ctx, "/workflows", q, &result); err != nil {
		return nil, fmt.Errorf("list workflows: %w", err)
	}

	return &result, nil
}

func (c *Client) GetWorkflow(ctx context.Context, id string) (*Workflow, error) {
	var result Workflow

	if err := c.api.Get(ctx, "/workflows/"+id, nil, &result); err != nil {
		return nil, fmt.Errorf("get workflow: %w", err)
	}

	return &result, nil
}

func (c *Client) ActivateWorkflow(ctx context.Context, id string) (*Workflow, error) {
	var result Workflow

	body := map[string]any{"active": true}

	if err := c.api.Patch(ctx, "/workflows/"+id, body, &result); err != nil {
		return nil, fmt.Errorf("activate workflow: %w", err)
	}

	return &result, nil
}

func (c *Client) DeactivateWorkflow(ctx context.Context, id string) (*Workflow, error) {
	var result Workflow

	body := map[string]any{"active": false}

	if err := c.api.Patch(ctx, "/workflows/"+id, body, &result); err != nil {
		return nil, fmt.Errorf("deactivate workflow: %w", err)
	}

	return &result, nil
}

func (c *Client) DeleteWorkflow(ctx context.Context, id string) error {
	if err := c.api.Delete(ctx, "/workflows/"+id, nil); err != nil {
		return fmt.Errorf("delete workflow: %w", err)
	}

	return nil
}

// --- Executions ---

func (c *Client) ListExecutions(ctx context.Context, workflowID, status string, limit int, cursor string) (*ExecutionList, error) {
	q := url.Values{}

	if workflowID != "" {
		q.Set("workflowId", workflowID)
	}

	if status != "" {
		q.Set("status", status)
	}

	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}

	if cursor != "" {
		q.Set("cursor", cursor)
	}

	var result ExecutionList

	if err := c.api.Get(ctx, "/executions", q, &result); err != nil {
		return nil, fmt.Errorf("list executions: %w", err)
	}

	return &result, nil
}

func (c *Client) GetExecution(ctx context.Context, id string) (*Execution, error) {
	var result Execution

	if err := c.api.Get(ctx, "/executions/"+id, nil, &result); err != nil {
		return nil, fmt.Errorf("get execution: %w", err)
	}

	return &result, nil
}

func (c *Client) DeleteExecution(ctx context.Context, id string) error {
	if err := c.api.Delete(ctx, "/executions/"+id, nil); err != nil {
		return fmt.Errorf("delete execution: %w", err)
	}

	return nil
}

func (c *Client) RetryExecution(ctx context.Context, id string) error {
	if err := c.api.Post(ctx, "/executions/"+id+"/retry", nil, nil); err != nil {
		return fmt.Errorf("retry execution: %w", err)
	}

	return nil
}

// --- Credentials ---

func (c *Client) ListCredentials(ctx context.Context) ([]Credential, error) {
	var result struct {
		Data []Credential `json:"data"`
	}

	if err := c.api.Get(ctx, "/credentials", nil, &result); err != nil {
		return nil, fmt.Errorf("list credentials: %w", err)
	}

	return result.Data, nil
}

// --- Tags ---

func (c *Client) ListTags(ctx context.Context) ([]Tag, error) {
	var result struct {
		Data []Tag `json:"data"`
	}

	if err := c.api.Get(ctx, "/tags", nil, &result); err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}

	return result.Data, nil
}

func (c *Client) CreateTag(ctx context.Context, name string) (*Tag, error) {
	var result Tag

	body := map[string]string{"name": name}

	if err := c.api.Post(ctx, "/tags", body, &result); err != nil {
		return nil, fmt.Errorf("create tag: %w", err)
	}

	return &result, nil
}

func (c *Client) DeleteTag(ctx context.Context, id string) error {
	if err := c.api.Delete(ctx, "/tags/"+id, nil); err != nil {
		return fmt.Errorf("delete tag: %w", err)
	}

	return nil
}

// --- Variables ---

func (c *Client) ListVariables(ctx context.Context) ([]Variable, error) {
	var result struct {
		Data []Variable `json:"data"`
	}

	if err := c.api.Get(ctx, "/variables", nil, &result); err != nil {
		return nil, fmt.Errorf("list variables: %w", err)
	}

	return result.Data, nil
}

func (c *Client) CreateVariable(ctx context.Context, key, value string) (*Variable, error) {
	var result Variable

	body := map[string]string{"key": key, "value": value}

	if err := c.api.Post(ctx, "/variables", body, &result); err != nil {
		return nil, fmt.Errorf("create variable: %w", err)
	}

	return &result, nil
}

func (c *Client) DeleteVariable(ctx context.Context, id string) error {
	if err := c.api.Delete(ctx, "/variables/"+id, nil); err != nil {
		return fmt.Errorf("delete variable: %w", err)
	}

	return nil
}

// --- Webhooks ---

func (c *Client) TriggerWebhook(ctx context.Context, path, method string, data any) (any, error) {
	webhookURL := c.rawBaseURL + "/webhook/" + path

	var result any

	webhookClient := api.NewClient(webhookURL)
	httpMethod := strings.ToUpper(strings.TrimSpace(method))
	if httpMethod == "" {
		httpMethod = http.MethodPost
	}

	var err error

	switch httpMethod {
	case http.MethodGet:
		err = webhookClient.Get(ctx, "", nil, &result)
	case http.MethodPost:
		err = webhookClient.Post(ctx, "", data, &result)
	case http.MethodPut:
		err = webhookClient.Put(ctx, "", data, &result)
	case http.MethodPatch:
		err = webhookClient.Patch(ctx, "", data, &result)
	case http.MethodDelete:
		err = webhookClient.Delete(ctx, "", &result)
	default:
		return nil, fmt.Errorf("unsupported webhook method: %s", httpMethod)
	}

	if err != nil {
		return nil, fmt.Errorf("trigger webhook: %w", err)
	}

	return result, nil
}

// --- Health ---

func (c *Client) Health(ctx context.Context) (*HealthStatus, error) {
	healthClient := api.NewClient(c.rawBaseURL)

	var result HealthStatus

	if err := healthClient.Get(ctx, "/healthz", nil, &result); err != nil {
		return nil, fmt.Errorf("health check: %w", err)
	}

	return &result, nil
}
