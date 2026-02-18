// Package web provides HTTP handlers and web UI
package web

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jredh-dev/nascent-nexus/toolbox/internal/agent"
)

// Server handles HTTP requests
type Server struct {
	agent     *agent.Agent
	mux       *http.ServeMux
	templates *template.Template
}

// NewServer creates a new web server
func NewServer(ag *agent.Agent) (*Server, error) {
	s := &Server{
		agent: ag,
		mux:   http.NewServeMux(),
	}

	// Parse templates
	tmplPath := filepath.Join("web", "templates", "*.html")
	tmpl, err := template.ParseGlob(tmplPath)
	if err != nil {
		// Try absolute path
		cwd, _ := os.Getwd()
		tmplPath = filepath.Join(cwd, "web", "templates", "*.html")
		tmpl, err = template.ParseGlob(tmplPath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse templates: %w", err)
		}
	}
	s.templates = tmpl

	// Register routes
	s.mux.HandleFunc("/", s.handleIndex)
	s.mux.HandleFunc("/api/message", s.handleMessage)
	s.mux.HandleFunc("/api/history", s.handleHistory)
	s.mux.HandleFunc("/api/reset", s.handleReset)
	s.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	return s, nil
}

// ServeHTTP implements http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// handleIndex serves the main page
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	err := s.templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// MessageRequest represents an incoming message
type MessageRequest struct {
	Message string `json:"message"`
}

// MessageResponse represents the agent's response
type MessageResponse struct {
	Response string `json:"response"`
	Error    string `json:"error,omitempty"`
}

// handleMessage processes user messages
func (s *Server) handleMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, MessageResponse{Error: "Invalid request"})
		return
	}

	if req.Message == "" {
		sendJSON(w, http.StatusBadRequest, MessageResponse{Error: "Message cannot be empty"})
		return
	}

	// Process with agent
	ctx := context.Background()
	response, err := s.agent.Process(ctx, req.Message)
	if err != nil {
		log.Printf("Agent error: %v", err)
		sendJSON(w, http.StatusInternalServerError, MessageResponse{Error: err.Error()})
		return
	}

	sendJSON(w, http.StatusOK, MessageResponse{Response: response})
}

// handleHistory returns conversation history
func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	history := s.agent.History()
	sendJSON(w, http.StatusOK, history)
}

// handleReset clears conversation history
func (s *Server) handleReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.agent.Reset()
	sendJSON(w, http.StatusOK, map[string]string{"status": "reset"})
}

// sendJSON sends a JSON response
func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
