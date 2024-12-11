package domain

type Answer struct {
	Id         int    `json:"-" db:"answer_id"`
	AnswerText string `json:"answer" db:"answer_text"`
	IsCorrect  bool   `json:"is_correct" db:"is_correct"`
	QuestionId int    `json:"question" db:"question_id"`
	QueryId    int    `json:"query" db:"query_id"`
}
