package auth

import (
	"net/http"
	"testing"
)

func TestSkip(t *testing.T) {
	cases := []struct {
		method string
		path   string
		skip   bool
	}{
		{"GET", "/healthz", true},
		{"POST", "/api/access/login", true},
		{"POST", "/api/access/register", true},
		{"POST", "/iot/devices/register/example-code", true},
		{"POST", "/devices/register/example-code", true},
		{"POST", "/iot/devices/register/", false},
		{"POST", "/devices/register/", false},
		{"GET", "/iot/devices/register/example-code", false},
		{"POST", "/internal/auth/token", true},
		{"GET", "/internal/auth/apikey", true},
		{"OPTIONS", "/api/profile", true},
		{"OPTIONS", "/workspace/workspaces", true},
		{"OPTIONS", "/iot/devices", true},
		{"GET", "/main/docs", true},
		{"GET", "/main/docs/swagger-ui.css", true},
		{"GET", "/workspace/docs", true},
		{"GET", "/workspace/docs/", true},
		{"GET", "/workspace/docs/index.html", true},
		{"GET", "/workspace/docs/doc.json", true},
		{"GET", "/iot/docs", true},
		{"GET", "/iot/redoc", true},
		{"GET", "/iot/openapi.json", true},
		{"GET", "/api/profile", false},
		{"GET", "/workspace/foo", false},
		{"GET", "/workspace/workspaces", false},
		{"GET", "/iot/bar", false},
	}
	for _, tc := range cases {
		req, _ := http.NewRequest(tc.method, tc.path, nil)
		if got := Skip(req); got != tc.skip {
			t.Errorf("%s %s: skip=%v want %v", tc.method, tc.path, got, tc.skip)
		}
	}
}
