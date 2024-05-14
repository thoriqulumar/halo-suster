package application

import (
	"helo-suster/config"
	"helo-suster/database"
	"helo-suster/pkg/log"
	"helo-suster/server"
	"helo-suster/version"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Start(cfg *config.Config) {
	// init logger
	logger, err := log.New(zapcore.DebugLevel, version.ServiceID, version.Version)
	if err != nil {
		panic(err)
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		logger.Error("error opening database", zap.Error(err))
		panic(err)
	}
	defer db.Close()

	s := server.NewServer(db, logger)
	s.RegisterRoute(cfg)

	logger.Fatal("failed run app", zap.Error(s.Start()))
}
