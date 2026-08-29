package main

import (
	"log"
	"net/http"
	"github.com/tiwazs/gnosis-gateway/internal/config"
	"github.com/tiwazs/gnosis-gateway/internal/proxy"
)

func main() {
	cfg := config.Load()

	MainService, err := proxy.New(cfg.MainService)

	if err != nil {
		log.Fatalf("Main server: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("ok"))
	})

	mux.Handle("/", MainService)
	log.Printf("gateway listening on %s -> %s", cfg.Listen, cfg.MainService)

	if err := http.ListenAndServe(cfg.Listen, mux); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}