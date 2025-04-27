package main

import (
	"learnDB/internal/config"
	mysql "learnDB/internal/dbRepository/MySQL"
	"learnDB/internal/testReviewer"
	trc "learnDB/internal/testReviewer/config"

	"github.com/jmoiron/sqlx"
)

func main() {
	config := config.MustLoad()
	trConfig := trc.MustLoadConfig("./static/excelConfig.json")

	db, err := sqlx.Connect("mysql", config.AdminCredential)
	repo, err := mysql.NewDB()

	reviewer := testReviewer.New()
}
