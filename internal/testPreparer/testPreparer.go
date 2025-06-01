package testPreparer

import (
	"encoding/xml"
	"io"

	"learnDB/internal/systemDBRepository/postgres"
)

type TestPreparer struct {
	repo postgres.Repository
}

func (p *TestPreparer) ReadXML(r io.Reader) (*Quiz, error) {
	q := new(Quiz)
	if err := xml.NewDecoder(r).Decode(q); err != nil {
		return nil, err
	}
	return q, nil
}

func (p *TestPreparer) WriteXML(w io.Writer, q *Quiz) error {
	err := xml.NewEncoder(w).Encode(q)
	return err
}

// func (p *TestPreparer) GetTest(name string) (*Quiz, error) {
// 	test, err := p.repo.GetTest(name)
// 	if err != nil {
// 		return nil, err
// 	}
// 	q := makeQuizPtr(test.Name)

// 	for _, task := range test.Tasks {
// 		if task.Name
// 		question := makeQuestion(task.Name)
// 		for _, ans := range task.Answers {
// 			var ansToPass answer
// 			if ans.CorrectAnswerJSONB != nil {
// 				ansToPass = makeAnswer(*ans.CorrectAnswerJSONB, *ans.Feedback)

// 			} else {
// 				ansToPass = makeAnswer(ans.Text, *ans.Feedback)
// 			}
// 			question.Answer = append(question.Answer, ansToPass)
// 		}
// 		question.DefaultGrade = float64(task.Points)
// 		q.Question = append()
// 	}

// }
