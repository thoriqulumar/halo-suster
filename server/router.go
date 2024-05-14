package server

import "helo-suster/config"

func (s *Server) RegisterRoute(cfg *config.Config) {
	s.app.Group("/v1")
}
