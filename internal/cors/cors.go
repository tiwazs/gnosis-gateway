package cors

import (
	"net/http"
)

const allowHeaders = "Authorization, Content-Type, x-access-token, ngrok-skip-browser-warning, X-Requested-With"
const allowMethods = "GET, PUT, POST, DELETE, PATCH, OPTIONS"

func applyCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", allowMethods)
	requested := r.Header.Get("Access-Control-Request-Headers")
	if requested != "" {
		w.Header().Set("Access-Control-Allow-Headers", requested)
	} else {
		w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
	}
	w.Header().Add("Vary", "Origin")
}

// Wrap answers OPTIONS at the gateway. Actual GET/POST CORS is left to
// Ingress so we never emit a second Access-Control-Allow-Origin.
func Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			applyCORS(w, r)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
