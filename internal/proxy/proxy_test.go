package proxy

import (
	"net/http"
	"testing"
)

func TestStripCORS(t *testing.T) {
	h := http.Header{}
	h.Add("Access-Control-Allow-Origin", "http://localhost:3000")
	h.Add("Access-Control-Allow-Origin", "*")
	h.Set("Content-Type", "application/json")
	stripCORS(h)
	if got := h.Values("Access-Control-Allow-Origin"); len(got) != 0 {
		t.Fatalf("expected CORS stripped, got %v", got)
	}
	if h.Get("Content-Type") != "application/json" {
		t.Fatal("stripped content-type")
	}
}
