package testreviewer

import (
	"errors"
	"sort"
	"sync"

	"github.com/pochkachaiki/learndb/internal/domain/answer"
	"github.com/pochkachaiki/learndb/internal/testreview/domain"

	mainDomain "github.com/pochkachaiki/learndb/internal/domain"
)

type DBManager interface {
	GetRepository(string) (mainDomain.DBRepository, bool)
}

type TestReviewer struct {
	manager    DBManager
	schemaName string
}

func New(manager DBManager, schemaName string) *TestReviewer {
	return &TestReviewer{
		manager:    manager,
		schemaName: schemaName,
	}
}

// Does this field should be protected somehow???
func (tr *TestReviewer) ChangeSchema(schemaName string) {
	tr.schemaName = schemaName
}

// Function to check if students answer is correct
func (tr *TestReviewer) checkAnswer(db string, sql string, correctAnswers []answer.CorrectAnswer, strictMode bool) domain.CheckResult {
	if sql == "" {
		return domain.CheckResult{
			Points: 0,
			Comment:  mainDomain.ErrEmptyAnswer,
		}
	}

	var checkResult domain.CheckResult

	// Omit "ok" for now
	repo, _ := tr.manager.GetRepository(db)
	qRes, err := repo.RunSelect(sql, tr.schemaName, 1)

	if err != nil {
		return domain.CheckResult{
			Points: 0,
			Comment:  err,
		}
	}

	sort.Slice(correctAnswers, func(i, j int) bool {
		return correctAnswers[j].Points > correctAnswers[i].Points
	})

	if len(qRes.Data) == 0 {
		return domain.CheckResult{
			Points: 0,
			Comment:  errors.Join(mainDomain.ErrIncorrectAnswer, errors.New("null returned")),
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
			checkResult.Points, anyAnswer = answer.StrictCheck(data, corrAns)
		} else {
			checkResult.Points, anyAnswer = answer.LightCheck(data, corrAns)
		}
		if anyAnswer {
			break
		}
	}

	if !anyAnswer {
		checkResult.Comment = errors.Join(mainDomain.ErrIncorrectAnswer, checkResult.Comment)
	}
	return checkResult
}

// Check a work of a single student
func (tr *TestReviewer) GradeTheWork(work *domain.Work) {

	for _, unit := range work.Units {

		// Set strict mode to false value to let the application
		// to find the correct answer among the student's answer columns
		res := tr.checkAnswer(work.DB, unit.Script, unit.CorrectAnswers, false)

		// Save information about student's points for answers
		unit.Review.Comment = errors.Join(unit.Review.Comment, res.Comment)
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

type WorkerPool struct {
}

func (wp *WorkerPool) Run(tasks <-chan func(domain.Work), size int) error{
	out := make(chan )
	for range size {
		go func() {
			
		}
	}
}
