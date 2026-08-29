package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct{
	Listen		  string
	MainService  string 
}

func Load() Config {
	_ = godotenv.Load()
	
	config := Config{
		Listen: ":3000",
		MainService: "http://localhost:4000",
	}

	if env_read := os.Getenv("LISTEN") ; env_read != "" {
		config.Listen = env_read;
	}

	if env_read := os.Getenv("MAIN_SERVICE") ; env_read != "" {
		config.MainService = env_read;
	}

	return config
}