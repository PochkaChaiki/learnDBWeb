// Package answer provides domain model for answers check
package answer

type Answer struct {
	ID         int    `json:"answer_id,omitempty" db:"answer_id"`
	AnswerText string `json:"answer" db:"answer_text"`
	IsCorrect  bool   `json:"is_correct,omitempty" db:"is_correct"`
	QuestionID int    `json:"question_id" db:"question_id"`
	QueryID    int    `json:"query_id" db:"query_id"`
}

type CorrectAnswer struct {
	Values []string `json:"values"`
	Points int      `json:"points"`
}
