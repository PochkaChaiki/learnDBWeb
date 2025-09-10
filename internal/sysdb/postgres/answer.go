package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pochkachaiki/learndb/internal/domain"
)

func (r *Repository) createXMLAnswer(ctx context.Context, tx pgx.Tx, ans *domain.XMLAnswer) error {
	return nil
}

func (r *Repository) CreateAnswer(answer *domain.XMLAnswer) error {
	const query = `INSERT INTO testing_info.answer 
                  (text, feedback) 
                  VALUES ($1, $2, $3) RETURNING id`

	err := r.db.QueryRowx(query,
		answer.Text, answer.Feedback).Scan(&answer.ID)
	if err != nil {
		return fmt.Errorf("failed to create xml answer: %w", err)
	}

	return nil
}

// Get single answer by id
func (r *Repository) ReadXMLAnswerByID(ctx context.Context, answerID int) (*domain.XMLAnswer, error) {
	const query = `SELECT id, text, feedback 
                  FROM testing_info.answer WHERE id = $1`
	answer := new(domain.XMLAnswer)

	if err := r.db.Get(answer, query, answerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("answer not found")
		}
		return nil, fmt.Errorf("failed to get xml answer: %w", err)
	}

	return answer, nil
}

func (r *Repository) ReadXMLAnswersByTaskID(ctx context.Context, taskID int) ([]domain.XMLAnswer, error) {
	const query = `SELECT a.id, a.text, a.feedback
                  FROM testing_info.answer a
                  WHERE a.task_id = $1`

	var answers []domain.XMLAnswer
	err := r.db.Select(&answers, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sml answers: %w", err)
	}
	return answers, nil
}

func (r *Repository) UpdateXMLAnswer(answer *domain.XMLAnswer) error {
	const query = `UPDATE testing_info.answer 
                  SET text = $1, feedback = $2
                  WHERE id = $3`

	_, err := r.db.Exec(query,
		answer.Text, answer.Feedback, answer.ID)
	if err != nil {
		return fmt.Errorf("failed to update xml answer: %w", err)
	}

	return nil
}

func (r *Repository) DeleteXMLAnswer(answerID int) error {
	const query = `DELETE FROM testing_info.answer WHERE id = $1`

	_, err := r.db.Exec(query, answerID)
	if err != nil {
		return fmt.Errorf("failed to delete xml answer: %w", err)
	}

	return nil
}
