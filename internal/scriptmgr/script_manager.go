package scriptmgr

import (
	"encoding/json"
	"errors"
	"math/rand"
	"regexp"

	"github.com/pochkachaiki/learndb/internal/domain"
	"github.com/pochkachaiki/learndb/internal/domain/answer"
)

const regexString = "(?i)[^;]+"

type SystemDBRepository interface {
	GetTaskByName(string) (*domain.Task, error)
	GetTasksByCategoryID(int) ([]domain.Task, error)
	GetDBMSList() ([]domain.DBMS, error)
	GetSchemasOfDBMS(string) ([]string, error)
}

type DBManager interface {
	GetRepository(string) (domain.DBRepository, bool)
}

type ScriptManager struct {
	manager   DBManager
	sysDBrepo SystemDBRepository
	limit     int
	scriptRe  *regexp.Regexp
}

func New(manager DBManager, sysDBRepo SystemDBRepository, limit int) *ScriptManager {
	scriptRe, err := regexp.Compile(regexString)
	if err != nil {
		return nil
	}

	return &ScriptManager{
		manager:   manager,
		limit:     limit,
		sysDBrepo: sysDBRepo,
		scriptRe:  scriptRe,
	}
}

func (s *ScriptManager) GetOpts() *ScriptOptions {
	return &ScriptOptions{}
}

func (s *ScriptManager) GetRandomTask() (*domain.Task, error) {

	tasks, err := s.sysDBrepo.GetTasksByCategoryID(1)
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, errors.New("no tasks returned")
	}
	ind := rand.Intn(len(tasks)) + 1

	task := tasks[ind]
	return &task, nil
}

func (s *ScriptManager) GetDBMSList() ([]domain.DBMS, error) {
	return s.sysDBrepo.GetDBMSList()
}

func (s *ScriptManager) GetSchemasOfDBMS(dbms string) ([]string, error) { // I NEED TO KEEP SCHEMAS AT DB
	return s.sysDBrepo.GetSchemasOfDBMS(dbms)
}

func (s *ScriptManager) RunScript(opts *ScriptOptions) []ScriptRunResult {
	// Omit "ok" for now
	repo, _ := s.manager.GetRepository(opts.DBMS)

	sqls := s.scriptRe.FindAllString(opts.Sql, -1)
	results := make([]ScriptRunResult, len(sqls))

	for _, sql := range sqls {
		qr, err := repo.RunSelect(sql, opts.SchemaName, s.limit)

		var res ScriptRunResult

		if err != nil {
			res.Error = err.Error()
		} else {
			res.Columns = qr.Columns
			res.Data = qr.Data
		}
		results = append(results, res)
	}

	return results
}

func (s *ScriptManager) ValidateScript(opts *ScriptOptions) (bool, string, error) {

	results := s.RunScript(opts)

	task, err := s.sysDBrepo.GetTaskByName(opts.TaskName)
	if err != nil {
		return false, "", err
	}

	var correctAnswers []answer.CorrectAnswer

	if err := json.Unmarshal(task.CorrectAnswerJSONB, &correctAnswers); err != nil {
		return false, "", err
	}

	if len(correctAnswers) == 0 {
		return false, "", errors.New("correctAnswers is nil")
	}

	check, comment := s.validate(results[0].Data, correctAnswers[0])
	return check, comment, nil
}

func (s *ScriptManager) validate(data [][]any, corrAns answer.CorrectAnswer) (bool, string) {

	if len(data) > 1 {
		return false, "Script returned more rows that was expected"
	}

	// First we try to validate on strictMode
	_, passed := answer.StrictCheck(data[0], corrAns)

	if passed {
		return true, "Excellent"
	}

	// Second we try to validate on strictMode
	_, passed = answer.LightCheck(data[0], corrAns)

	if passed {
		return true, "Ok, but returned more columns that was expected"
	}

	return false, "Script is incorrect"
}
