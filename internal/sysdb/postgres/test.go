package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pochkachaiki/learndb/internal/domain"
)

// Test CRUD

// Method to create test in db.
//
// Pass `Test` structure with empty `Categories` slice.
func (r *Repository) CreateTest(ctx context.Context, test *domain.Test) error {
	query := `INSERT INTO testing_info.test (name) VALUES ($1) RETURNING id`

	if err := r.db.QueryRow(ctx, query, test.Name).Scan(&test.ID); err != nil {
		return fmt.Errorf("failed to create test: %w", err)
	}

	return nil
}

func (r *Repository) AddCategoriesToTest(ctx context.Context, testID int, catIDs []int) error {
	query := `INSERT INTO testing_info.test_category (test_id, category_id) VALUES ($1, $2)`

	if err := r.addRelations(ctx, query, testID, catIDs); err != nil {
		return fmt.Errorf("failed to add categories to test: %w", err)
	}
	return nil
}

func (r *Repository) RemoveCategoriesFromTest(ctx context.Context, testID int, catIDs []int) error {
	query := `DELETE FROM testing_info.test_category WHERE test_id = $1 AND category_id = ANY ($2)`

	if _, err := r.db.Exec(ctx, query, testID, catIDs); err != nil {
		return fmt.Errorf("failed to remove categories from test: %w", err)
	}
	return nil
}

func (r *Repository) ReadTestByID(ctx context.Context, id int) (*domain.Test, error) {
	const query = `SELECT id, name FROM testing_info.test WHERE id = $1`
	test := new(domain.Test)

	if err := r.db.QueryRow(ctx, query, id).Scan(test); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("test not found")
		}
		return nil, fmt.Errorf("failed to get test: %w", err)
	}
	// Get categories for test
	cats, err := r.ReadCategoriesByTestId(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get test: %w", err)
	}

	test.Categories = cats

	return test, nil
}

func (r *Repository) ReadTestByName(ctx context.Context, name string) (*domain.Test, error) {
	query := `SELECT id, name FROM testing_info.test WHERE name = $1`
	test := new(domain.Test)

	if err := r.db.QueryRow(ctx, query, name).Scan(test); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("test not found")
		}
		return nil, fmt.Errorf("failed to get test: %w", err)
	}
	// Get categories for test
	cats, err := r.ReadCategoriesByTestId(ctx, test.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test: %w", err)
	}

	test.Categories = cats

	return test, nil
}

func (r *Repository) ReadTests(ctx context.Context) ([]domain.Test, error) {
	const query = "SELECT id, name FROM testing_info.test"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to read tests: %w", err)
	}
	tests, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Test])
	if err != nil {
		return nil, fmt.Errorf("failed to read tests: %w", err)
	}

	return tests, nil
}

// Method to update test in db.
//
// Pass `Test` structure with empty `Categories` slice.
func (r *Repository) UpdateTest(ctx context.Context, test *domain.Test) error {
	if test.Name != "" {
		query := `UPDATE testing_info.test SET name = $1 WHERE id = $2`
		_, err := r.db.Exec(ctx, query, test.Name, test.ID)

		if err != nil {
			return fmt.Errorf("failed to update test: %w", err)
		}
	}

	return nil
}

func (r *Repository) UpdateTestCategoryRelations(ctx context.Context, testID int, ids []int) error {
	query := "SELECT category_id FROM testing_info.test_category WHERE task_id = $1"
	rows, err := r.db.Query(ctx, query, testID)
	if err != nil {
		return fmt.Errorf("failed to read categories ids: %w", err)
	}
	catsInDB, err := pgx.CollectRows(rows, pgx.RowTo[int])
	if err != nil {
		return fmt.Errorf("failed to read categories ids: %w", err)
	}

	catsToAdd := difference(ids, catsInDB)
	catsToDelete := difference(catsInDB, ids)

	if err := r.AddCategoriesToTest(ctx, testID, catsToAdd); err != nil {
		return fmt.Errorf("failed to add categories to test: %w", err)
	}

	if err := r.RemoveCategoriesFromTest(ctx, testID, catsToDelete); err != nil {
		return fmt.Errorf("failed to remove categories from test: %w", err)
	}
	return nil
}

func (r *Repository) DeleteTest(ctx context.Context, testID int) error {
	const query = `DELETE FROM testing_info.test WHERE id = $1`

	_, err := r.db.Exec(ctx, query, testID)
	if err != nil {
		return fmt.Errorf("failed to delete test: %w", err)
	}

	return nil
}
