package domain

import "github.com/pochkachaiki/learndb/internal/domain/answer"

type CheckResult struct {
	Points int
	Error  error
}

// Struct to check it in "checkAnswer" function
type CheckUnit struct {
	Script         string
	CorrectAnswers []answer.CorrectAnswer
	Review         CheckResult
}
