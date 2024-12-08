package answer

import (
	"fmt"
	"learnDB/internal/storage"
)

type Answer struct {
	Id         int    `json:"-" db:"answer_id"`
	AnswerText string `json:"answer" db:"answer_text"`
	IsCorrect  bool   `json:"is_correct" db:"is_correct"`
	QuestionId int    `json:"question" db:"question_id"`
	QueryId    int    `json:"query" db:"query_id"`
}

func (ans *Answer) Insert(st *storage.Storage) error {
	_, err := st.DB.NamedExec("insert into answer(answer_text, is_correct, question_id, query_id) values (:answer_text, :is_correct, :question_id, :query_id);", ans)
	if err != nil {
		return fmt.Errorf("insert answer error: %s", err)
	}
	return nil
}

func (ans *Answer) Get(st *storage.Storage, id int) error {

	if err := st.DB.Get(ans, "select answer_id, answer_text, is_correct, question_id, query_id from answer where answer_id = ?;", id); err != nil {
		return fmt.Errorf("select answer error: %s", err)
	}
	return nil
}

func GetAll(st *storage.Storage) ([]Answer, error) {

	answers := []Answer{}

	if err := st.DB.Select(&answers, "select answer_id, answer_text, is_correct, question_id, query_id from answer;"); err != nil {
		return nil, fmt.Errorf("select answers error: %s", err)
	}
	return answers, nil
}

func (ans *Answer) Update(st *storage.Storage) error {
	_, err := st.DB.NamedExec("update answer set answer_text=:answer_text, is_correct=:is_correct, question_id=:question_id, query_id=:query_id where answer_id=:answer_id", ans)
	if err != nil {
		return fmt.Errorf("update answer error: %s", err)
	}
	return nil
}

func (ans *Answer) Delete(st *storage.Storage) error {
	if _, err := st.DB.Exec("delete from answer where answer_id = ?", ans.Id); err != nil {
		return fmt.Errorf("delete answer error: %s", err)
	}
	return nil
}
