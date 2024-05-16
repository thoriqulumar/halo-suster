package server

import (
	"helo-suster/config"
	"helo-suster/controller"
	"helo-suster/service"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func (s *Server) RegisterRoute(cfg *config.Config) {
	mainRoute := s.app.Group("/v1")

	registerImageRoute(mainRoute, s.db, cfg)
}

func registerImageRoute(e *echo.Group, db *sqlx.DB, cfg *config.Config) {
	ctr := controller.NewImageController(service.NewImageService())

	e.POST("/image", ctr.PostImage)
}
