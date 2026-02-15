package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestLogger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func TestLogging_CapturesExplicitStatus(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	handler := Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodPost, "/items", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	entry := parseLogEntry(t, &buf)
	assertLogField(t, entry, "method", "POST")
	assertLogField(t, entry, "path", "/items")
	assertLogFieldFloat(t, entry, "status", 201)
}

func TestLogging_ImplicitOK(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	handler := Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	entry := parseLogEntry(t, &buf)
	assertLogFieldFloat(t, entry, "status", 200)
}

func TestLogging_IncludesRequestID(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	inner := Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler := RequestID(inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	entry := parseLogEntry(t, &buf)
	rid, ok := entry["request_id"].(string)
	if !ok || rid == "" {
		t.Error("expected non-empty request_id in log entry")
	}
}

func TestLogging_IncludesDuration(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	handler := Logging(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	entry := parseLogEntry(t, &buf)
	dur, ok := entry["duration"].(string)
	if !ok || dur == "" {
		t.Error("expected non-empty duration in log entry")
	}
}

// helpers

func parseLogEntry(t *testing.T, buf *bytes.Buffer) map[string]interface{} {
	t.Helper()
	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse log JSON: %v\nraw: %s", err, buf.String())
	}
	return entry
}

func assertLogField(t *testing.T, entry map[string]interface{}, key, want string) {
	t.Helper()
	got, ok := entry[key].(string)
	if !ok {
		t.Errorf("log field %q missing or not a string", key)
		return
	}
	if got != want {
		t.Errorf("log field %q = %q, want %q", key, got, want)
	}
}

func assertLogFieldFloat(t *testing.T, entry map[string]interface{}, key string, want float64) {
	t.Helper()
	got, ok := entry[key].(float64)
	if !ok {
		t.Errorf("log field %q missing or not a number", key)
		return
	}
	if got != want {
		t.Errorf("log field %q = %v, want %v", key, got, want)
	}
}
