package providers

import "context"

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ToolDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"`
}

type ToolCall struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ChatResponse struct {
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type AIClient interface {
	Chat(ctx context.Context, messages []ChatMessage) (*ChatResponse, error)
	ChatWithTools(ctx context.Context, messages []ChatMessage, tools []ToolDefinition) (*ChatResponse, error)
	Generate(ctx context.Context, systemPrompt string, userPrompt string) (string, error)
}
