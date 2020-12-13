package main

import (
	"fmt"

	"github.com/Vysogota99/redis-implementation/internal/server/server"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("Error could not find .env file: %w", err))
	}
}

func main() {
	conf, err := server.NewConfig()
	if err != nil {
		panic(err)
	}

	server := server.NewServer(conf)
	if err := server.Start(); err != nil {
		panic(err)
	}
}
