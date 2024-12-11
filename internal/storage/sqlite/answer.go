package sqlite

import (
	"fmt"
	"learnDB/internal/domain"

	"github.com/jmoiron/sqlx"
)

type AnswerStorage struct {
	db *sqlx.DB
}

func NewAnswerStorage(db *sqlx.DB) *AnswerStorage {
	return &AnswerStorage{db: db}
}

func (r *AnswerStorage) Insert(ans *domain.Answer) error {
	_, err := r.db.NamedExec("insert into answer(answer_text, is_correct, question_id, query_id) values (:answer_text, :is_correct, :question_id, :query_id);", ans)
	if err != nil {
		return fmt.Errorf("insert answer error: %s", err)
	}
	return nil
}

func (r *AnswerStorage) Get(id int) (*domain.Answer, error) {
	ans := new(domain.Answer)
	if err := r.db.Get(ans, "select answer_id, answer_text, is_correct, question_id, query_id from answer where answer_id = ?;", id); err != nil {
		return nil, fmt.Errorf("select answer error: %s", err)
	}
	return ans, nil
}

func (r *AnswerStorage) GetAll() (answers []domain.Answer, err error) {

	if err = r.db.Select(&answers, "select answer_id, answer_text, is_correct, question_id, query_id from answer;"); err != nil {
		return nil, fmt.Errorf("select answers error: %s", err)
	}
	return answers, nil
}

func (r *AnswerStorage) Update(ans *domain.Answer) error {
	_, err := r.db.NamedExec("update answer set answer_text=:answer_text, is_correct=:is_correct, question_id=:question_id, query_id=:query_id where answer_id=:answer_id", ans)
	if err != nil {
		return fmt.Errorf("update answer error: %s", err)
	}
	return nil
}

func (r *AnswerStorage) Delete(id int) error {
	if _, err := r.db.Exec("delete from answer where answer_id = ?", id); err != nil {
		return fmt.Errorf("delete answer error: %s", err)
	}
	return nil
}
