package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/jredh-dev/nexus/services/cal/internal/database"
)

func testHandler(t *testing.T) *Handler {
	t.Helper()
	path := t.TempDir() + "/test.db"
	db, err := database.Open(path)
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
		os.Remove(path)
	})
	return New(db)
}

func testRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/cal/{token}.ics", h.Subscribe)
	r.Route("/api", func(r chi.Router) {
		r.Post("/feeds", h.CreateFeed)
		r.Get("/feeds", h.ListFeeds)
		r.Delete("/feeds/{id}", h.DeleteFeed)
		r.Get("/feeds/{id}/events", h.ListEvents)
		r.Post("/events", h.CreateEvent)
		r.Delete("/events/{id}", h.DeleteEvent)
	})
	return r
}

func TestCreateAndListFeeds(t *testing.T) {
	h := testHandler(t)
	r := testRouter(h)

	// Create feed
	body := `{"name":"Work Calendar"}`
	req := httptest.NewRequest(http.MethodPost, "/api/feeds", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("create feed: expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var created createFeedResp
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create response: %v", err)
	}
	if created.Name != "Work Calendar" {
		t.Errorf("expected name 'Work Calendar', got %q", created.Name)
	}
	if created.Token == "" {
		t.Error("expected non-empty token")
	}
	if created.URL == "" {
		t.Error("expected non-empty URL")
	}

	// List feeds
	req = httptest.NewRequest(http.MethodGet, "/api/feeds", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("list feeds: expected 200, got %d", w.Code)
	}

	var feeds []database.Feed
	if err := json.Unmarshal(w.Body.Bytes(), &feeds); err != nil {
		t.Fatalf("unmarshal feeds: %v", err)
	}
	if len(feeds) != 1 {
		t.Fatalf("expected 1 feed, got %d", len(feeds))
	}
}

func TestCreateEventAndSubscribe(t *testing.T) {
	h := testHandler(t)
	r := testRouter(h)

	// Create a feed
	body := `{"name":"Test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/feeds", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var feed createFeedResp
	if err := json.Unmarshal(w.Body.Bytes(), &feed); err != nil {
		t.Fatalf("unmarshal feed: %v", err)
	}

	// Create an event
	eventBody, _ := json.Marshal(map[string]interface{}{
		"feed_id":    feed.ID,
		"summary":    "Weekend Hackathon",
		"start":      "2026-02-21T10:00:00Z",
		"end":        "2026-02-21T18:00:00Z",
		"categories": "fun,code",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/events", bytes.NewReader(eventBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("create event: expected 201, got %d: %s", w.Code, w.Body.String())
	}

	// Subscribe to the feed
	req = httptest.NewRequest(http.MethodGet, "/cal/"+feed.Token+".ics", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("subscribe: expected 200, got %d", w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "text/calendar") {
		t.Errorf("expected Content-Type text/calendar, got %q", ct)
	}

	ics := w.Body.String()
	required := []string{
		"BEGIN:VCALENDAR",
		"BEGIN:VEVENT",
		"SUMMARY:Weekend Hackathon",
		"DTSTART:20260221T100000Z",
		"DTEND:20260221T180000Z",
		"CATEGORIES:fun,code",
		"END:VEVENT",
		"END:VCALENDAR",
	}
	for _, s := range required {
		if !strings.Contains(ics, s) {
			t.Errorf("iCal output missing %q", s)
		}
	}
}

func TestSubscribe_InvalidToken(t *testing.T) {
	h := testHandler(t)
	r := testRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/cal/nonexistent.ics", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 for invalid token, got %d", w.Code)
	}
}

func TestCreateFeed_ValidationErrors(t *testing.T) {
	h := testHandler(t)
	r := testRouter(h)

	// Empty name
	req := httptest.NewRequest(http.MethodPost, "/api/feeds", strings.NewReader(`{"name":""}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty name, got %d", w.Code)
	}

	// Invalid JSON
	req = httptest.NewRequest(http.MethodPost, "/api/feeds", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", w.Code)
	}
}

func TestCreateEvent_ValidationErrors(t *testing.T) {
	h := testHandler(t)
	r := testRouter(h)

	// Missing required fields
	req := httptest.NewRequest(http.MethodPost, "/api/events", strings.NewReader(`{"summary":"test"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing fields, got %d", w.Code)
	}

	// Bad date format
	body := `{"feed_id":"x","summary":"test","start":"not-a-date"}`
	req = httptest.NewRequest(http.MethodPost, "/api/events", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for bad date, got %d", w.Code)
	}
}

func TestDeleteFeedAndEvents(t *testing.T) {
	h := testHandler(t)
	r := testRouter(h)

	// Create feed
	req := httptest.NewRequest(http.MethodPost, "/api/feeds", strings.NewReader(`{"name":"Temp"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var feed createFeedResp
	if err := json.Unmarshal(w.Body.Bytes(), &feed); err != nil {
		t.Fatalf("unmarshal feed: %v", err)
	}

	// Delete it
	req = httptest.NewRequest(http.MethodDelete, "/api/feeds/"+feed.ID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}

	// Verify gone
	req = httptest.NewRequest(http.MethodGet, "/api/feeds", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var feeds []database.Feed
	if err := json.Unmarshal(w.Body.Bytes(), &feeds); err != nil {
		t.Fatalf("unmarshal feeds: %v", err)
	}
	if len(feeds) != 0 {
		t.Errorf("expected 0 feeds after delete, got %d", len(feeds))
	}
}
