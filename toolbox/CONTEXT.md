# CONTEXT.md

## Project Overview

- **Purpose**: A minimal, tool-focused AI assistant that only performs actions through explicitly defined tools
- **Vision**: Create a secure, constrained agent system where capabilities are controlled through tool access, not broad system permissions
- **Status**: Initial implementation phase

## Architecture Decisions

### Tech Stack
- **Language**: Go (stdlib-first, minimal dependencies)
- **Frontend**: Vanilla HTML/CSS/JS (no frameworks needed for simple text box UI)
- **Database**: PostgreSQL (via docker-compose)
- **Deployment**: Docker Compose (local development, easy to extend)
- **LLM**: Provider-agnostic (OpenAI, Anthropic, local models via API)

**Why Go?**
- Stdlib has everything we need (HTTP server, JSON, templates)
- Fast, single binary deployment
- Excellent docker support
- Type safety for tool definitions

### Design Patterns
- **Tool Registry Pattern**: Tools register themselves with capabilities and schemas
- **Provider Pattern**: LLM providers implement common interface (can swap Claude/GPT/local)
- **MCP-inspired**: Tools similar to Model Context Protocol, but simpler

### Trade-offs
- **No OpenCode fork**: OpenCode is TypeScript/Node-based with 100k+ lines. Building from scratch in Go gives us:
  - Full control over tool execution model
  - Minimal dependencies (stdlib > npm packages)
  - Clear security boundary (tools = explicit permissions)
  - Simpler to reason about and maintain
  
- **Web UI over TUI**: Unlike OpenCode's terminal focus, we prioritize:
  - Browser-based access (works from any device)
  - Easier to build data visualization from databases
  - No terminal emulator requirements

## Current State

### What's Working
- ✅ Go module initialized
- ✅ Directory structure created
- ✅ Project architecture documented

### What's Not
- ❌ No code written yet
- ❌ No docker infrastructure
- ❌ No LLM integration
- ❌ No tools implemented

### Next Steps
1. Create docker-compose.yml (PostgreSQL + toolbox service)
2. Implement tool registry system
3. Create basic web server with HTML form
4. Add first tool (database query)
5. Integrate LLM provider (start with OpenAI/Anthropic API)
6. Add file server tool
7. Add docker container management tool

## Task Tracking

### Available
- [ ] Add file server tool (read/write files in controlled directory)
- [ ] Add docker tool (inspect/manage containers)
- [ ] Add authentication/authorization
- [ ] Add conversation history persistence
- [ ] Add streaming responses
- [ ] Add tool permission system (require approval for dangerous ops)

### In Progress
- [ ] Writing CONTEXT.md (this file)

### Completed
- [x] Design architecture and directory structure
- [x] Create Go module

## Development Notes

### Setup Instructions
```bash
cd toolbox/
docker-compose up -d    # Start PostgreSQL
go run cmd/server/main.go

# Access UI at: http://localhost:8080
```

### Testing Strategy
- Unit tests for tool registry
- Integration tests for tool execution
- Table-driven tests (Go standard)
- 95% coverage target (per AGENTS.md)

### Tool Development Pattern
Tools must implement:
```go
type Tool interface {
    Name() string
    Description() string
    Schema() map[string]interface{}  // JSON schema for parameters
    Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)
}
```

### Security Model
- **Principle**: Tools define capability boundaries
- **No shell access**: No bash/shell tools by default
- **Explicit permissions**: Each tool declares what it can access
- **Read-only first**: Start with read-only tools, add write carefully
- **Database isolation**: Tools use connection pools, no raw SQL from user

## Open Questions
- Which LLM provider to use initially? (OpenAI for development, allow config)
- Should tools be hot-reloadable?
- Store conversation history in PostgreSQL or separate store?
- Support tool composition (tools calling tools)?
- Web UI framework or vanilla JS? (Start vanilla, add framework if needed)
