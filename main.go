package main

import (
	"log"
	"net/http"

	"github.com/tiwazs/gnosis-gateway/internal/auth"
	"github.com/tiwazs/gnosis-gateway/internal/config"
	appcors "github.com/tiwazs/gnosis-gateway/internal/cors"
	"github.com/tiwazs/gnosis-gateway/internal/proxy"
)

func main() {
	cfg := config.Load()

	MainService, err := proxy.New(cfg.MainService)
	if err != nil {
		log.Fatalf("Main server: %v", err)
	}

	WorkspaceService, err := proxy.New(cfg.WorkspaceService)
	if err != nil {
		log.Fatalf("Workspace server: %v", err)
	}

	IotService, err := proxy.New(cfg.IotService)
	if err != nil {
		log.Fatalf("Iot server: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("ok"))
	})

	// Endpoints

	mux.Handle("/workspace", WorkspaceService)
	mux.Handle("/workspace/", WorkspaceService)
	mux.Handle("/iot", IotService)
	mux.Handle("/iot/", IotService)

	mux.Handle("/", MainService)

	log.Printf("Workspace server: %s", cfg.WorkspaceService)
	log.Printf("Iot server: %s", cfg.IotService)
	log.Printf("Main server: %s", cfg.MainService)
	log.Printf("gateway listening on %s -> %s", cfg.Listen, cfg.MainService)

	if err := http.ListenAndServe(cfg.Listen, appcors.Wrap(auth.Wrap(mux, cfg))); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
