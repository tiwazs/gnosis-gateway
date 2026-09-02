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

type corsWriter struct {
	http.ResponseWriter
	req   *http.Request
	wrote bool
}

func (w *corsWriter) WriteHeader(code int) {
	if !w.wrote {
		applyCORS(w.ResponseWriter, w.req)
		w.wrote = true
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *corsWriter) Write(b []byte) (int, error) {
	if !w.wrote {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}

// Wrap answers OPTIONS and stamps exactly one Allow-Origin on GET/POST.
// Backends (FastAPI) must not also send CORS; the proxy strips those first.
func Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			applyCORS(w, r)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(&corsWriter{ResponseWriter: w, req: r}, r)
	})
}
