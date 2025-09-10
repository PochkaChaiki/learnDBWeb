package postgres

import (
	"context"
	"fmt"

	"github.com/pochkachaiki/learndb/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

var txOptions = pgx.TxOptions{
	IsoLevel:       pgx.ReadCommitted,
	AccessMode:     pgx.ReadWrite,
	DeferrableMode: pgx.Deferrable,
}

func (r *Repository) addRelations(ctx context.Context, query string, entityID int, idsToAdd []int) error {

	if len(idsToAdd) != 0 {

		batch := &pgx.Batch{}

		for _, id := range idsToAdd {
			batch.Queue(query, entityID, id).Exec(func(ct pgconn.CommandTag) error {
				if ct.RowsAffected() == 0 {
					return fmt.Errorf("batch insert error: %w", ct.String())
				}
				return nil
			})
		}

		if err := r.db.SendBatch(ctx, batch).Close(); err != nil {
			return fmt.Errorf("failed to add relation: %w", err)
		}
	}
	return nil
}

func (r *Repository) ReadSchemasOfDBMS(dbms string) ([]string, error) {
	const query = `SELECT ds.name
                  FROM testing_info.db_sample ds
                  WHERE ds.dbms = $1`

	var samples []string
	err := r.db.Select(&samples, query, dbms)
	if err != nil {
		return nil, fmt.Errorf("failed to get db samples: %w", err)
	}

	return samples, nil
}

func (r *Repository) ReadDBMSList() ([]domain.DBMS, error) {
	const query = `SELECT name FROM testing_info.dbms`

	var dbmsList []domain.DBMS
	err := r.db.Select(&dbmsList, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get dbms list: %w", err)
	}

	return dbmsList, nil
}

// Create data

// func (r *Repository) AddCategoryToTest(testID, catID int) error {
// 	const query = `INSERT INTO testing_info.test_category (test_id, category_id) VALUES ($1, $2)`

// 	_, err := r.db.Exec(query, testID, catID)
// 	if err != nil {
// 		return fmt.Errorf("failed to add category to test: %w", err)
// 	}

// 	return nil
// }

func (r *Repository) AddTaskToCategory(catID, taskID int) error {
	const query = `INSERT INTO testing_info.category_task (category_id, task_id) VALUES ($1, $2)`

	_, err := r.db.Exec(query, catID, taskID)
	if err != nil {
		return fmt.Errorf("failed to add task to category: %w", err)
	}

	return nil
}

// 3. Update data

// 4. Remove data

func (r *Repository) RemoveCategoryFromTest(testID, catID int) error {
	const query = `DELETE FROM testing_info.test_category WHERE test_id = $1 AND category_id = $2`

	_, err := r.db.Exec(query, testID, catID)
	if err != nil {
		return fmt.Errorf("failed to remove category from test: %w", err)
	}

	return nil
}

func (r *Repository) RemoveTaskFromCategory(catID, taskID int) error {
	const query = `DELETE FROM testing_info.category_task WHERE cat_id = $1 AND task_id = $2`

	_, err := r.db.Exec(query, catID, taskID)
	if err != nil {
		return fmt.Errorf("failed to remove task from category: %w", err)
	}

	return nil
}

func (r *Repository) RemoveAnswerFromTask(taskID, answerID int) error {
	const query = `DELETE FROM testing_info.task_answer WHERE task_id = $1 AND answer_id = $2`

	_, err := r.db.Exec(query, taskID, answerID)
	if err != nil {
		return fmt.Errorf("failed to remove answer from task: %w", err)
	}

	return nil
}

func (r *Repository) RemoveUserGroupFromTask(group string, taskID int) error {
	const query = `DELETE FROM testing_info.user_groups_task WHERE user_group = $1 AND task_id = $2`

	_, err := r.db.Exec(query, group, taskID)
	if err != nil {
		return fmt.Errorf("failed to remove user group from task: %w", err)
	}

	return nil
}

func (r *Repository) RemoveDBSampleFromTask(taskID int, sampleName string) error {
	const query = `DELETE FROM testing_info.task_db_sample WHERE task_id = $1 AND db_sample_name = $2`

	_, err := r.db.Exec(query, taskID, sampleName)
	if err != nil {
		return fmt.Errorf("failed to remove db sample from task: %w", err)
	}

	return nil
}

func (r *Repository) AddUserToGroup(userID int, group string) error {
	const query = `INSERT INTO user_management.user_user_groups (user_id, user_group) VALUES ($1, $2)`

	_, err := r.db.Exec(query, userID, group)
	if err != nil {
		return fmt.Errorf("failed to add user to group: %w", err)
	}

	return nil
}

func (r *Repository) RemoveUserFromGroup(userID int, group string) error {
	const query = `DELETE FROM user_management.user_user_groups WHERE user_id = $1 AND user_group = $2`

	_, err := r.db.Exec(query, userID, group)
	if err != nil {
		return fmt.Errorf("failed to remove user from group: %w", err)
	}

	return nil
}
