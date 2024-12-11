package sqlite

import (
	"fmt"
	"learnDB/internal/domain"

	"github.com/jmoiron/sqlx"
)

type QuestionStorage struct {
	db *sqlx.DB
}

func NewQuestionStorage(db *sqlx.DB) *QuestionStorage {
	return &QuestionStorage{db: db}
}

func (r *QuestionStorage) Insert(q *domain.Question) error {
	_, err := r.db.NamedExec("insert into question(question_text, correct_answer, db_sample_id) values (:question_text, :correct_answer, :db_sample_id);", q)
	if err != nil {
		return fmt.Errorf("insert question error: %s", err)
	}
	return nil
}

func (r *QuestionStorage) Get(id int) (*domain.Question, error) {
	q := new(domain.Question)
	if err := r.db.Get(q, "select question_id, question_text, correct_answer, db_sample_id from question where question_id = ?;", id); err != nil {
		return nil, fmt.Errorf("select question error: %s", err)
	}
	return q, nil
}

func (r *QuestionStorage) GetAll() (questions []domain.Question, err error) {

	if err := r.db.Select(&questions, "select question_id, question_text, correct_answer, db_sample_id from question;"); err != nil {
		return nil, fmt.Errorf("select questions error: %s", err)
	}
	return questions, nil
}

func (r *QuestionStorage) Update(q *domain.Question) error {
	_, err := r.db.NamedExec("update question set question_text=:question_text, correct_answer=:correct_answer, db_sample_id=:db_sample_id where question_id=:question_id", q)
	if err != nil {
		return fmt.Errorf("update question error: %s", err)
	}
	return nil
}

func (r *QuestionStorage) Delete(id int) error {
	if _, err := r.db.Exec("delete from question where question_id = ?", id); err != nil {
		return fmt.Errorf("delete question error: %s", err)
	}
	return nil
}
