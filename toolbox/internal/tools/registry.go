// Package tools provides the tool registry and execution framework
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// Tool defines the interface that all tools must implement
type Tool interface {
	// Name returns the unique identifier for this tool
	Name() string

	// Description returns a human-readable description of what this tool does
	Description() string

	// Schema returns the JSON schema for the tool's parameters
	Schema() map[string]interface{}

	// Execute runs the tool with the provided parameters
	Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)
}

// Registry manages available tools
type Registry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

// NewRegistry creates a new tool registry
func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

// Register adds a tool to the registry
func (r *Registry) Register(tool Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := tool.Name()
	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("tool %q already registered", name)
	}

	r.tools[name] = tool
	return nil
}

// Get retrieves a tool by name
func (r *Registry) Get(name string) (Tool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tool, exists := r.tools[name]
	if !exists {
		return nil, fmt.Errorf("tool %q not found", name)
	}

	return tool, nil
}

// List returns all registered tools
func (r *Registry) List() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}

	return tools
}

// ToolInfo contains metadata about a tool for LLM consumption
type ToolInfo struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema"`
}

// ListInfo returns tool information suitable for sending to an LLM
func (r *Registry) ListInfo() []ToolInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	info := make([]ToolInfo, 0, len(r.tools))
	for _, tool := range r.tools {
		info = append(info, ToolInfo{
			Name:        tool.Name(),
			Description: tool.Description(),
			Schema:      tool.Schema(),
		})
	}

	return info
}

// Execute runs a tool by name with the provided parameters
func (r *Registry) Execute(ctx context.Context, name string, params map[string]interface{}) (interface{}, error) {
	tool, err := r.Get(name)
	if err != nil {
		return nil, err
	}

	return tool.Execute(ctx, params)
}

// MarshalJSON serializes tool info for JSON responses
func (ti ToolInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		Schema      map[string]interface{} `json:"parameters"`
	}{
		Name:        ti.Name,
		Description: ti.Description,
		Schema:      ti.Schema,
	})
}
