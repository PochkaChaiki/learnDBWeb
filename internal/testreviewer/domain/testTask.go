package domain

import "github.com/pochkachaiki/learndb/internal/domain/answer"

type TestTask struct {
	QuestionText   string
	Answer         string
	CorrectAnswers []answer.CorrectAnswer
}

type Task struct {
	Review CheckResult
	TestTask
}
