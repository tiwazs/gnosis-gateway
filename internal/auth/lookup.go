package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type identity struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type cacheEntry struct {
	id      identity
	expires time.Time
}

type apiKeyLookup struct {
	mainService string
	secret      string
	client      *http.Client
	mu          sync.Mutex
	cache       map[string]cacheEntry
}

func newAPIKeyLookup(mainService, secret string) *apiKeyLookup {
	return &apiKeyLookup{
		mainService: trimSlash(mainService),
		secret:      secret,
		client:      &http.Client{Timeout: 5 * time.Second},
		cache:       make(map[string]cacheEntry),
	}
}

func trimSlash(s string) string {
	if len(s) > 0 && s[len(s)-1] == '/' {
		return s[:len(s)-1]
	}
	return s
}

func (l *apiKeyLookup) lookup(key string) (identity, error) {
	l.mu.Lock()
	if entry, ok := l.cache[key]; ok && time.Now().Before(entry.expires) {
		id := entry.id
		l.mu.Unlock()
		return id, nil
	}
	l.mu.Unlock()

	endpoint := fmt.Sprintf("%s/internal/auth/apikey?key=%s", l.mainService, url.QueryEscape(key))
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return identity{}, err
	}
	req.Header.Set("X-Gateway-Secret", l.secret)

	resp, err := l.client.Do(req)
	if err != nil {
		return identity{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusUnauthorized {
		return identity{}, errInvalidAPIKey
	}
	if resp.StatusCode != http.StatusOK {
		return identity{}, fmt.Errorf("apikey lookup status %d", resp.StatusCode)
	}

	var id identity
	if err := json.NewDecoder(resp.Body).Decode(&id); err != nil {
		return identity{}, err
	}
	if id.ID == "" {
		return identity{}, errInvalidAPIKey
	}

	l.mu.Lock()
	l.cache[key] = cacheEntry{id: id, expires: time.Now().Add(30 * time.Second)}
	l.mu.Unlock()
	return id, nil
}
