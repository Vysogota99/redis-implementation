package server

import (
	"fmt"
	"os"
)

// Config ...
type Config struct {
	serverPort string
}

// NewConfig - helper to init config
func NewConfig() (*Config, error) {
	serverPort, exists := os.LookupEnv("SERVER_PORT")
	if !exists {
		return nil, fmt.Errorf("No SERVER_PORT in .env")
	}

	return &Config{
		serverPort: serverPort,
	}, nil
}
