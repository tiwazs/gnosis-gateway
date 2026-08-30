package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Listen                string
	MainService           string
	WorkspaceService      string
	IotService            string
	TokenKey              string
	GatewayInternalSecret string
}

func Load() Config {
	_ = godotenv.Load()

	config := Config{
		Listen:           ":3000",
		MainService:      "http://localhost:4000",
		WorkspaceService: "http://localhost:4000",
		IotService:       "http://localhost:3000",
	}

	if v := os.Getenv("LISTEN"); v != "" {
		config.Listen = v
	}
	if v := os.Getenv("MAIN_SERVICE"); v != "" {
		config.MainService = v
	}
	if v := os.Getenv("WORKSPACE_SERVICE"); v != "" {
		config.WorkspaceService = v
	}
	if v := os.Getenv("IOT_SERVICE"); v != "" {
		config.IotService = v
	}
	config.TokenKey = os.Getenv("TOKEN_KEY")
	config.GatewayInternalSecret = os.Getenv("GATEWAY_INTERNAL_SECRET")

	return config
}
