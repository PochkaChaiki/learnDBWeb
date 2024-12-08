package question

import (
	"fmt"
	"learnDB/internal/storage"
)

type Question struct {
	Id            int    `json:"-" db:"question_id"`
	QuestionText  string `json:"question_text" db:"question_text"`
	CorrectAnswer string `json:"correct_answer" db:"correct_answer"`
	DBSampleId    int    `json:"-" db:"db_sample_id"`
}

func (q *Question) Insert(st *storage.Storage) error {
	_, err := st.DB.NamedExec("insert into question(question_text, correct_answer, db_sample_id) values (:question_text, :correct_answer, :db_sample_id);", q)
	if err != nil {
		return fmt.Errorf("insert question error: %s", err)
	}
	return nil
}

func (q *Question) Get(st *storage.Storage, id int) error {

	if err := st.DB.Get(q, "select question_id, question_text, correct_answer, db_sample_id from question where question_id = ?;", id); err != nil {
		return fmt.Errorf("select question error: %s", err)
	}
	return nil
}

func GetAll(st *storage.Storage) ([]Question, error) {

	questions := []Question{}

	if err := st.DB.Select(&questions, "select question_id, question_text, correct_answer, db_sample_id from question;"); err != nil {
		return nil, fmt.Errorf("select questions error: %s", err)
	}
	return questions, nil
}

func (q *Question) Update(st *storage.Storage) error {
	_, err := st.DB.NamedExec("update question set question_text=:question_text, correct_answer=:correct_answer, db_sample_id=:db_sample_id where question_id=:question_id", q)
	if err != nil {
		return fmt.Errorf("update question error: %s", err)
	}
	return nil
}

func (q *Question) Delete(st *storage.Storage) error {
	if _, err := st.DB.Exec("delete from question where question_id = ?", q.Id); err != nil {
		return fmt.Errorf("delete question error: %s", err)
	}
	return nil
}
