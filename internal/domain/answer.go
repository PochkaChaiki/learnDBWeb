package domain

import (
	"math/rand"
	"time"
)

type Answer struct {
	Id         int    `json:"answer_id,omitempty" db:"answer_id"`
	AnswerText string `json:"answer" db:"answer_text"`
	IsCorrect  bool   `json:"is_correct,omitempty" db:"is_correct"`
	QuestionId int    `json:"question_id" db:"question_id"`
	QueryId    int    `json:"query_id" db:"query_id"`
}

func (ans *Answer) CheckAnswer() bool {
	rand := rand.New(rand.NewSource(time.Now().Unix()))
	return rand.Int()%2 != 0
}
