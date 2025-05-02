package main

import (
	"learnDB/internal/config"
	"learnDB/internal/dbManager"
	"learnDB/internal/testReviewer"
	trc "learnDB/internal/testReviewer/config"
	"learnDB/internal/testReviewer/repository/excel"
	"log/slog"
)

const (
	readPath  = "/home/pochka/projects/learnDB/static/excelConfig.json"
	readSh    = "Sheet1"
	writePath = "/home/pochka/projects/learnDB/static/reviewedWorks.xlsx"
	writeSh   = "Sheet1"
)

func main() {
	config := config.MustLoad()
	testReviewerConf := trc.MustLoadConfig("./static/excelConfig.json")

	repo := make(map[string]testReviewer.DBRepository)
	for name, connStr := range config.Databases {
		manager, err := dbManager.New(name, connStr)
		if err != nil {
			slog.Warn("cannot create dbManager to db", slog.String("db", name), slog.Any("error", err))
		}

		repo[name] = manager
	}

	reviewer := testReviewer.New(repo)

	exrepo := excel.New(readPath, readSh, writePath, writeSh)

	studentWorks, err := exrepo.Read(testReviewerConf)
	if err != nil {
		slog.Error("cannot read excel", slog.Any("error", err))
		return
	}

	reviewedWorks := reviewer.CheckTest(studentWorks)

	exrepo.Write(reviewedWorks)
}
