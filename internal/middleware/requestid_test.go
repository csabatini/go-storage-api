package middleware

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

var uuidV4Re = regexp.MustCompile(
	`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`,
)

func TestRequestID_GeneratesValidUUID(t *testing.T) {
	var captured string
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = RequestIDFromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if !uuidV4Re.MatchString(captured) {
		t.Errorf("context request ID %q does not match UUID v4 format", captured)
	}

	header := rr.Header().Get(headerXRequestID)
	if header != captured {
		t.Errorf("response header %q != context value %q", header, captured)
	}
}

func TestRequestID_PreservesIncomingHeader(t *testing.T) {
	incoming := "my-custom-id-12345"

	var captured string
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = RequestIDFromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(headerXRequestID, incoming)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if captured != incoming {
		t.Errorf("expected context ID %q, got %q", incoming, captured)
	}
	if rr.Header().Get(headerXRequestID) != incoming {
		t.Errorf("expected response header %q, got %q", incoming, rr.Header().Get(headerXRequestID))
	}
}

func TestRequestID_ContextRoundTrip(t *testing.T) {
	var captured string
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = RequestIDFromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if captured == "" {
		t.Fatal("RequestIDFromContext returned empty string")
	}
	if rr.Header().Get(headerXRequestID) != captured {
		t.Errorf("header %q != context %q", rr.Header().Get(headerXRequestID), captured)
	}
}

func TestRequestID_Uniqueness(t *testing.T) {
	ids := make(map[string]bool)

	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	for i := 0; i < 100; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		id := rr.Header().Get(headerXRequestID)
		if ids[id] {
			t.Fatalf("duplicate request ID on iteration %d: %s", i, id)
		}
		ids[id] = true
	}
}

func TestNewUUIDv4_Format(t *testing.T) {
	for i := 0; i < 50; i++ {
		id := newUUIDv4()

		if len(id) != 36 {
			t.Errorf("iteration %d: expected length 36, got %d for %q", i, len(id), id)
		}
		if !uuidV4Re.MatchString(id) {
			t.Errorf("iteration %d: %q does not match UUID v4 pattern %s", i, id, uuidV4Re.String())
		}
	}
}

func TestRequestIDFromContext_EmptyContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	id := RequestIDFromContext(req.Context())
	if id != "" {
		t.Errorf("expected empty string from bare context, got %q", id)
	}
}
