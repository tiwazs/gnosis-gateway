package cors

import (
	"net/http"
)

const allowHeaders = "Authorization, Content-Type, x-access-token, ngrok-skip-browser-warning, X-Requested-With"
const allowMethods = "GET, PUT, POST, DELETE, PATCH, OPTIONS"

// Wrap answers CORS at the gateway (including OPTIONS) so Gin/FastAPI 404s never
// reach the browser as a failed preflight.
func Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
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

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
