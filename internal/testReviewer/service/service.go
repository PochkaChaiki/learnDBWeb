package service

import (
	"encoding/json"
	"fmt"
	domainAnswer "learnDB/internal/domain/answer"
	"learnDB/internal/testReviewer/config"
	"learnDB/internal/testReviewer/domain"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type ExcelReader struct {
}

func (er *ExcelReader) Read(bookname string, sheetname string, xlsx *config.ExcelConfig) ([]domain.StudentWork, error) {

	f, err := excelize.OpenFile(bookname)
	if err != nil {
		return nil, fmt.Errorf("excel open error: %w", err)
	}

	defer f.Close()

	rows, err := f.GetRows(sheetname)
	if err != nil {
		return nil, fmt.Errorf("excel read error: %s", err)
	}

	swSlice := make([]domain.StudentWork, 0, len(rows))

	for i := 2; i <= len(rows); i++ {
		row := strconv.Itoa(i)

		var name string
		for _, col := range xlsx.Name {
			namePart, err := f.GetCellValue(sheetname, col+row)
			if err != nil {
				return nil, fmt.Errorf("excel read error: %v; cell: %s", err, col+row)
			}
			name += namePart
		}
		group, err := f.GetCellValue(sheetname, xlsx.Group+row)
		if err != nil {
			return nil, fmt.Errorf("excel read error: %v; cell: %s", err, xlsx.Group+row)
		}
		sTasks := make([]domain.TestTask, 0, len(xlsx.Tasks))
		for _, task := range xlsx.Tasks {
			question, err := f.GetCellValue(sheetname, task.Question+row)
			if err != nil {
				return nil, fmt.Errorf("excel read error: %v; cell: %s", err, task.Question+row)
			}
			answer, err := f.GetCellValue(sheetname, task.Answer+row)
			if err != nil {
				return nil, fmt.Errorf("excel read error: %v; cell: %s", err, task.Answer+row)
			}
			corrAns, err := f.GetCellValue(sheetname, task.CorrectAnswer+row)
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

func (er *ExcelReader) Write(bookname string, sheetname string, rw domain.ReviewedWorks) error {
	f := excelize.NewFile()
	defer f.Close()

	if sheetname != "Sheet1" {
		_, err := f.NewSheet(sheetname)
		if err != nil {
			return fmt.Errorf("excel write error: cannot create sheet \"%s\" error: %v", sheetname, err)
		}
	}

	for i, sw := range rw {
		row := strconv.Itoa(i + 1)
		if err := f.SetCellStr(sheetname, "A"+row, sw.Name); err != nil {
			return fmt.Errorf("excel write error: set name \"%v\" error: %v", sw.Name, err)
		}
		if err := f.SetCellStr(sheetname, "B"+row, sw.Group); err != nil {
			return fmt.Errorf("excel write error: set group \"%v\" error: %v", sw.Group, err)
		}
		if err := f.SetCellInt(sheetname, "C"+row, sw.TotalGrade); err != nil {
			return fmt.Errorf("excel write error: set total grade \"%v\" error: %v", sw.TotalGrade, err)
		}

		offset := 0
		for j := 0; j < len(sw.Tasks); j++ {
			task := sw.Tasks[j]
			review := sw.WorkReview[j]

			// Omitting error cause it will mess the code while not having much affect on algorithm
			cell, _ := excelize.CoordinatesToCellName(4+j+offset, i+1)
			offset++
			if err := f.SetCellStr(sheetname, cell, task.QuestionText); err != nil {
				return fmt.Errorf("excel write error: task %v, error %v", task, err)
			}

			cell, _ = excelize.CoordinatesToCellName(4+j+offset, i+1)
			offset++
			if err := f.SetCellStr(sheetname, cell, task.Answer); err != nil {
				return fmt.Errorf("excel write error: task %v, error %v", task, err)
			}

			cell, _ = excelize.CoordinatesToCellName(4+j+offset, i+1)
			offset++
			corrAnss, err := json.Marshal(task.CorrectAnswers)
			if err != nil {
				return fmt.Errorf("excel write error: task %v, json marshal error: %v", task, err)
			}
			if err := f.SetCellStr(sheetname, cell, string(corrAnss)); err != nil {
				return fmt.Errorf("excel write error: correct answer %v, error %v", corrAnss, err)
			}

			cell, _ = excelize.CoordinatesToCellName(4+j+offset, i+1)
			offset++
			if err := f.SetCellInt(sheetname, cell, review.Points); err != nil {
				return fmt.Errorf("excel write error: review %v, error %v", review, err)
			}

			cell, _ = excelize.CoordinatesToCellName(4+j+offset, i+1)
			if err := f.SetCellStr(sheetname, cell, review.Error.Error()); err != nil {
				return fmt.Errorf("excel write error: review %v, error %v", review, err)
			}

		}
	}
	return f.SaveAs(bookname)
}
