// Package agent provides a mock LLM provider for testing
package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/jredh-dev/nascent-nexus/toolbox/internal/tools"
)

// MockProvider is a simple mock LLM for testing
type MockProvider struct{}

// NewMockProvider creates a mock provider
func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

// Complete implements Provider
func (m *MockProvider) Complete(ctx context.Context, messages []Message, availableTools []tools.ToolInfo) (*Response, error) {
	// Get last user message
	var lastUserMsg string
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			lastUserMsg = strings.ToLower(messages[i].Content)
			break
		}
	}

	// Simple keyword-based mock responses
	if strings.Contains(lastUserMsg, "task") && strings.Contains(lastUserMsg, "list") {
		return &Response{
			Content: "I'll query the database for tasks",
			ToolCalls: []ToolCall{
				{
					ID:       "call_1",
					ToolName: "query_database",
					Parameters: map[string]interface{}{
						"query": "SELECT * FROM tasks ORDER BY created_at DESC",
					},
				},
			},
		}, nil
	}

	if strings.Contains(lastUserMsg, "note") {
		return &Response{
			Content: "I'll query the database for notes",
			ToolCalls: []ToolCall{
				{
					ID:       "call_2",
					ToolName: "query_database",
					Parameters: map[string]interface{}{
						"query": "SELECT * FROM notes ORDER BY created_at DESC",
					},
				},
			},
		}, nil
	}

	// Check if this is a tool result message
	if strings.Contains(lastUserMsg, "tool result") {
		return &Response{
			Content: fmt.Sprintf("Here are the results from the database query: %s", lastUserMsg),
			Done:    true,
		}, nil
	}

	// Default response with list of available tools
	toolNames := make([]string, len(availableTools))
	for i, tool := range availableTools {
		toolNames[i] = tool.Name
	}

	return &Response{
		Content: fmt.Sprintf("I'm a mock LLM provider. Available tools: %v. Try asking about 'tasks' or 'notes'.", toolNames),
		Done:    true,
	}, nil
}
