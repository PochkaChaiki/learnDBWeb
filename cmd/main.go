package main

import (
	apiserver "learnDB/internal/apiServer"
	"learnDB/internal/config"
	"learnDB/internal/controller"
	"learnDB/internal/service"
	"learnDB/internal/storage"
	"learnDB/internal/storage/sqlite"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	config := config.MustLoad()

	db := sqlx.MustConnect("sqlite3", config.StoragePath)
	defer db.Close()
	storage := storage.Storage{
		AnswerStorage:   sqlite.NewAnswerStorage(db),
		DBStorage:       sqlite.NewDBStorage(db),
		DBSampleStorage: sqlite.NewDBSampleStorage(db),
		QueryStorage:    sqlite.NewQueryStorage(db),
		QuestionStorage: sqlite.NewQuestionStorage(db),
	}

	service := service.New(&storage)
	controller := controller.New(service)
	server := apiserver.New(config.Address, controller)

	server.Run()

}
