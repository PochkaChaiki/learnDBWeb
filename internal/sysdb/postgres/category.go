package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pochkachaiki/learndb/internal/domain"
)

// Category CRUD

// Method to create `Category` in db.
//
// Pass `Category` structure with empty `Tasks` slice.
func (r *Repository) CreateCategory(ctx context.Context, cat *domain.Category) error {
	var query = `INSERT INTO testing_info.category (name) VALUES ($1) RETURNING id`
	if err := r.db.QueryRow(ctx, query, cat.Name).Scan(&cat.ID); err != nil {
		return fmt.Errorf("failed to create test: %w", err)
	}
	return nil
}

func (r *Repository) AddTasksToCategory(ctx context.Context, catID int, ids []int) error {
	query := `INSERT INTO testing_info.category_task (category_id, task_id) VALUES ($1, $2)`
	if err := r.addRelations(ctx, query, catID, ids); err != nil {
		return fmt.Errorf("failed to add tasks to category: %w", err)
	}
	return nil
}

func (r *Repository) RemoveTasksFromCategory(ctx context.Context, catID int, ids []int) error {
	query := `DELETE FROM testing_info.category_task WHERE category_id = $1 AND task_id = ANY ($2)`
	if _, err := r.db.Exec(ctx, query, catID, ids); err != nil {
		return fmt.Errorf("failed to remove tasks from category: %w", err)
	}
	return nil
}

func (r *Repository) ReadCategoryByID(ctx context.Context, id int) (*domain.Category, error) {
	const query = `SELECT id, name FROM testing_info.category WHERE id = $1`
	cat := new(domain.Category)

	if err := r.db.QueryRow(ctx, query, id).Scan(cat); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	tasks, err := r.ReadTasksByCategoryID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	cat.Tasks = tasks

	return cat, nil
}

func (r *Repository) ReadCategoryByName(ctx context.Context, name string) (*domain.Category, error) {
	const query = `SELECT id, name FROM testing_info.category WHERE name = $1`
	cat := new(domain.Category)

	if err := r.db.QueryRow(ctx, query, name).Scan(cat); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	tasks, err := r.ReadTasksByCategoryID(ctx, cat.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	cat.Tasks = tasks

	return cat, nil
}

func (r *Repository) ReadCategoriesByTestId(ctx context.Context, testID int) ([]domain.Category, error) {
	const query = `SELECT id, name FROM testing_info.category c 
	JOIN testing_info.test_category tc on c.id = tc.category_id
	WHERE tc.test_id = $1`

	rows, err := r.db.Query(ctx, query, testID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	categories, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Category])
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	for i := range categories {
		tasks, err := r.ReadTasksByCategoryID(ctx, categories[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get category: %w", err)
		}
		categories[i].Tasks = tasks
	}
	return categories, nil
}

func (r *Repository) ReadCategories(ctx context.Context) ([]domain.Category, error) {
	const query = "SELECT id, name FROM testing_info.category"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to read categories: %w", err)
	}
	cats, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Category])
	if err != nil {
		return nil, fmt.Errorf("failed to read categories: %w", err)
	}

	return cats, nil
}

// Method to update category in db.
//
// Pass `Category` structure with empty `Tasks` slice.
func (r *Repository) UpdateCategory(ctx context.Context, cat *domain.Category) error {

	if cat.Name != "" {
		query := `UPDATE testing_info.category SET name = $1 WHERE id = $2`
		_, err := r.db.Exec(ctx, query, cat.Name, cat.ID)

		if err != nil {
			return fmt.Errorf("failed to update category: %w", err)
		}
	}
	return nil

}

func (r *Repository) UpdateCategoryTaskRelations(ctx context.Context, catID int, ids []int) error {
	query := `SELECT task_id FROM testing_info.category_task WHERE category_id = $1`
	rows, err := r.db.Query(ctx, query, catID)
	if err != nil {
		return fmt.Errorf("failed to get tasks: %w", err)
	}
	tasksInDB, err := pgx.CollectRows(rows, pgx.RowTo[int])
	if err != nil {
		return fmt.Errorf("failed to get tasks: %w", err)
	}

	tasksToAdd := difference(ids, tasksInDB)
	tasksToDelete := difference(tasksInDB, ids)

	if err := r.AddTasksToCategory(ctx, catID, tasksToAdd); err != nil {
		return fmt.Errorf("failed to add tasks to category: %w", err)
	}

	if err := r.RemoveCategoriesFromTest(ctx, catID, tasksToDelete); err != nil {
		return fmt.Errorf("failed to remove tasks from category: %w", err)
	}

	return nil
}

func (r *Repository) DeleteCategory(ctx context.Context, catID int) error {
	const query = `DELETE FROM testing_info.category WHERE id = $1`

	_, err := r.db.Exec(ctx, query, catID)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}
