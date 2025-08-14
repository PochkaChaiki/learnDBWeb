package excel

import "github.com/pochkachaiki/learndb/internal/domain/answer"

type StudentWork struct {
	Name  string
	Group string
	Tasks []Task
}

type Task struct {
	Question       string
	Answer         string
	CorrectAnswers []answer.CorrectAnswer
}
