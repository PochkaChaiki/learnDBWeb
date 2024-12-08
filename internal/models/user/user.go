package user

import (
	"fmt"
	"learnDB/internal/storage"
)

type User struct {
	Id       int    `json:"-" db:"user_id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

func (u *User) Insert(st *storage.Storage) error {
	_, err := st.DB.NamedExec("insert into user(username, password) values (:username, :password);", u)
	if err != nil {
		return fmt.Errorf("insert user error: %s", err)
	}
	return nil
}

func (u *User) Get(st *storage.Storage, id int) error {

	if err := st.DB.Get(u, "select user_id, username, password from user where user_id = ?;", id); err != nil {
		return fmt.Errorf("select user error: %s", err)
	}
	return nil
}

func GetAll(st *storage.Storage) ([]User, error) {

	users := []User{}

	if err := st.DB.Select(&users, "select user_id, username, password from user;"); err != nil {
		return nil, fmt.Errorf("select users error: %s", err)
	}
	return users, nil
}

func (u *User) Update(st *storage.Storage) error {
	_, err := st.DB.NamedExec("update user set username=:username, password=:password where user_id=:user_id", u)
	if err != nil {
		return fmt.Errorf("update user error: %s", err)
	}
	return nil
}

func (u *User) Delete(st *storage.Storage) error {
	if _, err := st.DB.Exec("delete from user where user_id = ?", u.Id); err != nil {
		return fmt.Errorf("delete user error: %s", err)
	}
	return nil
}
