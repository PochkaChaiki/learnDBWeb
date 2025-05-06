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

// const (
// 	readPath  = "/home/pochka/projects/learnDB/static/excelConfig.json"
// 	readSh    = "Sheet1"
// 	writePath = "/home/pochka/projects/learnDB/static/reviewedWorks.xlsx"
// 	writeSh   = "Sheet1"
// )

func main() {
	config := config.MustLoad()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

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

	s := http.Server{
		Addr:    config.Address,
		Handler: mux,
	}

	if err := s.ListenAndServe(); err != nil {
		logger.Error("server run error", slog.Any("error", err))
	} else {
		logger.Info("server is running", slog.String("address", config.Address))
	}

}
