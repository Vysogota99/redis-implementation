package server

// Server ...
type Server struct {
	conf   *Config
	router *router
}

// NewServer - helper to init server
func NewServer(conf *Config) *Server {
	return &Server{
		conf: conf,
	}
}

// Start - start the server
func (s *Server) Start() error {
	s.router = newRouter(s.conf.serverPort)
	s.router.setup().Run(s.conf.serverPort)

	return nil
}
