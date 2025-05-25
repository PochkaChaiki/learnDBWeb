package testReviewer

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"

	mainDomain "learnDB/internal/domain"
	"learnDB/internal/domain/answer"
	"learnDB/internal/testReviewer/domain"
)

type TestReviewer struct {
	repo       map[string]DBRepository
	schemaName string
}

func New(repo map[string]DBRepository, schemaName string) *TestReviewer {
	return &TestReviewer{
		repo:       repo,
		schemaName: schemaName,
	}
}

// Function to check if students answer is correct
func (tr *TestReviewer) checkAnswer(db string, sql string, correctAnswers []answer.CorrectAnswer, strictMode bool) *domain.CheckResult {
	if sql == "" {
		return &domain.CheckResult{
			Points: 0,
			Error:  mainDomain.ErrEmptyAnswer,
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
			Error:  errors.Join(mainDomain.ErrIncorrectAnswer, errors.New("null returned")),
		}
	}

	data := qRes.Data[0]

	anyAnswer := false // this variable is used to print "incorrect answer"

	//
	// Going down the slice to find most valuable (grades with most points)
	// correct answer that matches student's answer
	//
	for _, corrAns := range correctAnswers {

		if strictMode {
			checkResult.Points, anyAnswer = strictCheck(data, corrAns)
		} else {
			checkResult.Points, anyAnswer = lightCheck(data, corrAns)
		}
		if anyAnswer {
			break
		}
	}

	if !anyAnswer {
		checkResult.Error = errors.Join(mainDomain.ErrIncorrectAnswer, checkResult.Error)
	}
	return &checkResult
}

func strictCheck(data []any, corrAns answer.CorrectAnswer) (int, bool) {
	if len(corrAns.Values) != len(data) {
		return 0, false
	}

	ansIsCorrect := true
	// If every part of student's answer is present at the correct answer than the first one is correct
	for _, ans := range data {
		var ansToCheck string
		switch ans := ans.(type) {
		case string:
			ansToCheck = ans
		case []uint8:
			ansToCheck = string(ans)
		default:
			ansToCheck = fmt.Sprint(ans)
		}

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
		return corrAns.Points, true
	}
	return 0, false
}

func lightCheck(data []any, corrAns answer.CorrectAnswer) (int, bool) {
	coef := float64(len(corrAns.Values)) / float64(len(data))
	if coef > 1 {
		return 0, false
	}

	ansIsCorrect := true
	for _, corrAnsPart := range corrAns.Values {
		partIsCorrect := false
		for _, ans := range data {
			var ansToCheck string
			switch ans := ans.(type) {
			case string:
				ansToCheck = ans
			case []uint8:
				ansToCheck = string(ans)
			default:
				ansToCheck = fmt.Sprint(ans)
			}
			partIsCorrect = partIsCorrect || strings.EqualFold(ansToCheck, corrAnsPart)
			if partIsCorrect {
				break
			}
		}
		ansIsCorrect = ansIsCorrect && partIsCorrect
	}
	if ansIsCorrect {
		return int(math.Round(float64(corrAns.Points) * coef)), true
	}
	return 0, false
}

// Check the work of single student
func (tr *TestReviewer) GradeTheWork(work *domain.Work) {

	for _, unit := range work.Units {

		// Set strict mode to false value to let the application
		// to find the correct answer among the student's answer columns
		res := tr.checkAnswer(work.DB, unit.Script, unit.CorrectAnswers, false)

		// Save information about student's points for answers
		unit.Review.Error = errors.Join(unit.Review.Error, res.Error)
		unit.Review.Points += res.Points
		if unit.Review.Points < 0 {
			unit.Review.Points = 0
		}

	}
}

func (tr *TestReviewer) CheckTest(sw domain.Works) {
	var wg sync.WaitGroup
	for i := range sw {
		wg.Add(1)
		go func(work *domain.Work) {
			defer wg.Done()

			tr.GradeTheWork(work)

			totalGrade := 0
			for _, t := range work.Units {
				totalGrade += t.Review.Points
			}

			work.TotalGrade += totalGrade

		}(sw[i])
	}
	wg.Wait()

}

// // Function to debug
// func (tr *TestReviewer) CheckTest(sw domain.Works) {
// 	for _, work := range sw {
// 		tr.GradeTheWork(work)

// 		totalGrade := 0
// 		for _, t := range work.Units {
// 			totalGrade += t.Review.Points
// 		}

// 		work.TotalGrade += totalGrade

// 	}
// }
