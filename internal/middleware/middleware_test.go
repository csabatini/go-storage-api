package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// orderMiddleware appends a marker before and after calling next.
func orderMiddleware(name string, order *[]string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			*order = append(*order, name+"-before")
			next.ServeHTTP(w, r)
			*order = append(*order, name+"-after")
		})
	}
}

func TestChain_Ordering(t *testing.T) {
	var order []string

	chained := Chain(
		orderMiddleware("A", &order),
		orderMiddleware("B", &order),
		orderMiddleware("C", &order),
	)

	handler := chained(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expected := "A-before,B-before,C-before,handler,C-after,B-after,A-after"
	got := strings.Join(order, ",")
	if got != expected {
		t.Errorf("chain order:\n  got  %s\n  want %s", got, expected)
	}
}

func TestChain_Empty(t *testing.T) {
	called := false
	handler := Chain()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if !called {
		t.Error("handler should be called with empty chain")
	}
}

func TestChain_SingleMiddleware(t *testing.T) {
	var order []string

	handler := Chain(
		orderMiddleware("A", &order),
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expected := "A-before,handler,A-after"
	got := strings.Join(order, ",")
	if got != expected {
		t.Errorf("single middleware order:\n  got  %s\n  want %s", got, expected)
	}
}
