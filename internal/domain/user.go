package domain

type User struct {
	Id       int    `json:"user_id,omitempty" db:"user_id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}
