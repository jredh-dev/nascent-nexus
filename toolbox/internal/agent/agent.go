// Package agent provides LLM agent orchestration
package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jredh-dev/nascent-nexus/toolbox/internal/tools"
)

// Provider defines the interface for LLM providers (OpenAI, Anthropic, etc.)
type Provider interface {
	// Complete sends a message to the LLM and returns the response
	Complete(ctx context.Context, messages []Message, availableTools []tools.ToolInfo) (*Response, error)
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"` // "user", "assistant", "system"
	Content string `json:"content"`
}

// Response represents an LLM response
type Response struct {
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
	Done      bool       `json:"done"` // true if no more tool calls needed
}

// ToolCall represents a request to execute a tool
type ToolCall struct {
	ID         string                 `json:"id"`
	ToolName   string                 `json:"tool_name"`
	Parameters map[string]interface{} `json:"parameters"`
}

// Agent orchestrates conversations with tool execution
type Agent struct {
	provider Provider
	registry *tools.Registry
	history  []Message
}

// New creates a new agent
func New(provider Provider, registry *tools.Registry) *Agent {
	return &Agent{
		provider: provider,
		registry: registry,
		history:  make([]Message, 0),
	}
}

// Process handles a user message and executes tools as needed
func (a *Agent) Process(ctx context.Context, userMessage string) (string, error) {
	// Add user message to history
	a.history = append(a.history, Message{
		Role:    "user",
		Content: userMessage,
	})

	// Get available tools
	availableTools := a.registry.ListInfo()

	// Loop until LLM is done (handles multi-step tool usage)
	maxIterations := 10
	for i := 0; i < maxIterations; i++ {
		// Get LLM response
		resp, err := a.provider.Complete(ctx, a.history, availableTools)
		if err != nil {
			return "", fmt.Errorf("LLM completion failed: %w", err)
		}

		// If no tool calls, we're done
		if len(resp.ToolCalls) == 0 {
			a.history = append(a.history, Message{
				Role:    "assistant",
				Content: resp.Content,
			})
			return resp.Content, nil
		}

		// Execute tool calls
		toolResults := make([]string, 0, len(resp.ToolCalls))
		for _, call := range resp.ToolCalls {
			result, err := a.registry.Execute(ctx, call.ToolName, call.Parameters)
			if err != nil {
				toolResults = append(toolResults, fmt.Sprintf("Error executing %s: %v", call.ToolName, err))
				continue
			}

			// Serialize result
			resultJSON, _ := json.Marshal(result)
			toolResults = append(toolResults, fmt.Sprintf("%s result: %s", call.ToolName, string(resultJSON)))
		}

		// Add tool results to history
		a.history = append(a.history, Message{
			Role:    "assistant",
			Content: fmt.Sprintf("Tool calls: %s", resp.Content),
		})
		a.history = append(a.history, Message{
			Role:    "user",
			Content: fmt.Sprintf("Tool results: %v", toolResults),
		})
	}

	return "", fmt.Errorf("max iterations reached without completion")
}

// History returns the conversation history
func (a *Agent) History() []Message {
	return a.history
}

// Reset clears conversation history
func (a *Agent) Reset() {
	a.history = make([]Message, 0)
}
