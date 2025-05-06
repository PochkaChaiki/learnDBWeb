package testReviewer

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

	"learnDB/internal/domain/answer"
	"learnDB/internal/testReviewer/domain"
)

type TestReviewer struct {
	repo       map[string]DBRepository
	schemaName string
}

func New(repo map[string]DBRepository, schemaName string) *TestReviewer {
	return &TestReviewer{repo: repo, schemaName: schemaName}
}

// Function to check if students answer is correct
func (tr *TestReviewer) checkAnswer(db string, sql string, correctAnswers []answer.CorrectAnswer) *domain.CheckResult {
	if sql == "" {
		return &domain.CheckResult{
			Points: 0,
			Error:  errors.New("empty answer"),
		}
	}

	var checkResult domain.CheckResult
	qRes, err := tr.repo[db].RunSelect(sql, tr.schemaName, 1)

	if err != nil {
		return &domain.CheckResult{
			Points: 0,
			Error:  err,
		}
	}

	sort.Slice(correctAnswers, func(i, j int) bool {
		return correctAnswers[j].Points > correctAnswers[i].Points
	})

	if len(qRes.Data) == 0 {
		return &domain.CheckResult{
			Points: 0,
			Error:  errors.New("null returned"),
		}
	}

	data := qRes.Data[0]

	anyAnswer := false // this variable is used to print "incorrect answer"

	// Going down the slice to find most valuable (grades with most points)
	// correct answer that matches student's answer
	for _, corrAns := range correctAnswers {

		if len(corrAns.Values) != len(data) {
			continue
		}

		ansIsCorrect := true
		// If every part of student's answer is present at the correct answer than the first one is correct
		for _, ans := range data {
			ansToCheck := fmt.Sprint(ans)
			partIsCorrect := false
			for _, corrAnsPart := range corrAns.Values {
				partIsCorrect = partIsCorrect || strings.EqualFold(ansToCheck, corrAnsPart)
				if partIsCorrect {
					break
				}
			}
			ansIsCorrect = ansIsCorrect && partIsCorrect
		}

		if ansIsCorrect {
			checkResult.Points = corrAns.Points
			anyAnswer = true
			break
		}
	}

	if !anyAnswer {
		errors.Join(checkResult.Error, errors.New("incorrect answer"))
	}
	return &checkResult
}

// Check the work of single student
func (tr *TestReviewer) GradeTheWork(db string, tt []domain.TestTask) domain.WorkReview {
	wr := make(domain.WorkReview, 0, len(tt))

	// ch := make(chan *domain.CheckResult)
	// var wg sync.WaitGroup
	for _, task := range tt {
		// wg.Add(1)
		// go func(task *domain.TestTask, ch chan *domain.CheckResult) {
		// defer wg.Done()
		// Get script from cell
		sql := RetrieveScript(task.Answer)

		// Run Script
		res := tr.checkAnswer(db, sql, task.CorrectAnswers)

		// ch <- res
		// }(&task, ch)
		// Save information about student and theirs points for answers

		// go func() {
		// 	wg.Wait()
		// 	close(ch)
		// }()

		// for res := range ch {
		wr = append(wr, *res)
		// }
	}

	return wr
}

func (tr *TestReviewer) CheckTest(sw []domain.StudentWork) domain.ReviewedWorks {
	res := make(domain.ReviewedWorks, 0, len(sw))
	// ch := make(chan *domain.ReviewedStudentWork)

	// var wg sync.WaitGroup
	for i := range sw {
		// wg.Add(1)
		// go func(work *domain.StudentWork, ch chan *domain.ReviewedStudentWork) {
		// defer wg.Done()
		log.Printf(sw[i].Name)
		work := &sw[i]
		wr := tr.GradeTheWork(work.DB, work.Tasks)
		totalGrade := 0
		for _, cr := range wr {
			totalGrade += cr.Points
		}
		// ch <- &domain.ReviewedStudentWork{
		res = append(res, domain.ReviewedStudentWork{
			TotalGrade:  totalGrade,
			StudentWork: *work,
			WorkReview:  wr,
		})
		// }(&sw[i], ch)
	}

	// go func() {
	// 	wg.Wait()
	// 	close(ch)
	// }()

	// for rsw := range ch {
	// 	res = append(res, *rsw)
	// }

	return res
}
