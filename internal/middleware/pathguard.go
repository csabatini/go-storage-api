package middleware

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// errorResponse mirrors the api.ErrorResponse JSON shape but is defined
// locally so pathguard has no dependency on internal/api.
type errorResponse struct {
	Error string `json:"error"`
}

// PathGuard rejects requests whose "path" query parameter contains directory
// traversal sequences (..) or null bytes. Valid paths are normalized with
// path.Clean before the request continues.
func PathGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := r.URL.Query().Get("path")
		if raw == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Decode to catch double-encoded traversal (%252e%252e).
		decoded, err := url.QueryUnescape(raw)
		if err != nil {
			writeErrorJSON(w, http.StatusBadRequest, "invalid path encoding")
			return
		}

		if containsTraversal(decoded) || containsNullByte(decoded) {
			writeErrorJSON(w, http.StatusBadRequest, "invalid path")
			return
		}

		// Normalize and replace the query parameter.
		cleaned := path.Clean(decoded)
		q := r.URL.Query()
		q.Set("path", cleaned)
		r.URL.RawQuery = q.Encode()

		next.ServeHTTP(w, r)
	})
}

func containsTraversal(s string) bool {
	return strings.Contains(s, "..")
}

func containsNullByte(s string) bool {
	return strings.ContainsRune(s, '\x00')
}

func writeErrorJSON(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{Error: msg})
}
