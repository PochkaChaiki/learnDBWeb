package domain

type DB struct {
	Id   int    `json:"-" db:"db_id"`
	Name string `json:"db_name" db:"db_name"`
}
