package excel

import answer "learnDB/internal/domain/answer"

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
