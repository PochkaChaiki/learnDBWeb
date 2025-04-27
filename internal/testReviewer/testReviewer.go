package testReviewer

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"learnDB/internal/domain/answer"
	"learnDB/internal/testReviewer/domain"
)

func RetrieveScript(ans string) string {
	start := "select"
	end := ";"
	lowerCasedAns := strings.ToLower(ans)

	startIndex := strings.Index(lowerCasedAns, start)
	if startIndex == -1 {
		return ""
	}

	if startIndex != 0 {
		switch symbolBefore := []rune(lowerCasedAns)[startIndex-1]; symbolBefore {
		case '"':
			end = "\""
		case '(':
			end = ")"
		case '[':
			end = "]"
		case '\'':
			end = "'"
		default:
		}
	}

	if endIndex := strings.Index(lowerCasedAns[startIndex:], end); endIndex != -1 {
		return lowerCasedAns[startIndex : endIndex+1]
	}

	return lowerCasedAns[startIndex:]
}

type TestReviewer struct {
	repo DBRepository
}

func New(repo DBRepository) *TestReviewer {
	return &TestReviewer{repo: repo}
}

// Function to check if students answer is correct
func (tr *TestReviewer) checkAnswer(sql string, correctAnswers []answer.CorrectAnswer) *domain.CheckResult {

	var checkResult domain.CheckResult
	qRes, err := tr.repo.RunSelect(sql, 1)

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
func (tr *TestReviewer) GradeTheWork(tt []domain.TestTask) domain.WorkReview {
	wr := make(domain.WorkReview, 0, len(tt))

	for _, task := range tt {
		// Get script from cell
		sql := RetrieveScript(task.Answer)

		// Run Script
		res := tr.checkAnswer(sql, task.CorrectAnswers)

		// Save information about student and theirs points for answers
		wr = append(wr, *res)
	}

	return wr
}

func (tr *TestReviewer) CheckTest(sw []domain.StudentWork) domain.ReviewedWorks {
	res := make(domain.ReviewedWorks, 0, len(sw))
	for i := range sw {
		work := &sw[i]
		wr := tr.GradeTheWork(work.Tasks)
		totalGrade := 0
		for _, cr := range wr {
			totalGrade += cr.Points
		}
		res = append(res, domain.ReviewedStudentWork{
			TotalGrade:  totalGrade,
			StudentWork: *work,
			WorkReview:  wr,
		})
	}
	return res
}
