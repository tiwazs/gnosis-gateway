package auth

import (
	"net/http"
	"strings"
)

func Skip(r *http.Request) bool {
	if r.Method == http.MethodOptions {
		return true
	}
	path := r.URL.Path
	if path == "/healthz" && r.Method == http.MethodGet {
		return true
	}
	if r.Method == http.MethodPost && (path == "/api/access/login" || path == "/api/access/register") {
		return true
	}
	if r.Method == http.MethodPost && isDeviceRegister(path) {
		return true
	}
	if strings.HasPrefix(path, "/internal/auth") {
		return true
	}
	if isPublicDocs(path) {
		return true
	}
	return false
}

// Swagger/OpenAPI UIs and their static assets. API routes stay authenticated.
func isPublicDocs(path string) bool {
	for _, prefix := range []string{
		"/main/docs",
		"/workspace/docs",
		"/iot/docs",
		"/iot/redoc",
		"/iot/openapi.json",
	} {
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return true
		}
	}
	return false
}

func isDeviceRegister(path string) bool {
	for _, prefix := range []string{"/iot/devices/register/", "/devices/register/"} {
		if strings.HasPrefix(path, prefix) && path != prefix {
			return true
		}
	}
	return false
}
