package main

import (
	"learnDB/internal/config"
	"learnDB/internal/dbManager"
	"learnDB/internal/testReviewer"
	"learnDB/internal/testReviewer/web"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	config := config.MustLoad()

	logFile, err := os.Create("./logfiles")
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	logger := slog.New(slog.NewJSONHandler(logFile, nil))

	repo := make(map[string]testReviewer.DBRepository)
	for name, connStr := range config.Databases {
		logger.Info("connecting to database", slog.String("db", name))
		manager, err := dbManager.New(name, connStr)
		if err != nil {
			logger.Warn("cannot create dbManager to db", slog.String("db", name), slog.Any("error", err))
		}

		repo[name] = manager
	}

	reviewer := testReviewer.New(repo, "olympics")

	controller := web.NewController(logger, reviewer)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /upload", controller.ExcelHandler)
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	s := http.Server{
		Addr:    config.Address,
		Handler: mux,
	}

	logger.Info("start server", slog.String("address", config.Address))
	if err := s.ListenAndServe(); err != nil {
		logger.Error("server run error", slog.Any("error", err))
	}

}
