package main

import (
	"learnDB/internal/api"
	"learnDB/internal/config"
	"log/slog"
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

	srv := api.New(config, logger)
	srv.Run()
}
