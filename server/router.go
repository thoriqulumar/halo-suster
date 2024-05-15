package server

import (
	"helo-suster/config"
	"helo-suster/controller"
	"helo-suster/repo"
	"helo-suster/service"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (s *Server) RegisterRoute(cfg *config.Config) {
	mainRoute := s.app.Group("/v1")

	registerMedicalRoute(mainRoute, s.db, cfg, s.validator, s.logger)
}

func registerMedicalRoute(e *echo.Group, db *sqlx.DB, cfg *config.Config, validate *validator.Validate, logger *zap.Logger) {
	ctr := controller.NewMedicalController(service.NewMedicalService(repo.NewMedicalRepo(db), logger), validate)

	e.POST("/medical/patient", ctr.PostPatient)
}