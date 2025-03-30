package excelReader

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

func RetrieveScript(ans string) string {
	scriptStart := "select"
	scriptEnd := ";"
	lowerCasedAns := strings.ToLower(ans)

	scriptStartIndex := strings.Index(lowerCasedAns, scriptStart)
	if scriptStartIndex == -1 {
		return ""
	}

	if scriptStartIndex != 0 {
		switch symbolBefore := []rune(lowerCasedAns)[scriptStartIndex-1]; symbolBefore {
		case '"':
			scriptEnd = "\""
		case '(':
			scriptEnd = "("
		case '[':
			scriptEnd = "["
		case '\'':
			scriptEnd = "'"
		default:
		}
	}

	scriptEndIndex := strings.Index(lowerCasedAns[scriptStartIndex:], scriptEnd)

	script := lowerCasedAns[scriptStartIndex:scriptEndIndex]
	return script
}

type Repository interface {
	RunScript(sql string, limit int) (*QueryResult, error)
}

type TestChecker struct {
	repo Repository
}

func New(repo Repository) *TestChecker {
	return &TestChecker{repo: repo}
}

type QueryResult struct {
	Columns []string
	Data    [][]any
}

type CorrectAnswer struct {
	Values []string
	Points int
}

// type CorrectAnswer []CorrectAnswer

type CheckResult struct {
	Points int
	Error  error
}

// Function to check if students answer is correct
func (tc *TestChecker) checkAnswer(sql string, correctAnswers []CorrectAnswer) *CheckResult {

	var checkResult CheckResult
	qRes, err := tc.repo.RunScript(sql, 1)

	sort.Slice(correctAnswers, func(i, j int) bool {
		return correctAnswers[j].Points > correctAnswers[i].Points
	})

	errors.Join(checkResult.Error, err)

	data := qRes.Data[0]

	// Going down the list to find most valuable (grades with most points) correct answer
	for _, corrAns := range correctAnswers {

		if len(corrAns.Values) != len(data) {
			continue
		}

		ansIsCorrect := true
		// If every student's part of answer is present at correct answer than student's answer is correct
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

type TestTask struct {
	QuestionText   string
	CorrectAnswers []CorrectAnswer
	Answer         string
}

type StudentWork struct {
	Name  string
	Group string
	Tasks []TestTask
}

type WorkReview []CheckResult

// Check the work of single student
func (tc *TestChecker) GradeTheWork(tt []TestTask) WorkReview {
	// So main questions are:
	//   1) How to extract script from excel
	//   2) How to save check results? : 1] in excel
	wr := make(WorkReview, 0, len(tt))

	for _, task := range tt {
		// Get script from cell
		sql := RetrieveScript(task.Answer)

		// Run Script
		res := tc.checkAnswer(sql, task.CorrectAnswers)

		// Save information about student and theirs points for answers
		wr = append(wr, *res)
	}

	return wr
}

type ReviewedStudentWork struct {
	TotalGrade int
	StudentWork
	WorkReview WorkReview
}

type ReviewedWorks []ReviewedStudentWork

func (tc *TestChecker) CheckTest(sw []*StudentWork) ReviewedWorks {
	res := make(ReviewedWorks, 0, len(sw))
	for _, work := range sw {
		wr := tc.GradeTheWork(work.Tasks)
		totalGrade := 0
		for _, cr := range wr {
			totalGrade += cr.Points
		}
		res = append(res, ReviewedStudentWork{
			TotalGrade:  totalGrade,
			StudentWork: *work,
			WorkReview:  wr,
		})
	}
	return res
}
