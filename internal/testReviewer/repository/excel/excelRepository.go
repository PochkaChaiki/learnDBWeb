package excel

import (
	"encoding/json"
	"fmt"
	domainAnswer "learnDB/internal/domain/answer"
	"learnDB/internal/testReviewer/config"
	"learnDB/internal/testReviewer/domain"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type ExcelRepository struct {
	bookRead   string
	sheetRead  string
	bookWrite  string
	sheetWrite string
}

func (er *ExcelRepository) Read(xlsx *config.ExcelConfig) ([]domain.StudentWork, error) {

	f, err := excelize.OpenFile(er.bookRead)
	if err != nil {
		return nil, fmt.Errorf("excel open error: %w", err)
	}

	defer f.Close()

	rows, err := f.GetRows(er.sheetRead)
	if err != nil {
		return nil, fmt.Errorf("excel read error: %s", err)
	}

	swSlice := make([]domain.StudentWork, 0, len(rows))

	for i := 2; i <= len(rows); i++ {
		row := strconv.Itoa(i)

		var name string
		for _, col := range xlsx.Name {
			namePart, err := f.GetCellValue(er.sheetRead, col+row)
			if err != nil {
				return nil, fmt.Errorf("excel read error: %v; cell: %s", err, col+row)
			}
			name += namePart
		}
		group, err := f.GetCellValue(er.sheetRead, xlsx.Group+row)
		if err != nil {
			return nil, fmt.Errorf("excel read error: %v; cell: %s", err, xlsx.Group+row)
		}
		sTasks := make([]domain.TestTask, 0, len(xlsx.Tasks))
		for _, task := range xlsx.Tasks {
			question, err := f.GetCellValue(er.sheetRead, task.Question+row)
			if err != nil {
				return nil, fmt.Errorf("excel read error: %v; cell: %s", err, task.Question+row)
			}
			answer, err := f.GetCellValue(er.sheetRead, task.Answer+row)
			if err != nil {
				return nil, fmt.Errorf("excel read error: %v; cell: %s", err, task.Answer+row)
			}
			corrAns, err := f.GetCellValue(er.sheetRead, task.CorrectAnswer+row)
			if err != nil {
				return nil, fmt.Errorf("excel read error: %v; cell: %s", err, task.CorrectAnswer+row)
			}
			corrAnswers := make([]domainAnswer.CorrectAnswer, 0)
			err = json.Unmarshal([]byte(corrAns), &corrAnswers)
			if err != nil {
				return nil, fmt.Errorf("excel read error: json unmarshall error: %v; try unmarshall: %s", err, corrAns)
			}

			sTasks = append(sTasks, domain.TestTask{
				QuestionText:   question,
				Answer:         answer,
				CorrectAnswers: corrAnswers,
			})
		}
		swSlice = append(swSlice, domain.StudentWork{
			Name:  name,
			Group: group,
			Tasks: sTasks,
		})
	}

	return swSlice, nil

}

func (er *ExcelRepository) Write(rw domain.ReviewedWorks) error {
	f := excelize.NewFile()
	defer f.Close()

	if er.sheetWrite != "Sheet1" {
		_, err := f.NewSheet(er.sheetWrite)
		if err != nil {
			return fmt.Errorf("excel write error: cannot create sheet \"%s\" error: %v", er.sheetWrite, err)
		}
	}

	for i, sw := range rw {
		row := strconv.Itoa(i + 1)
		if err := f.SetCellStr(er.sheetWrite, "A"+row, sw.Name); err != nil {
			return fmt.Errorf("excel write error: set name \"%v\" error: %v", sw.Name, err)
		}
		if err := f.SetCellStr(er.sheetWrite, "B"+row, sw.Group); err != nil {
			return fmt.Errorf("excel write error: set group \"%v\" error: %v", sw.Group, err)
		}
		if err := f.SetCellInt(er.sheetWrite, "C"+row, sw.TotalGrade); err != nil {
			return fmt.Errorf("excel write error: set total grade \"%v\" error: %v", sw.TotalGrade, err)
		}

		offset := 0
		for j := 0; j < len(sw.Tasks); j++ {
			task := sw.Tasks[j]
			review := sw.WorkReview[j]

			// Omitting error cause it will mess the code while not having much affect on algorithm
			cell, _ := excelize.CoordinatesToCellName(4+j+offset, i+1)
			offset++
			if err := f.SetCellStr(er.sheetWrite, cell, task.QuestionText); err != nil {
				return fmt.Errorf("excel write error: task %v, error %v", task, err)
			}

			cell, _ = excelize.CoordinatesToCellName(4+j+offset, i+1)
			offset++
			if err := f.SetCellStr(er.sheetWrite, cell, task.Answer); err != nil {
				return fmt.Errorf("excel write error: task %v, error %v", task, err)
			}

			cell, _ = excelize.CoordinatesToCellName(4+j+offset, i+1)
			offset++
			corrAnss, err := json.Marshal(task.CorrectAnswers)
			if err != nil {
				return fmt.Errorf("excel write error: task %v, json marshal error: %v", task, err)
			}
			if err := f.SetCellStr(er.sheetWrite, cell, string(corrAnss)); err != nil {
				return fmt.Errorf("excel write error: correct answer %v, error %v", corrAnss, err)
			}

			cell, _ = excelize.CoordinatesToCellName(4+j+offset, i+1)
			offset++
			if err := f.SetCellInt(er.sheetWrite, cell, review.Points); err != nil {
				return fmt.Errorf("excel write error: review %v, error %v", review, err)
			}

			cell, _ = excelize.CoordinatesToCellName(4+j+offset, i+1)
			if err := f.SetCellStr(er.sheetWrite, cell, review.Error.Error()); err != nil {
				return fmt.Errorf("excel write error: review %v, error %v", review, err)
			}

		}
	}
	return f.SaveAs(er.bookWrite)
}
