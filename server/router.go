package server

import (
	"halo-suster/config"
	"halo-suster/controller"
	"halo-suster/middleware"
	"halo-suster/repo"
	"halo-suster/service"

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
	auth := middleware.Authentication(cfg.JWTSecret)
	e.POST("/image", auth(ctr.PostImage))
}

func registerMedicalRoute(e *echo.Group, db *sqlx.DB, cfg *config.Config, validate *validator.Validate, logger *zap.Logger) {
	ctr := controller.NewMedicalController(service.NewMedicalService(repo.NewMedicalRepo(db), logger), validate)

	auth := middleware.Authentication(cfg.JWTSecret)
	e.POST("/medical/patient", auth(ctr.PostPatient))
	e.POST("/medical/record", auth(ctr.PostMedicalReport))
	e.GET("/medical/patient", auth(ctr.GetPatient))
	e.GET("/medical/record", auth(ctr.GetMedicalRecord))
}

func registerStaffRoute(e *echo.Group, db *sqlx.DB, cfg *config.Config, validate *validator.Validate) {
	ctr := controller.NewStaffController(service.NewStaffService(cfg, repo.NewStaffRepo(db)), validate)

	auth := middleware.AuthenticationIT(cfg.JWTSecret)
	e.POST("/user/it/register", ctr.RegisterIT)
	e.POST("/user/it/login", ctr.LoginIT)

	e.POST("/user/nurse/register", auth(ctr.RegisterNurse))
	e.POST("/user/nurse/login", ctr.LoginNurse)

	e.GET("/user", auth(ctr.GetUser))

	e.PUT("/user/nurse/:id", auth(ctr.UpdateNurse))
	e.DELETE("/user/nurse/:id", auth(ctr.DeleteNurse))
	e.POST("/user/nurse/:id/access", auth(ctr.GrantAccessNurse))
}
