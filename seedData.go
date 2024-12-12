package main

import (
	"learnDB/internal/config"
	"learnDB/internal/domain"
	"learnDB/internal/storage/sqlite"
	"learnDB/internal/utils"
	"sync"

	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Storage[T any] interface {
	Insert(*T) error
	Get(int) (*T, error)
	GetAll() ([]T, error)
	Update(*T) error
	Delete(int) error
}

type SafeSeeding struct {
	mu sync.Mutex
	wg sync.WaitGroup
}

func main() {

	config := config.MustLoad()

	db := sqlx.MustConnect("sqlite3", config.StoragePath)

	sf := &SafeSeeding{mu: sync.Mutex{}, wg: sync.WaitGroup{}}

	sf.wg.Add(1)
	go sf.SeedUser(sqlite.NewUserStorage(db), config.Salt)

	sf.wg.Add(1)
	go sf.SeedDB(sqlite.NewDBStorage(db))

	sf.wg.Add(1)
	go sf.SeedDBSample(sqlite.NewDBSampleStorage(db))

	sf.wg.Wait()
}

func (sf *SafeSeeding) SeedUser(s Storage[domain.User], salt string) {

	users := []domain.User{
		{Username: "admin", Password: OmitErrorString(utils.SaltAndHashString("password", salt))},
		{Username: "user", Password: OmitErrorString(utils.SaltAndHashString("qwerty", salt))},
	}
	for _, u := range users {
		if err := s.Insert(&u); err != nil {
			sf.mu.Lock()
			log.Fatalf("seed data error: %s", err)
			sf.mu.Unlock()
		}
	}
	sf.wg.Done()
}

func (sf *SafeSeeding) SeedDB(s Storage[domain.DB]) {
	db := domain.DB{Name: "sqlite3"}
	if err := s.Insert(&db); err != nil {
		sf.mu.Lock()
		log.Fatalf("seed db error: %s", err)
		sf.mu.Unlock()
	}
	sf.wg.Done()
}

func (sf *SafeSeeding) SeedDBSample(s Storage[domain.DBSample]) {
	dbSample := domain.DBSample{Description: "some db", Filepath: "/home/pochka/", DBId: 1}
	if err := s.Insert(&dbSample); err != nil {
		sf.mu.Lock()
		log.Fatalf("seed data error: %s", err)
		sf.mu.Unlock()
	}
	sf.wg.Done()
}

func OmitErrorString(i string, err error) string {
	return i
}
