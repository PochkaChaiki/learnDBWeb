package domain

import (
	"time"
)

type Query struct {
	Id         int       `json:"-" db:"query_id"`
	Script     string    `json:"script" db:"script"`
	Info       string    `json:"info" db:"info"`
	ExecutedAt time.Time `json:"executed_at" db:"executed_at"`
	UserId     int       `json:"-" db:"user_id"`
	DBId       int       `json:"-" db:"db_id"`
}
