package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWrapOptionsDoesNotReachNext(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		called = true
	})
	req := httptest.NewRequest(http.MethodOptions, "/workspace/workspaces", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Headers", "authorization")
	rec := httptest.NewRecorder()

	Wrap(next).ServeHTTP(rec, req)

	if called {
		t.Fatal("OPTIONS should not be proxied")
	}
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status=%d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Fatalf("allow-origin=%q", rec.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestWrapGETDoesNotAddCORS(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("[]"))
	})
	req := httptest.NewRequest(http.MethodGet, "/iot/devices", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()

	Wrap(next).ServeHTTP(rec, req)

	if rec.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("GET should not set CORS (ingress does); got %q", rec.Header().Get("Access-Control-Allow-Origin"))
	}
}
