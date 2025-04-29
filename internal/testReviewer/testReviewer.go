package testReviewer

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"learnDB/internal/domain/answer"
	"learnDB/internal/testReviewer/domain"
)

type TestReviewer struct {
	repo map[string]DBRepository
}

func New(repo map[string]DBRepository) *TestReviewer {
	return &TestReviewer{repo: repo}
}

// Function to check if students answer is correct
func (tr *TestReviewer) checkAnswer(db string, sql string, correctAnswers []answer.CorrectAnswer) *domain.CheckResult {

	var checkResult domain.CheckResult
	qRes, err := tr.repo[db].RunSelect(sql, 1)

	sort.Slice(correctAnswers, func(i, j int) bool {
		return correctAnswers[j].Points > correctAnswers[i].Points
	})

	errors.Join(checkResult.Error, err)

	data := qRes.Data[0]

	// Going down the slice to find most valuable (grades with most points)
	// correct answer that matches student's answer
	for _, corrAns := range correctAnswers {

		if len(corrAns.Values) != len(data) {
			continue
		}

		ansIsCorrect := true
		// If every part of student's answer is present at the correct answer than the first one is correct
		for _, ans := range data {
			ansToCheck, ok := ans.(string)
			if !ok {
				errors.Join(checkResult.Error, fmt.Errorf("type assertion error: %v to string", data[0]))
			}
			partIsCorrect := false
			for _, corrAnsPart := range corrAns.Values {
				partIsCorrect = partIsCorrect || strings.EqualFold(ansToCheck, corrAnsPart)
				if partIsCorrect {
					ansIsCorrect = ansIsCorrect && partIsCorrect
					break
				}
			}
		}

		if ansIsCorrect {
			checkResult.Points = corrAns.Points
			break
		}
	}

	return &checkResult
}

// Check the work of single student
func (tr *TestReviewer) GradeTheWork(db string, tt []domain.TestTask) domain.WorkReview {
	wr := make(domain.WorkReview, 0, len(tt))

	ch := make(chan *domain.CheckResult)
	var wg sync.WaitGroup
	for _, task := range tt {
		wg.Add(1)
		go func(task *domain.TestTask, ch chan *domain.CheckResult) {
			defer wg.Done()
			// Get script from cell
			sql := RetrieveScript(task.Answer)

			// Run Script
			res := tr.checkAnswer(db, sql, task.CorrectAnswers)
			ch <- res
		}(&task, ch)
		// Save information about student and theirs points for answers

		go func() {
			wg.Wait()
			close(ch)
		}()

		for res := range ch {
			wr = append(wr, *res)
		}
	}

	return wr
}

func (tr *TestReviewer) CheckTest(sw []domain.StudentWork) domain.ReviewedWorks {
	res := make(domain.ReviewedWorks, 0, len(sw))
	ch := make(chan *domain.ReviewedStudentWork)

	var wg sync.WaitGroup
	for i := range sw {
		wg.Add(1)
		go func(work *domain.StudentWork, ch chan *domain.ReviewedStudentWork) {
			defer wg.Done()
			wr := tr.GradeTheWork(work.DB, work.Tasks)
			totalGrade := 0
			for _, cr := range wr {
				totalGrade += cr.Points
			}
			ch <- &domain.ReviewedStudentWork{
				TotalGrade:  totalGrade,
				StudentWork: *work,
				WorkReview:  wr,
			}
		}(&sw[i], ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for rsw := range ch {
		res = append(res, *rsw)
	}

	return res
}
