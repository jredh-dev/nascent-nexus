# Toolbox

**A minimal, tool-focused AI assistant built from scratch in Go**

> Unlike OpenCode (TypeScript, 100k+ lines), Toolbox is a clean-room implementation with a simple goal: an AI assistant that ONLY does what you explicitly allow through tools.

## Why Toolbox?

- **Tool-constrained**: No unrestricted system access - only tools you register
- **Minimal dependencies**: Built with Go stdlib where possible
- **Docker-first**: Everything runs in docker-compose
- **AGPL-3.0**: Copyleft license ensures improvements stay open source
- **Simple architecture**: Easy to understand and extend

## Quick Start

```bash
# 1. Start infrastructure
cd toolbox/
docker-compose up -d

# 2. Run toolbox (local development)
go run cmd/server/main.go

# 3. Access web UI
open http://localhost:8080
```

## What Can It Do?

**Out of the box**, Toolbox can:
- Query your PostgreSQL database (SELECT queries only for safety)
- Execute tools you define
- Provide a clean web UI for interaction

**By design**, Toolbox CANNOT:
- Run arbitrary commands on your system
- Access files outside configured directories
- Make network requests without explicit tools

## Architecture

```
toolbox/
├── cmd/server/          # Main entry point
├── internal/
│   ├── agent/           # LLM orchestration
│   ├── tools/           # Tool registry
│   └── web/             # HTTP handlers
├── tools/               # User-defined tools
│   ├── database/        # PostgreSQL query tool
│   ├── files/           # File server (TODO)
│   └── docker/          # Docker management (TODO)
├── web/                 # Frontend (HTML/CSS/JS)
├── docker-compose.yml
└── CONTEXT.md          # Development documentation
```

## Adding Your Own Tools

Tools implement a simple interface:

```go
type Tool interface {
    Name() string
    Description() string
    Schema() map[string]interface{}  // JSON schema
    Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)
}
```

**Example**: See `tools/database/query.go` for a complete implementation.

## Configuration

Environment variables:

```bash
PORT=8080                          # Web server port
DATABASE_URL=postgres://...        # PostgreSQL connection
OPENAI_API_KEY=sk-...             # (Future) OpenAI API key
ANTHROPIC_API_KEY=sk-ant-...      # (Future) Anthropic API key
```

## Current Limitations

- **Mock LLM**: Currently uses a simple keyword-based mock provider
- **No real LLM**: Need to add OpenAI/Anthropic integration
- **Limited tools**: Only database tool implemented
- **No authentication**: Web UI is open (add auth before production)

## Roadmap

- [ ] Real LLM provider integration (OpenAI/Anthropic)
- [ ] File server tool (read/write controlled directory)
- [ ] Docker tool (inspect/manage containers)
- [ ] Tool permission system (require approval for dangerous ops)
- [ ] Conversation history persistence
- [ ] Streaming responses
- [ ] Authentication/authorization

## Why Not Fork OpenCode?

OpenCode is excellent, but:

1. **TypeScript/Node ecosystem**: 100k+ lines, heavy npm dependencies
2. **Broad capabilities**: Terminal-focused, many features we don't need
3. **MIT license**: Permissive (Toolbox uses AGPL-3.0 for copyleft)
4. **Learning curve**: Complex codebase vs. our simple Go implementation

Building from scratch gave us:
- Full control over tool execution model
- Minimal dependencies (Go stdlib > npm packages)
- Clear security boundaries
- Simpler maintenance

## License

AGPL-3.0 - See LICENSE file

This ensures that any improvements or derivative works remain open source.

## Development

See [CONTEXT.md](CONTEXT.md) for detailed development documentation, architecture decisions, and current project state.

## Questions?

- Check [CONTEXT.md](CONTEXT.md) for project vision and technical decisions
- Review tool implementations in `tools/` directory
- Examine agent orchestration in `internal/agent/`

---

**Built with ❤️ and a commitment to copyleft licensing**
