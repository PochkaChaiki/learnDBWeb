package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/pochkachaiki/learndb/internal/api/controller"
	"github.com/pochkachaiki/learndb/internal/config"
	"github.com/pochkachaiki/learndb/internal/dbmgr"
	"github.com/pochkachaiki/learndb/internal/testreview"
)

func main() {
	config := config.MustLoad()

	logFile, err := os.Create("./logfiles")
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	logger := slog.New(slog.NewJSONHandler(logFile, nil))

	manager := dbmgr.New()

	for name, connStr := range config.Databases {
		logger.Info("connecting to database", slog.String("db", name))
		err := manager.AddRepository(name, connStr)
		if err != nil {
			logger.Warn("cannot create db repository", slog.String("db", name), slog.Any("error", err))
		}
	}

	reviewer := testreview.New(manager, "olympics")

	controller := controller.NewTestReviewerController(logger, reviewer)

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
