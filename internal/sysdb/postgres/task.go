package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pochkachaiki/learndb/internal/domain"
)

// Task CRUD

// Method to create `Task` in db.
//
// Pass `Task` structure with empty `UserGroups` and `DBSamples` slice.
// Pass user groups' and db sample's ids alone.
func (r *Repository) CreateTask(ctx context.Context, task *domain.Task, ug []int, dbs []int) error {

	tx, err := r.db.BeginTx(ctx, txOptions)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	query := `INSERT INTO testing_info.task 
			  (name, question, points, general_feedback, correct_answer_jsonb)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`

	if err := tx.QueryRow(
		ctx,
		query,
		task.Name,
		task.Question,
		task.Points,
		task.GeneralFeedback,
		task.CorrectAnswerJSONB).Scan(&task.ID); err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	for _, ans := range task.XMLAnswers {
		if err := r.createXMLAnswer(ctx, tx, &ans); err != nil {
			return fmt.Errorf("failed to insert answer: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) AddUserGroupsToTask(ctx context.Context, taskID int, ids []int) error {
	query := `INSERT INTO testing_info.user_group_task (task_id, user_group_id) VALUES ($1, $2)`
	if err := r.addRelations(ctx, query, taskID, ids); err != nil {
		return fmt.Errorf("failed to add user_groups: %w", err)
	}
	return nil
}

func (r *Repository) AddDBSampleToTask(ctx context.Context, taskID int, ids []int) error {
	query := `INSERT INTO testing_info.task_db_sample (task_id, db_sample_id) VALUES ($1, $2)`
	if err := r.addRelations(ctx, query, taskID, ids); err != nil {
		return fmt.Errorf("failed to add db_samples: %w", err)
	}
	return nil
}

func (r *Repository) getTaskDependencies(ctx context.Context, task *domain.Task) error {
	// Read answers
	answers, err := r.ReadAnswersByTaskID(ctx, task.ID)
	if err != nil {
		return fmt.Errorf("failed to get answers: %w", err)
	}
	task.XMLAnswers = answers

	// Read user groups
	groups, err := r.ReadUserGroupsByTaskID(ctx, task.ID)
	if err != nil {
		return fmt.Errorf("failed to get user groups: %w", err)
	}
	task.UserGroups = groups

	// Read db samples
	dbSamples, err := r.ReadDBSamplesByTaskID(ctx, task.ID)
	if err != nil {
		return fmt.Errorf("failed to get db samples: %w", err)
	}
	task.DBSamples = dbSamples
	return nil
}

func (r *Repository) ReadTaskByID(ctx context.Context, id int) (*domain.Task, error) {
	const query = `SELECT id, name, question, points, general_feedback, correct_answer_jsonb 
                  FROM testing_info.task WHERE id = $1`
	task := new(domain.Task)

	if err := r.db.QueryRow(ctx, query, id).Scan(task); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to read task: %w", err)
	}

	if err := r.getTaskDependencies(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (r *Repository) ReadTaskByName(ctx context.Context, name string) (*domain.Task, error) {
	const query = `SELECT id, name, question, points, general_feedback, correct_answer_jsonb 
                  FROM testing_info.task WHERE name = $1`
	task := new(domain.Task)

	if err := r.db.QueryRow(ctx, query, name).Scan(task); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to read task: %w", err)
	}

	if err := r.getTaskDependencies(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (r *Repository) ReadTasksByTestID(ctx context.Context, testID int) ([]domain.Task, error) {
	const query = `SELECT t.id, t.name, t.question, t.points, t.general_feedback, t.correct_answer_jsonb
                  FROM testing_info.task t
                  JOIN testing_info.category_task ct ON t.id = ct.task_id
				  JOIN testing_info.category c on ct.category_id = c.id
				  JOIN testing_info.test_category tc on c.id = tc.category_id
                  WHERE tc.test_id = $1`

	rows, err := r.db.Query(ctx, query, testID)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks: %w", err)
	}
	tasks, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Task])
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks: %w", err)
	}

	for i := range tasks {
		if err := r.getTaskDependencies(ctx, &tasks[i]); err != nil {
			return nil, err
		}
	}
	return tasks, nil
}

func (r *Repository) ReadTasksByCategoryID(ctx context.Context, catID int) ([]domain.Task, error) {
	const query = `SELECT t.id, t.name, t.question, t.points, t.general_feedback, t.correct_answer_jsonb
    			   FROM testing_info.task t
    			   JOIN testing_info.category_task ct ON t.id = ct.task_id
    			   WHERE ct.category_id = $1`

	rows, err := r.db.Query(ctx, query, catID)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks: %w", err)
	}

	tasks, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Task])
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks: %w", err)
	}

	for i := range tasks {
		if err := r.getTaskDependencies(ctx, &tasks[i]); err != nil {
			return nil, err
		}
	}
	return tasks, nil
}

// Method to update `Task` in db.
//
// Pass `Task` structure with empty `UserGroups` and `DBSamples` slice.
// Pass user groups' and db sample's ids alone.
func (r *Repository) UpdateTask(ctx context.Context, task *domain.Task, ug []int, dbs []int) error {
	tx, err := r.db.BeginTx(ctx, txOptions)
	if err != nil {
		return fmt.Errorf("failed to update tasks: %w", err)
	}

	defer tx.Rollback(ctx)

	var args []any
	i := 1
	query := `UPDATE testing_info.task SET `
	if task.Name != nil {
		query += fmt.Sprintf("name = $%d, ", i)
		args = append(args, task.Name)
		i++
	}
	if task.Question != "" {
		query += fmt.Sprintf("question = $%d, ", i)
		args = append(args, task.Question)
		i++
	}
	if task.Points != 0 {
		query += fmt.Sprintf("points = $%d, ", i)
		args = append(args, task.Points)
		i++
	}
	if task.GeneralFeedback != nil {
		query += fmt.Sprintf("general_feedback = $%d, ", i)
		args = append(args, task.Question)
		i++
	}
	if task.CorrectAnswerJSONB != nil {
		query += fmt.Sprintf("correct_answer_jsonb = $%d, ", i)
		args = append(args, task.Question)
		i++
	}

	query += fmt.Sprintf("WHERE id = $%d", i)
	args = append(args, task.ID)

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	query = `SELECT id FROM testing_info.answer WHERE task_id = $1`
	rows, err := tx.Query(ctx, query, task.ID)
	if err != nil {
		return fmt.Errorf("failed to read answers: %w", err)
	}
	ansInDB, err := pgx.CollectRows(rows, pgx.RowTo[int])
	if err != nil {
		return fmt.Errorf("failed to read answers: %w", err)
	}

	var answers []int
	var ansToAdd []domain.XMLAnswer

	for _, ans := range task.XMLAnswers {
		if ans.ID == 0 {
			ansToAdd = append(ansToAdd, ans)
			continue
		}
		answers = append(answers, ans.ID)
	}

	ansToDelete := difference(ansInDB, answers)

	query = `DELETE FROM testing_info.answer WHERE id = ANY ($1)`

	if _, err := tx.Exec(ctx, query, ansToDelete); err != nil {
		return fmt.Errorf("failed to delete answers: %w", err)
	}

	for _, ans := range ansToAdd {
		if err := r.createAnswer(ctx, tx, &ans); err != nil {
			return fmt.Errorf("failed to create new answers: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

func (r *Repository) DeleteTask(ctx context.Context, taskID int) error {
	const query = `DELETE FROM testing_info.task WHERE id = $1`

	_, err := r.db.Exec(ctx, query, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}
