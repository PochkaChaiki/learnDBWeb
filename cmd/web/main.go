package main

import (
	"learnDB/internal/config"
	apiserver "learnDB/internal/web/app/apiServer"
	"learnDB/internal/web/controller"
	"learnDB/internal/web/service"
	"learnDB/internal/web/storage"
	"learnDB/internal/web/storage/sqlite"

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
		UserStorage:     sqlite.NewUserStorage(db),
	}

	apiService := service.New(&storage)
	apiController := controller.New(apiService)
	authService := service.NewAuthService(storage.UserStorage, config.Salt, []byte(config.SecretKey), config.ExpirationTime, config.AdminCredential)
	authController := controller.NewAuthController(authService)
	server := apiserver.New(config.Address, []byte(config.SecretKey), apiController, authController)

	server.Run()

}
