package testReader

import (
	"encoding/json"
	"fmt"
	domainAnswer "learnDB/internal/domain/answer"
	"learnDB/internal/testReviewer/domain"
	"strconv"

	"github.com/xuri/excelize/v2"
)

// type ExcelReader struct {
// }

// func (er *ExcelReader) Read(bookname string, xlsx *ExcelConfig) ([]domain.StudentWork, error) {
func Read(bookname string, sheetname string, xlsx *ExcelConfig) ([]domain.StudentWork, error) {
	// bookname := "test.xlsx"

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
