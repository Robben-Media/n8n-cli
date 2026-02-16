package n8n

import "time"

type Workflow struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Tags        []Tag     `json:"tags,omitempty"`
	Nodes       any       `json:"nodes,omitempty"`
	Connections any       `json:"connections,omitempty"`
}

type WorkflowList struct {
	Data       []Workflow `json:"data"`
	NextCursor string     `json:"nextCursor"`
}

type Execution struct {
	ID           string    `json:"id"`
	Finished     bool      `json:"finished"`
	Mode         string    `json:"mode"`
	Status       string    `json:"status"`
	StartedAt    time.Time `json:"startedAt"`
	StoppedAt    time.Time `json:"stoppedAt"`
	WorkflowID   string    `json:"workflowId"`
	WorkflowData any       `json:"workflowData,omitempty"`
}

type ExecutionList struct {
	Data       []Execution `json:"data"`
	NextCursor string      `json:"nextCursor"`
}

type Credential struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Tag struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Variable struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type HealthStatus struct {
	Status string `json:"status"`
}
