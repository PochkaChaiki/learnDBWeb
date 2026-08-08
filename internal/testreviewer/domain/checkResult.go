package domain

import (
	"github.com/pochkachaiki/learndb/internal/domain/answer"
	"github.com/pochkachaiki/learndb/internal/domain/comments"
)

type CheckResult struct {
	Points  int
	Comment comments.Comment
}

// Struct to check it in "checkAnswer" function
type CheckUnit struct {
	Script         string
	CorrectAnswers []answer.CorrectAnswer
	Review         CheckResult
}
