package server

import (
	"fmt"
	"os"
	"strconv"
)

// Config ...
type Config struct {
	serverPort                      string
	redisAddr                       string
	sessionKey                      string
	sessionMaxNumberIDLEConnections int
	sessionName                     string
}

// NewConfig - helper to init config
func NewConfig() (*Config, error) {
	serverPort, exists := os.LookupEnv("SERVER_PORT")
	if !exists {
		return nil, fmt.Errorf("No SERVER_PORT in .env")
	}

	redisAddr, exists := os.LookupEnv("REDIS_ADDR")
	if !exists {
		return nil, fmt.Errorf("No REDIS_ADDR in .env")
	}

	sessionKey, exists := os.LookupEnv("SESSION_KEY")
	if !exists {
		return nil, fmt.Errorf("No SESSION_KEY in .env")
	}

	maxIDLEconn, exists := os.LookupEnv("MAX_IDLE_SESSION_CONN")
	if !exists {
		return nil, fmt.Errorf("No MAX_IDLE_SESSION_CONN in .env")
	}
	imaxIDLEconn, err := strconv.Atoi(maxIDLEconn)
	if err != nil {
		return nil, err
	}

	return &Config{
		serverPort:                      serverPort,
		redisAddr:                       redisAddr,
		sessionMaxNumberIDLEConnections: imaxIDLEconn,
		sessionKey:                      sessionKey,
		sessionName:                     "auth",
	}, nil
}
