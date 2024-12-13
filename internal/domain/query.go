package domain

import (
	"math/rand"
	"time"
)

type Query struct {
	Id         int       `json:"query_id,omitempty" db:"query_id"`
	Script     string    `json:"script" db:"script"`
	Info       string    `json:"info,omitempty" db:"info"`
	ExecutedAt time.Time `json:"executed_at" db:"executed_at"`
	UserId     int       `json:"user_id,omitempty" db:"user_id"`
	DBId       int       `json:"db_id" db:"db_id"`
}

func (q *Query) RunQuery() {
	rand := rand.New(rand.NewSource(time.Now().Unix()))
	if rand.Int()%3 == 0 {
		q.Info = "Error: something bad happened"
		return
	}
	q.Info = "OK"
}
