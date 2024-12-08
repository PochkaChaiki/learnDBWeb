package main

import (
	"learnDB/internal/config"
	"learnDB/internal/models/db"
	"learnDB/internal/models/dbSample"
	"learnDB/internal/models/user"
	"learnDB/internal/storage"
	"log"
)

func main() {

	config := config.MustLoad()
	storage := storage.New(config.StoragePath)
	if err := SeedData(storage); err != nil {
		log.Fatalf(err.Error())
	}

}

func SeedData(st *storage.Storage) error {
	users := []user.User{
		{Username: "admin", Password: "password"},
		{Username: "user", Password: "qwerty"},
	}
	for _, u := range users {
		if err := u.Insert(st); err != nil {
			log.Fatalf("seed data error: %s", err)
		}
	}

	db := db.DB{Name: "sqlite3"}
	if err := db.Insert(st); err != nil {
		log.Fatalf("seed data error: %s", err)
	}
	dbSample := dbSample.DBSample{Description: "some db", Filepath: "/home/pochka/", DBId: 1}
	if err := dbSample.Insert(st); err != nil {
		log.Fatalf("seed data error: %s", err)
	}
	return nil
}
