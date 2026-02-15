package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPathGuard_BlocksTraversal(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{"double dot", "../etc/passwd"},
		{"leading slash with traversal", "/../../../etc/passwd"},
		{"mid-path traversal", "files/../../../secret"},
		{"encoded double dot", "%2e%2e/etc/passwd"},
		{"encoded slash and dots", "%2e%2e%2fetc%2fpasswd"},
		{"mixed encoding", "..%2fetc/passwd"},
	}

	handler := PathGuard(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not have been called")
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build the URL with the raw path value to preserve encoding.
			req := httptest.NewRequest(http.MethodGet, "/files?path="+tt.path, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d", rr.Code)
			}

			var body errorResponse
			if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if body.Error == "" {
				t.Error("expected non-empty error message")
			}
		})
	}
}

func TestPathGuard_BlocksNullBytes(t *testing.T) {
	handler := PathGuard(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not have been called")
	}))

	t.Run("encoded null byte", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/files?path=file%00.txt", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})
}

func TestPathGuard_AllowsValidPaths(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"simple filename", "readme.txt", "readme.txt"},
		{"nested path", "docs/guide/intro.md", "docs/guide/intro.md"},
		{"trailing slash cleaned", "docs/guide/", "docs/guide"},
		{"double slash cleaned", "docs//guide", "docs/guide"},
		{"dot current dir", "./readme.txt", "readme.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath string
			handler := PathGuard(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Query().Get("path")
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "/files?path="+tt.path, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("expected 200, got %d", rr.Code)
			}
			if gotPath != tt.expected {
				t.Errorf("expected cleaned path %q, got %q", tt.expected, gotPath)
			}
		})
	}
}

func TestPathGuard_NoPathParam(t *testing.T) {
	called := false
	handler := PathGuard(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if !called {
		t.Error("handler should have been called when no path param present")
	}
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}
