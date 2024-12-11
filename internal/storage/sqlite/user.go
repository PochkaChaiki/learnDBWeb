package sqlite

import (
	"fmt"
	"learnDB/internal/domain"

	"github.com/jmoiron/sqlx"
)

type UserStorage struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (r *UserStorage) Insert(u *domain.User) error {
	_, err := r.db.NamedExec("insert into user(username, password) values (:username, :password);", u)
	if err != nil {
		return fmt.Errorf("insert user error: %s", err)
	}
	return nil
}

func (r *UserStorage) Get(id int) (*domain.User, error) {
	u := new(domain.User)
	if err := r.db.Get(u, "select user_id, username, password from user where user_id = ?;", id); err != nil {
		return nil, fmt.Errorf("select user error: %s", err)
	}
	return u, nil
}

func (r *UserStorage) GetAll() (users []domain.User, err error) {
	if err := r.db.Select(&users, "select user_id, username, password from user;"); err != nil {
		return nil, fmt.Errorf("select users error: %s", err)
	}
	return users, nil
}

func (r *UserStorage) Update(u *domain.User) error {
	_, err := r.db.NamedExec("update user set username=:username, password=:password where user_id=:user_id", u)
	if err != nil {
		return fmt.Errorf("update user error: %s", err)
	}
	return nil
}

func (r *UserStorage) Delete(id int) error {
	if _, err := r.db.Exec("delete from user where user_id = ?", id); err != nil {
		return fmt.Errorf("delete user error: %s", err)
	}
	return nil
}
