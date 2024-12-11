package domain

type DBSample struct {
	Id          int    `json:"-" db:"db_sample_id"`
	Description string `json:"description" db:"description"`
	Filepath    string `json:"-" db:"filepath"`
	DBId        int    `json:"-" db:"db_id"`
}
