package scriptManager

import (
	"encoding/json"
	"errors"
	"learnDB/internal/domain"
	"learnDB/internal/domain/answer"
	"math/rand"
)

type SystemDBRepository interface {
	GetTaskByName(string) (*domain.Task, error)
	GetTasksByCategoryID(int) ([]domain.Task, error)
}

type DBManager interface {
	GetRepository(string) (domain.DBRepository, bool)
}

type ScriptManager struct {
	manager   DBManager
	sysDBrepo SystemDBRepository
	limit     int
}

func New(manager DBManager, sysDBRepo SystemDBRepository, limit int) *ScriptManager {
	return &ScriptManager{
		manager:   manager,
		limit:     limit,
		sysDBrepo: sysDBRepo,
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

func (s *ScriptManager) RunScript(opts *ScriptOptions) *ScriptRunResult {
	res := new(ScriptRunResult)

	// Omit "ok" for now
	repo, _ := s.manager.GetRepository(opts.DBMS)
	qr, err := repo.RunSelect(opts.Sql, opts.SchemaName, s.limit)

	if err != nil {
		res.Error = err.Error()
	} else {
		res.Columns = qr.Columns
		res.Data = qr.Data
	}

	return res
}

func (s *ScriptManager) ValidateScript(opts *ScriptOptions) (bool, string, error) {

	// Omit "ok" for now
	repo, _ := s.manager.GetRepository(opts.DBMS)
	qr, err := repo.RunSelect(opts.Sql, opts.SchemaName, s.limit)
	if err != nil {
		return false, "", err
	}

	task, err := s.sysDBrepo.GetTaskByName(opts.TaskName)
	if err != nil {
		return false, "", err
	}

	var correctAnswers []answer.CorrectAnswer

	if err := json.Unmarshal([]byte(*task.CorrectAnswerJSONB), &correctAnswers); err != nil {
		return false, "", err
	}

	if len(correctAnswers) == 0 {
		return false, "", errors.New("correctAnswers is nil")
	}

	check, comment := s.validate(qr.Data, correctAnswers[0])
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
