package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
)

type contextKey string

const requestIDKey contextKey = "request_id"

const headerXRequestID = "X-Request-ID"

// RequestID injects a UUID v4 request ID into the context and response header.
// If the incoming request already has an X-Request-ID header, it is preserved.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(headerXRequestID)
		if id == "" {
			id = newUUIDv4()
		}

		w.Header().Set(headerXRequestID, id)
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestIDFromContext extracts the request ID stored by the RequestID middleware.
func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

// newUUIDv4 generates a RFC 4122 version 4 UUID using crypto/rand.
func newUUIDv4() string {
	var uuid [16]byte
	_, err := rand.Read(uuid[:])
	if err != nil {
		// crypto/rand should never fail on supported platforms;
		// fall back to a zero UUID rather than panicking.
		return "00000000-0000-4000-8000-000000000000"
	}

	// Set version 4 (bits 4-7 of byte 6).
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	// Set variant 10xx (bits 6-7 of byte 8).
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])
}
