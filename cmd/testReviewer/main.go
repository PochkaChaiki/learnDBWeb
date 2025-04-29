package main

import (
	"learnDB/internal/config"
	"learnDB/internal/dbManager"
	"learnDB/internal/testReviewer"
	trc "learnDB/internal/testReviewer/config"
	"learnDB/internal/testReviewer/repository/excel"
	"log/slog"

	"github.com/jmoiron/sqlx"
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
	for _, database := range config.Databases {
		db, err := sqlx.Connect(database.Name, database.ConnectionString)
		if err != nil {
			slog.Warn("cannot connect to db", slog.String("db", database.Name), slog.Any("error", err))
		}
		var manager testReviewer.DBRepository
		switch database.Name {
		case "mysql":
			manager, err = dbManager.NewMySQL(db)
		case "postgres":
			manager, err = dbManager.NewPostgreSQL(db)
		default:
			slog.Warn("database unrecognisable", slog.String("db", database.Name))
			continue
		}
		if err != nil {
			slog.Warn("cannot connect to db", slog.String("db", database.Name), slog.Any("error", err))
		}
		repo[database.Name] = manager
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
