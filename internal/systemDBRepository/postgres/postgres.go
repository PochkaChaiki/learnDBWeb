package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"learnDB/internal/domain"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// 1. Reading data

// Read Test
func (r *Repository) GetTest(name string) (*domain.Test, error) {
	const query = `SELECT id, name FROM testing_info.test WHERE name = $1`
	test := new(domain.Test)

	if err := r.db.Get(test, query, name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("test not found")
		}
		return nil, fmt.Errorf("failed to get test: %w", err)
	}
	// Get categories for test
	cats, err := r.GetCategoriesByTestId(test.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test: %w", err)
	}

	test.Categories = cats

	return test, nil
}

// Read category
func (r *Repository) getCategoryDependencies(cat *domain.Category) error {
	tasks, err := r.GetTasksByCategoryID(cat.ID)
	if err != nil {
		return err
	}
	cat.Tasks = tasks
	return nil
}

func (r *Repository) GetCategory(catID int) (*domain.Category, error) {
	const query = `SELECT id, name FROM testing_info.category WHERE id = $1`
	cat := new(domain.Category)

	if err := r.db.Get(cat, query, catID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	if err := r.getCategoryDependencies(cat); err != nil {
		return nil, err
	}

	return cat, nil
}

func (r *Repository) GetCategoriesByTestId(testID int) ([]domain.Category, error) {
	const query = `SELECT id, name FROM testing_info.category c 
	JOIN testing_info.test_category tc on c.id = tc.category_id
	WHERE tc.test_id = $1`
	var categories []domain.Category

	if err := r.db.Select(&categories, query, testID); err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	for i := range categories {
		if err := r.getCategoryDependencies(&categories[i]); err != nil {
			return nil, err
		}
	}
	return categories, nil
}

// Read Task

func (r *Repository) getTaskDependencies(task *domain.Task) error {
	// Read answers
	answers, err := r.GetAnswersByTaskID(task.ID)
	if err != nil {
		return fmt.Errorf("failed to get answers: %w", err)
	}
	task.Answers = answers

	// Read user groups
	groups, err := r.GetUserGroupsByTaskID(task.ID)
	if err != nil {
		return fmt.Errorf("failed to get user groups: %w", err)
	}
	task.UserGroups = groups

	// Read db samples
	dbSamples, err := r.GetDBSamplesByTaskID(task.ID)
	if err != nil {
		return fmt.Errorf("failed to get db samples: %w", err)
	}
	task.DBSamples = dbSamples
	return nil
}

func (r *Repository) GetTaskByName(name string) (*domain.Task, error) {
	const query = `SELECT id, name, question, points, general_feedback, correct_answer_jsonb 
                  FROM testing_info.task WHERE name = $1`
	task := new(domain.Task)

	if err := r.db.Get(task, query, name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	if err := r.getTaskDependencies(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (r *Repository) GetTasksByTestID(testID int) ([]domain.Task, error) {
	const query = `SELECT t.id, t.name, t.question, t.points, t.general_feedback, t.correct_answer_jsonb
                  FROM testing_info.task t
                  JOIN testing_info.category_task ct ON t.id = ct.task_id
				  JOIN testing_info.category c on ct.category_id = c.id
				  JOIN testing_info.test_category tc on c.id = tc.category_id
                  WHERE tc.test_id = $1`

	var tasks []domain.Task
	err := r.db.Select(&tasks, query, testID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	for i := range tasks {
		if err := r.getTaskDependencies(&tasks[i]); err != nil {
			return tasks, nil
		}
	}
	return tasks, nil
}

func (r *Repository) GetTasksByCategoryID(catID int) ([]domain.Task, error) {
	const query = `SELECT t.id, t.name, t.question, t.points, t.general_feedback, t.correct_answer_jsonb
    FROM testing_info.task t
    JOIN testing_info.category_task ct ON t.id = ct.task_id
    WHERE ct.category_id = $1`
	var tasks []domain.Task
	if err := r.db.Select(&tasks, query, catID); err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	for i := range tasks {
		if err := r.getTaskDependencies(&tasks[i]); err != nil {
			return tasks, nil
		}
	}
	return tasks, nil
}

// Get single answer by id
func (r *Repository) GetAnswerByID(answerID int) (*domain.Answer, error) {
	const query = `SELECT id, text, feedback 
                  FROM testing_info.answer WHERE id = $1`
	answer := new(domain.Answer)

	if err := r.db.Get(answer, query, answerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("answer not found")
		}
		return nil, fmt.Errorf("failed to get answer: %w", err)
	}

	return answer, nil
}

func (r *Repository) GetAnswersByTaskID(taskID int) ([]domain.Answer, error) {
	const query = `SELECT a.id, a.text, a.feedback
                  FROM testing_info.answer a
                  WHERE a.task_id = $1`

	var answers []domain.Answer
	err := r.db.Select(&answers, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get answers: %w", err)
	}
	return answers, nil
}

func (r *Repository) GetUserGroupsByTaskID(taskID int) ([]string, error) {
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

func (r *Repository) GetDBSamplesByTaskID(taskID int) ([]domain.DBSample, error) {
	const query = `SELECT ds.id, ds.name, ds.path, ds.description, ds.dbms
                  FROM testing_info.db_sample ds
                  JOIN testing_info.task_db_sample tds ON ds.id = tds.db_sample_id
                  WHERE tds.task_id = $1`

	var samples []domain.DBSample
	err := r.db.Select(&samples, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get db samples: %w", err)
	}

	return samples, nil
}

func (r *Repository) GetDBMSList() ([]domain.DBMS, error) {
	const query = `SELECT name FROM testing_info.dbms`

	var dbmsList []domain.DBMS
	err := r.db.Select(&dbmsList, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get dbms list: %w", err)
	}

	return dbmsList, nil
}

func (r *Repository) GetUserGroups() ([]domain.UserGroup, error) {
	const query = `SELECT name FROM user_management.user_groups`

	var groups []domain.UserGroup
	err := r.db.Select(&groups, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}

	return groups, nil
}

// Create data

func (r *Repository) CreateTest(test *domain.Test) error {
	const query = `INSERT INTO testing_info.test (name) VALUES ($1) RETURNING id`

	err := r.db.QueryRowx(query, test.Name).Scan(&test.ID)
	if err != nil {
		return fmt.Errorf("failed to create test: %w", err)
	}

	for _, cat := range test.Categories {
		err := r.AddCategoryToTest(test.ID, cat.ID)
		if err != nil {
			return fmt.Errorf("failed to add task to test: %w", err)
		}
	}

	return nil
}

func (r *Repository) AddCategoryToTest(testID, catID int) error {
	const query = `INSERT INTO testing_info.test_category (test_id, category_id) VALUES ($1, $2)`

	_, err := r.db.Exec(query, testID, catID)
	if err != nil {
		return fmt.Errorf("failed to add category to test: %w", err)
	}

	return nil
}

func (r *Repository) CreateCategory(cat *domain.Category) error {
	const query = `INSERT INTO testing_info.category (name) VALUES ($1) RETURNING id`
	if err := r.db.QueryRowx(query, cat.Name).Scan(&cat.ID); err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	for _, task := range cat.Tasks {
		if err := r.AddTaskToCategory(cat.ID, task.ID); err != nil {
			return fmt.Errorf("failed to add task to category: %w", err)
		}
	}
	return nil
}

func (r *Repository) AddTaskToCategory(catID, taskID int) error {
	const query = `INSERT INTO testing_info.category_task (category_id, task_id) VALUES ($1, $2)`

	_, err := r.db.Exec(query, catID, taskID)
	if err != nil {
		return fmt.Errorf("failed to add task to category: %w", err)
	}

	return nil
}

func (r *Repository) CreateTask(task *domain.Task) error {
	const query = `INSERT INTO testing_info.task 
                  (name, question, points, general_feedback, correct_answer_jsonb) 
                  VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.QueryRowx(query,
		task.Name, task.Question, task.Points, task.GeneralFeedback, task.CorrectAnswerJSONB).Scan(&task.ID)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	for _, answer := range task.Answers {
		err := r.AddAnswerToTask(task.ID, answer.ID)
		if err != nil {
			return fmt.Errorf("failed to add answer to task: %w", err)
		}
	}

	for _, group := range task.UserGroups {
		err := r.AddUserGroupToTask(group, task.ID)
		if err != nil {
			return fmt.Errorf("failed to add user group to task: %w", err)
		}
	}

	for _, sample := range task.DBSamples {
		err := r.AddDBSampleToTask(task.ID, sample.Name)
		if err != nil {
			return fmt.Errorf("failed to add db sample to task: %w", err)
		}
	}

	return nil
}

func (r *Repository) CreateAnswer(answer *domain.Answer) error {
	const query = `INSERT INTO testing_info.answer 
                  (text, feedback) 
                  VALUES ($1, $2, $3) RETURNING id`

	err := r.db.QueryRowx(query,
		answer.Text, answer.Feedback).Scan(&answer.ID)
	if err != nil {
		return fmt.Errorf("failed to create answer: %w", err)
	}

	return nil
}

func (r *Repository) AddAnswerToTask(taskID, answerID int) error {
	const query = `INSERT INTO testing_info.task_answer (task_id, answer_id) VALUES ($1, $2)`

	_, err := r.db.Exec(query, taskID, answerID)
	if err != nil {
		return fmt.Errorf("failed to add answer to task: %w", err)
	}

	return nil
}

func (r *Repository) AddUserGroupToTask(group string, taskID int) error {
	const query = `INSERT INTO testing_info.user_groups_task (user_group, task_id) VALUES ($1, $2)`

	_, err := r.db.Exec(query, group, taskID)
	if err != nil {
		return fmt.Errorf("failed to add user group to task: %w", err)
	}

	return nil
}

func (r *Repository) AddDBSampleToTask(taskID int, sampleName string) error {
	const query = `INSERT INTO testing_info.task_db_sample (task_id, db_sample_name) VALUES ($1, $2)`

	_, err := r.db.Exec(query, taskID, sampleName)
	if err != nil {
		return fmt.Errorf("failed to add db sample to task: %w", err)
	}

	return nil
}

// 3. Update data

func (r *Repository) UpdateTest(test *domain.Test) error {
	const query = `UPDATE testing_info.test SET name = $1 WHERE id = $2`

	_, err := r.db.Exec(query, test.Name, test.ID)
	if err != nil {
		return fmt.Errorf("failed to update test: %w", err)
	}

	return nil
}

func (r *Repository) UpdateCategory(cat *domain.Category) error {
	const query = `UPDATE testing_info.category SET name = $1 WHERE id = $2`
	if _, err := r.db.Exec(query, cat.Name, cat.ID); err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

func (r *Repository) UpdateTask(task *domain.Task) error {
	const query = `UPDATE testing_info.task 
                  SET name = $1, question = $2, points = $3, general_feedback = $4, correct_answer_jsonb = $5
                  WHERE id = $6`

	_, err := r.db.Exec(query,
		task.Name, task.Question, task.Points, task.GeneralFeedback, task.CorrectAnswerJSONB, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

func (r *Repository) UpdateAnswer(answer *domain.Answer) error {
	const query = `UPDATE testing_info.answer 
                  SET text = $1, feedback = $2
                  WHERE id = $3`

	_, err := r.db.Exec(query,
		answer.Text, answer.Feedback, answer.ID)
	if err != nil {
		return fmt.Errorf("failed to update answer: %w", err)
	}

	return nil
}

// 4. Remove data

func (r *Repository) DeleteTest(testID int) error {
	const query = `DELETE FROM testing_info.test WHERE id = $1`

	_, err := r.db.Exec(query, testID)
	if err != nil {
		return fmt.Errorf("failed to delete test: %w", err)
	}

	return nil
}

func (r *Repository) DeleteCategory(catID int) error {
	const query = `DELETE FROM testing_info.category WHERE id = $1`

	_, err := r.db.Exec(query, catID)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

func (r *Repository) DeleteTask(taskID int) error {
	const query = `DELETE FROM testing_info.task WHERE id = $1`

	_, err := r.db.Exec(query, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

func (r *Repository) DeleteAnswer(answerID int) error {
	const query = `DELETE FROM testing_info.answer WHERE id = $1`

	_, err := r.db.Exec(query, answerID)
	if err != nil {
		return fmt.Errorf("failed to delete answer: %w", err)
	}

	return nil
}

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

func (r *Repository) CreateDBSample(sample *domain.DBSample) error {
	const query = `INSERT INTO testing_info.db_sample 
                  (name, path, description, dbms) 
                  VALUES ($1, $2, $3, $4)`

	_, err := r.db.Exec(query,
		sample.Name, sample.Path, sample.Description, sample.DBMS)
	if err != nil {
		return fmt.Errorf("failed to create db sample: %w", err)
	}

	return nil
}

func (r *Repository) UpdateDBSample(sample *domain.DBSample) error {
	const query = `UPDATE testing_info.db_sample 
                  SET path = $1, description = $2, dbms = $3
                  WHERE name = $4`

	_, err := r.db.Exec(query,
		sample.Path, sample.Description, sample.DBMS, sample.Name)
	if err != nil {
		return fmt.Errorf("failed to update db sample: %w", err)
	}

	return nil
}

func (r *Repository) DeleteDBSample(name string) error {
	const query = `DELETE FROM testing_info.db_sample WHERE name = $1`

	_, err := r.db.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to delete db sample: %w", err)
	}

	return nil
}

func (r *Repository) CreateUserGroup(name string) error {
	const query = `INSERT INTO user_management.user_groups (name) VALUES ($1)`

	_, err := r.db.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to create user group: %w", err)
	}

	return nil
}

func (r *Repository) DeleteUserGroup(name string) error {
	const query = `DELETE FROM user_management.user_groups WHERE name = $1`

	_, err := r.db.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to delete user group: %w", err)
	}

	return nil
}

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
