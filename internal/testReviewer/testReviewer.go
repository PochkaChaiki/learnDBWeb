package testReviewer

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	mainDomain "learnDB/internal/domain"
	"learnDB/internal/domain/answer"
	"learnDB/internal/testReviewer/domain"
)

type TestReviewer struct {
	repo               map[string]DBRepository
	schemaName         string
	availableConnToken chan bool
}

func New(repo map[string]DBRepository, schemaName string) *TestReviewer {
	return &TestReviewer{
		repo:               repo,
		schemaName:         schemaName,
		availableConnToken: make(chan bool, 100),
	}
}

// Function to check if students answer is correct
func (tr *TestReviewer) checkAnswer(db string, sql string, correctAnswers []answer.CorrectAnswer) *domain.CheckResult {
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

		if len(corrAns.Values) != len(data) {
			continue
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
			checkResult.Points = corrAns.Points
			anyAnswer = true
			break
		}
	}

	if !anyAnswer {
		checkResult.Error = errors.Join(mainDomain.ErrIncorrectAnswer, checkResult.Error)
	}
	return &checkResult
}

// Check the work of single student
func (tr *TestReviewer) GradeTheWork(work *domain.StudentWork) {
	// reviews := make([]domain.TaskReview, 0, len(tt))

	for i := range work.Tasks {
		task := &work.Tasks[i]
		// Get script from cell
		sql := RetrieveScripts(task.Answer)
		if len(sql) == 0 {
			task.Review = domain.CheckResult{
				Points: 0,
				Error:  mainDomain.ErrEmptyAnswer,
			}
			// reviews = append(reviews, domain.TaskReview{
			// 	Task: review,
			// 	Review: domain.CheckResult{
			// 		Points: 0,
			// 		Error:  mainDomain.ErrEmptyAnswer,
			// 	},
			// })
			continue
		}

		// Run Script
		tr.availableConnToken <- true
		res := tr.checkAnswer(work.DB, sql[len(sql)-1], task.CorrectAnswers)
		<-tr.availableConnToken

		if len(sql) > 1 && res.Points > 0 {
			res.Points--
		}
		// Save information about student and theirs points for answers
		task.Review = *res

		// reviews = append(reviews, domain.TaskReview{
		// 	Task:   task,
		// 	Review: *res,
		// })

	}
	// return reviews
}

func (tr *TestReviewer) CheckTest(sw domain.Works) {
	// res := make(domain.ReviewedWorks, 0, len(sw))
	// reviewedWorks := make(chan *domain.WorkReview)

	var wg sync.WaitGroup
	for i := range sw {
		wg.Add(1)
		go func(work *domain.StudentWork) {
			defer wg.Done()

			// tr.availableConnToken <- true
			tr.GradeTheWork(work)
			// <-tr.availableConnToken

			totalGrade := 0
			for _, t := range work.Tasks {
				totalGrade += t.Review.Points
			}

			work.TotalGrade += totalGrade

			// reviewedWorks <- &domain.WorkReview{
			// 	Name:       work.Name,
			// 	Group:      work.Group,
			// 	DB:         work.DB,
			// 	TotalGrade: totalGrade,
			// 	Tasks:      reviews,
			// }
		}(sw[i])
	}

	// go func() {
	wg.Wait()
	// 	close(reviewedWorks)
	// }()

	// for rw := range reviewedWorks {
	// 	res = append(res, *rw)
	// }

	// return res
}

// // Function to debug
// func (tr *TestReviewer) CheckTest(sw []domain.StudentWork) domain.ReviewedWorks {
// 	res := make(domain.ReviewedWorks, 0, len(sw))

// 	for _, work := range sw {

// 		reviews := tr.GradeTheWork(work.DB, work.Tasks)
// 		totalGrade := 0
// 		for _, cr := range reviews {
// 			totalGrade += cr.Review.Points
// 		}

// 		log.Println(work.Name)

// 		res = append(res, domain.WorkReview{
// 			Name:       work.Name,
// 			Group:      work.Group,
// 			DB:         work.DB,
// 			TotalGrade: totalGrade,
// 			Tasks:      reviews,
// 		})
// 	}
// 	return res
// }
