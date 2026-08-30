package auth

import (
	"net/http"
	"strings"
)

func extractToken(r *http.Request) string {
	if authorization := r.Header.Get("Authorization"); authorization != "" {
		if strings.HasPrefix(strings.ToLower(authorization), "bearer ") {
			return strings.TrimSpace(authorization[7:])
		}
		return strings.TrimSpace(authorization)
	}
	if accessToken := r.Header.Get("x-access-token"); accessToken != "" {
		return strings.TrimSpace(accessToken)
	}
	if queryToken := r.URL.Query().Get("token"); queryToken != "" {
		return strings.TrimSpace(queryToken)
	}
	return ""
}

func looksLikeJWT(token string) bool {
	return strings.Count(token, ".") == 2
}
