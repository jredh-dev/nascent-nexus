# CONTEXT.md

## Project Overview

- **Purpose**: Personal AI assistant system that acts as a "mind" - remembering, reminding, and coordinating across all services. Accessible via SMS for maximum ubiquity.
- **Vision**: 
  - Personal secretary that connects to any/all services and coordinates between them
  - IFTTT/BetterCloud-like integration hub
  - Reddit-style question system (ask questions of your own crowd)
  - Voice chat with general lobby
  - SMS textbot with LLM-based security questions for identity verification
  - Recording system for people who are scared, want to be remembered, or need to talk
  - Human-in-the-loop for sensitive conversations (AI asks humans for advice)
  - "The Big Question" - daily question asked to community, answers optionally contribute to LLM training (with credit/currency rewards)
  - Eventually: consultancy offering around "socioprogrammatics" (tacit prompts utilizing underutilized LLM features)
- **Status**: Phase 1 - Initial setup and SMS "hello world"
- **Monetization**: Open source (eventually), but prototyping privately first. Plans for credits/currency system for community contributions.
- **Privacy**: Minimal identity ties for early versions, but owner name on all official accounts (Twilio, etc.). No nefarious use protection needed.

## Architecture Decisions

- **Tech Stack**:
  - **Go**: Minimal footprint, efficient energy usage, compiles to native binary, low memory, fast startup
  - **PostgreSQL**: High-relation, high-complexity database with flexible multi-format support (JSONB)
  - **Twilio**: SMS provider (API-driven, minimal manual steps after initial setup)
  - **ngrok**: Local development webhook receiver (secure tunnel to localhost)
  - **Docker Compose**: For PostgreSQL only (keeps laptop clean)
- **Design Patterns**:
  - **Minimal footprint**: Energy-efficient at all times (reducing CPU load, memory usage)
  - **Monorepo**: All projects/apps isolated by deployment
  - **Local-first**: Runs on local laptop/PC initially, can scale to dedicated server later
  - **API-driven**: Twilio fully managed via API (minimal manual steps)
- **Trade-offs**:
  - **Go over Node.js/Python**: Native compilation = smallest memory footprint, fastest startup
  - **PostgreSQL over MongoDB/Redis**: More powerful relations + JSONB flexibility for multi-format data
  - **Local over cloud**: Privacy, control, prototyping without recurring costs
  - **ngrok over VPS**: Faster iteration during development, easy HTTPS for webhooks

## Current State

### What's Working
- ✅ Repository created (`jredh-dev/nascent-nexus`)
- ✅ Workstream initialized
- ✅ CONTEXT.md documented

### What's Not
- ❌ SMS webhook endpoint (not yet built)
- ❌ Twilio account setup (manual step required)
- ❌ PostgreSQL setup (Docker Compose not yet configured)
- ❌ ngrok configuration (not yet set up)

### Next Steps
1. Create initial project structure (Go module, directories)
2. Set up CHANGELOG.md (Keep a Changelog format)
3. Add LICENSE (AGPL-3.0)
4. Create basic Go HTTP server with `/sms` webhook endpoint
5. Add ngrok config file
6. Document Twilio setup instructions (manual steps)
7. Test end-to-end: text "hello" → receive "world"

## Task Tracking

### Available
- [ ] Initialize Go module and project structure
- [ ] Create CHANGELOG.md
- [ ] Add AGPL-3.0 LICENSE
- [ ] Build basic HTTP server with health check
- [ ] Add `/sms` webhook endpoint (returns "world" for any message)
- [ ] Create ngrok configuration
- [ ] Write Twilio setup instructions
- [ ] Add Docker Compose for PostgreSQL (for future use)
- [ ] Update README.md with project description

### In Progress
- [ ] Setting up project foundation

### Completed
- [x] Repository created on GitHub
- [x] Workstream initialized
- [x] CONTEXT.md created with full project vision

## Development Notes

### Setup Instructions (Once Built)
```bash
# Clone and navigate
cd work/workstream_6ec1de07-3748-4b22-8c8c-58c186fcf15e/nascent-nexus

# Build
go build -o bin/nascent-nexus cmd/server/main.go

# Run
./bin/nascent-nexus

# In another terminal: Start ngrok
ngrok http 8080

# Configure Twilio webhook URL to ngrok URL + /sms
```

### Testing Strategy
- Manual testing via SMS (text the Twilio number, verify "world" response)
- Eventually: Unit tests for webhook handler
- Eventually: Integration tests for Twilio API

### Deployment
- **Phase 1**: Local laptop + ngrok (prototyping)
- **Phase 2**: Dedicated PC server (for friends/family testing)
- **Phase 3**: Cloud deployment (when ready for broader use)

## Session Log

- **[2025-02-17]**: Project inception. Created repository, initialized workstream, documented full vision in CONTEXT.md. Ready to build SMS "hello world" endpoint.
