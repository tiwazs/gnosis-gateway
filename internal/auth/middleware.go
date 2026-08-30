package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tiwazs/gnosis-gateway/internal/config"
)

var errInvalidAPIKey = errors.New("invalid api key")

type Claims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func Wrap(next http.Handler, cfg config.Config) http.Handler {
	lookup := newAPIKeyLookup(cfg.MainService, cfg.GatewayInternalSecret)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if Skip(r) {
			next.ServeHTTP(w, r)
			return
		}

		r.Header.Del("X-User-Id")
		r.Header.Del("X-User-Email")
		r.Header.Del("X-Auth-Method")

		token := extractToken(r)
		if token == "" {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "A token is required for authentication"})
			return
		}

		if looksLikeJWT(token) {
			if cfg.TokenKey == "" {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "TOKEN_KEY is not configured"})
				return
			}
			claims, err := parseJWT(token, cfg.TokenKey)
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
				return
			}
			r.Header.Set("X-User-Id", claims.UserID)
			r.Header.Set("X-User-Email", claims.Email)
			r.Header.Set("X-Auth-Method", "jwt")
			next.ServeHTTP(w, r)
			return
		}

		if cfg.GatewayInternalSecret == "" {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "GATEWAY_INTERNAL_SECRET is not configured"})
			return
		}
		id, err := lookup.lookup(token)
		if err != nil {
			if errors.Is(err, errInvalidAPIKey) {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid API key"})
				return
			}
			writeJSON(w, http.StatusBadGateway, map[string]string{"error": "API key lookup failed"})
			return
		}
		r.Header.Set("X-User-Id", id.ID)
		r.Header.Set("X-User-Email", id.Email)
		r.Header.Set("X-Auth-Method", "apikey")
		next.ServeHTTP(w, r)
	})
}

func parseJWT(tokenString, secret string) (*Claims, error) {
	claims := &Claims{}
	parsed, err := jwt.ParseWithClaims(tokenString, claims, func(parsed *jwt.Token) (interface{}, error) {
		if _, isHMAC := parsed.Method.(*jwt.SigningMethodHMAC); !isHMAC {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil || parsed == nil || !parsed.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	if claims.UserID == "" {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}

func writeJSON(w http.ResponseWriter, status int, body map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
