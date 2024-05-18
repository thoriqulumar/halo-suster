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

	registerImageRoute(mainRoute, cfg, s.logger)
	registerMedicalRoute(mainRoute, s.db, cfg, s.validator, s.logger)
	registerStaffRoute(mainRoute, s.db, cfg, s.validator)
}

func registerImageRoute(e *echo.Group, cfg *config.Config, logger *zap.Logger) {
	ctr := controller.NewImageController(service.NewImageService(cfg, logger))

	e.POST("/image", ctr.PostImage)
}

func registerMedicalRoute(e *echo.Group, db *sqlx.DB, cfg *config.Config, validate *validator.Validate, logger *zap.Logger) {
	ctr := controller.NewMedicalController(service.NewMedicalService(repo.NewMedicalRepo(db), logger), validate)

	e.POST("/medical/patient", ctr.PostPatient)
	e.POST("/medical/record", ctr.PostMedicalReport)
}

func registerStaffRoute(e *echo.Group, db *sqlx.DB, cfg *config.Config, validate *validator.Validate) {
	ctr := controller.NewStaffController(service.NewStaffService(cfg, repo.NewStaffRepo(db)), validate)

	e.POST("/user/it/register", ctr.Register)
}
