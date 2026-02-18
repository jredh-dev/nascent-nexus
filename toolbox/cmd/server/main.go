// Toolbox - A minimal, tool-focused AI assistant
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jredh-dev/nascent-nexus/toolbox/internal/agent"
	"github.com/jredh-dev/nascent-nexus/toolbox/internal/tools"
	"github.com/jredh-dev/nascent-nexus/toolbox/internal/web"
	"github.com/jredh-dev/nascent-nexus/toolbox/tools/database"
)

func main() {
	log.Println("🧰 Starting Toolbox...")

	// Get configuration from environment
	port := getEnv("PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "postgres://toolbox:toolbox_dev_password@localhost:5432/toolbox?sslmode=disable")

	// Create tool registry
	registry := tools.NewRegistry()

	// Register database tool
	dbTool, err := database.NewQueryTool(dbURL)
	if err != nil {
		log.Fatalf("Failed to create database tool: %v", err)
	}
	defer dbTool.Close()

	if err := registry.Register(dbTool); err != nil {
		log.Fatalf("Failed to register database tool: %v", err)
	}
	log.Println("✓ Registered database tool")

	// Create LLM provider (mock for now)
	// TODO: Replace with real provider (OpenAI/Anthropic) when API keys are available
	provider := agent.NewMockProvider()
	log.Println("✓ Using mock LLM provider (replace with real provider later)")

	// Create agent
	ag := agent.New(provider, registry)
	log.Println("✓ Created agent")

	// Create web server
	server, err := web.NewServer(ag)
	if err != nil {
		log.Fatalf("Failed to create web server: %v", err)
	}
	log.Println("✓ Created web server")

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("🚀 Toolbox running at http://localhost:%s", port)
	log.Printf("📚 Documentation: See CONTEXT.md for project overview")

	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
