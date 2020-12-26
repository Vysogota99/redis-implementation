package server

import (
	"github.com/Vysogota99/redis-implementation/internal/server/store"
	"github.com/gorilla/sessions"

	"gopkg.in/boj/redistore.v1"
)

// Server ...
type Server struct {
	conf         *Config
	router       *router
	redis        store.RedisImpl
	sessionStore sessions.Store
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

	if err := s.initSessionStore(); err != nil {
		return err
	}

	s.router = newRouter(s.conf.serverPort, s.conf.sessionName, s.redis, s.sessionStore)
	s.router.setup().Run(s.conf.serverPort)

	return nil
}

func (s *Server) initSessionStore() error {
	store, err := redistore.NewRediStore(s.conf.sessionMaxNumberIDLEConnections, "tcp", s.conf.redisAddr, "", []byte(s.conf.sessionKey))
	if err != nil {
		return err
	}

	s.sessionStore = store
	return nil
}
