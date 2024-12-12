package domain

import (
	"time"
)

type Query struct {
	Id         int       `json:"query_id,omitempty" db:"query_id"`
	Script     string    `json:"script" db:"script"`
	Info       string    `json:"info" db:"info"`
	ExecutedAt time.Time `json:"executed_at" db:"executed_at"`
	UserId     int       `json:"user_id" db:"user_id"`
	DBId       int       `json:"db_id" db:"db_id"`
}
