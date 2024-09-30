package server

import (
	"fmt"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/RedrikShuhartRed/EfMobSongLib/config"
	"github.com/RedrikShuhartRed/EfMobSongLib/db"
	"github.com/RedrikShuhartRed/EfMobSongLib/external"
	"github.com/RedrikShuhartRed/EfMobSongLib/logger"
	"github.com/RedrikShuhartRed/EfMobSongLib/songLib-api/song/handler"
	service "github.com/RedrikShuhartRed/EfMobSongLib/songLib-api/song/service"
	"github.com/RedrikShuhartRed/EfMobSongLib/songLib-api/song/storer"
)

func RunServer() {
	logger := logger.InitLogger()
	defer logger.Sync()
	logger.Info("Start zap logger")
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal("Failed to load environment variable", zap.Error(err))
	}
	cfg := config.NewConfig()
	logger.Info(fmt.Sprintf("Your environment variable %v", cfg))

	db, err := db.NewDatabase(logger, cfg)
	if err != nil {
		logger.Fatal("Error creating database", zap.Error(err))
		return
	}
	defer db.Close()
	logger.Info("successfully connected to the database")

	logger.Info("Running migrations...")
	err = db.RunMigrate()
	if err != nil {
		logger.Fatal("error running migrations", zap.Error(err))
		return
	}
	logger.Info("Successfully migrations")

	st := storer.NewPgStorer(db.GetDB(), logger)
	logger.Info("Successfully created song storer")
	srv := service.NewSongService(st)
	logger.Info("Successfully created song service")
	hdl := handler.NewHandler(srv, logger)
	logger.Info("Successfully created song handler")
	r := handler.RegisterRoutes(hdl)
	logger.Info("Successfully registered routes")
	go external.StartServer(cfg)
	logger.Info(fmt.Sprintf("Started external server on port %s", cfg.ExternalPort))
	logger.Info(fmt.Sprintf("Started music library server on port %s", cfg.LibraryPort))
	err = r.Run(":" + cfg.LibraryPort)
	if err != nil {
		logger.Fatal("Error starting server", zap.Error(err))
	}

}
