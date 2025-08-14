package main

import (
	"log/slog"
	"os"

	"github.com/pochkachaiki/learndb/internal/api"
	"github.com/pochkachaiki/learndb/internal/config"
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
