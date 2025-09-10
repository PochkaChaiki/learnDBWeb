package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/pochkachaiki/learndb/internal/domain"
)

func (r *Repository) GetUser(userID int) (*domain.User, error) {
	const query = `SELECT id, login, password FROM user_management.user WHERE id = $1`
	var user domain.User

	err := r.db.Get(&user, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	groups, err := r.GetUserGroupsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}
	user.Groups = groups

	return &user, nil
}

func (r *Repository) GetUserGroupsByUserID(userID int) ([]string, error) {
	const query = `SELECT user_group 
                  FROM user_management.user_user_groups 
                  WHERE user_id = $1`

	var groups []string
	err := r.db.Select(&groups, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}

	return groups, nil
}

func (r *Repository) CreateUser(user *domain.User) error {
	const query = `INSERT INTO user_management.user (login, password) VALUES ($1, $2) RETURNING id`

	err := r.db.QueryRowx(query, user.Login, user.Password).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	for _, group := range user.Groups {
		err := r.AddUserToGroup(user.ID, group)
		if err != nil {
			return fmt.Errorf("failed to add user to group: %w", err)
		}
	}

	return nil
}

func (r *Repository) UpdateUser(user *domain.User) error {
	const query = `UPDATE user_management.user SET login = $1, password = $2 WHERE id = $3`

	_, err := r.db.Exec(query, user.Login, user.Password, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *Repository) DeleteUser(userID int) error {
	const query = `DELETE FROM user_management.user WHERE id = $1`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
