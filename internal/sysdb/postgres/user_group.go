package postgres

import (
	"fmt"

	"github.com/pochkachaiki/learndb/internal/domain"
)

func (r *Repository) CreateUserGroup(name string) error {
	const query = `INSERT INTO user_management.user_groups (name) VALUES ($1)`

	_, err := r.db.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to create user group: %w", err)
	}

	return nil
}

func (r *Repository) ReadUserGroupsByTaskID(taskID int) ([]string, error) {
	const query = `SELECT user_group 
                  FROM testing_info.user_groups_task 
                  WHERE task_id = $1`

	var groups []string
	err := r.db.Select(&groups, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}

	return groups, nil
}

func (r *Repository) ReadUserGroups() ([]domain.UserGroup, error) {
	const query = `SELECT name FROM user_management.user_groups`

	var groups []domain.UserGroup
	err := r.db.Select(&groups, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}

	return groups, nil
}

func (r *Repository) DeleteUserGroup(name string) error {
	const query = `DELETE FROM user_management.user_groups WHERE name = $1`

	_, err := r.db.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to delete user group: %w", err)
	}

	return nil
}
