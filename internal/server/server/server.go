package server

import (
	"github.com/Vysogota99/redis-implementation/internal/server/store"
)

// Server ...
type Server struct {
	conf   *Config
	router *router
	redis  *store.Redis
}

// NewServer - helper to init server
func NewServer(conf *Config) *Server {
	return &Server{
		conf: conf,
	}
}

// Start - start the server
func (s *Server) Start() error {
	redis, err := store.New(s.conf.redisAddr)
	if err != nil {
		return err
	}

	s.redis = redis

	s.router = newRouter(s.conf.serverPort, s.redis)
	s.router.setup().Run(s.conf.serverPort)

	return nil
}
