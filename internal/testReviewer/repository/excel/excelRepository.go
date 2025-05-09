package excel

import (
	"encoding/json"
	"errors"
	"fmt"
	mainDomain "learnDB/internal/domain"
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

func New(bookread string, sheetread string, bookwrite string, sheetwrite string) *ExcelRepository {
	return &ExcelRepository{
		bookRead:   bookread,
		sheetRead:  sheetread,
		bookWrite:  bookwrite,
		sheetWrite: sheetwrite,
	}
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
		db, err := f.GetCellValue(er.sheetRead, xlsx.DB+row)
		if err != nil {
			return nil, fmt.Errorf("excel read error: %v; cell: %s", err, xlsx.DB+row)
		}
		// REMOVE IT AFTER CHECKING ----------------------------------------------------------------------------
		switch db {
		case "postgresql":
			db = "postgres"
		case "Мне помогали":
			db = "postgres"
		case "-":
			db = "postgres"
		default:
		}
		// -----------------------------------------------------------------------------------------------------
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
			// REMOVE IT AFTER CHECKING ----------------------------------------------------------------------------
			corrAns = fmt.Sprintf("[{\"values\":[\"%s\"], \"points\": 3}]", corrAns)
			// -----------------------------------------------------------------------------------------------------
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
			DB:    db,
		})
	}

	return swSlice, nil

}

func (er *ExcelRepository) Write(reviewedWorks domain.ReviewedWorks) error {
	const initialCell int = 5
	f := excelize.NewFile()
	defer f.Close()

	if er.sheetWrite != "Sheet1" {
		_, err := f.NewSheet(er.sheetWrite)
		if err != nil {
			return fmt.Errorf("excel write error: cannot create sheet \"%s\" error: %v", er.sheetWrite, err)
		}
	}

	// Set header
	if err := f.SetCellStr(er.sheetWrite, "A1", "Name"); err != nil {
		return fmt.Errorf("excel write error: set header name \"%v\" error: %v", "Name", err)
	}
	if err := f.SetCellStr(er.sheetWrite, "B1", "Group"); err != nil {
		return fmt.Errorf("excel write error: set header name \"%v\" error: %v", "Group", err)
	}
	if err := f.SetCellStr(er.sheetWrite, "C1", "TotalGrade"); err != nil {
		return fmt.Errorf("excel write error: set header name \"%v\" error: %v", "TotalGrade", err)
	}
	if err := f.SetCellStr(er.sheetWrite, "D1", "Database"); err != nil {
		return fmt.Errorf("excel write error: set header name \"%v\" error: %v", "Database", err)
	}

	offset := 0
	for i := 0; i < len(reviewedWorks[0].Tasks); i++ {
		cell, _ := excelize.CoordinatesToCellName(initialCell+i+offset, 1)
		num := strconv.Itoa(i + 1)
		if err := f.SetCellStr(er.sheetWrite, cell, "Question"+num); err != nil {
			return fmt.Errorf("excel write error: set header name \"%v\" error: %v", "Question"+num, err)
		}
		offset++
		cell, _ = excelize.CoordinatesToCellName(initialCell+i+offset, 1)
		if err := f.SetCellStr(er.sheetWrite, cell, "Answer"+num); err != nil {
			return fmt.Errorf("excel write error: set header name \"%v\" error: %v", "Answer"+num, err)
		}
		offset++
		cell, _ = excelize.CoordinatesToCellName(initialCell+i+offset, 1)
		if err := f.SetCellStr(er.sheetWrite, cell, "CorrectAnswer"+num); err != nil {
			return fmt.Errorf("excel write error: set header name \"%v\" error: %v", "CorrectAnswer"+num, err)
		}
		offset++
		cell, _ = excelize.CoordinatesToCellName(initialCell+i+offset, 1)
		if err := f.SetCellStr(er.sheetWrite, cell, "Points"+num); err != nil {
			return fmt.Errorf("excel write error: set header name \"%v\" error: %v", "Points"+num, err)
		}
		offset++
		cell, _ = excelize.CoordinatesToCellName(initialCell+i+offset, 1)
		if err := f.SetCellStr(er.sheetWrite, cell, "Error"+num); err != nil {
			return fmt.Errorf("excel write error: set header name \"%v\" error: %v", "Error"+num, err)
		}
	}

	// Putting students' works reviews to excel
	for i, workReview := range reviewedWorks {
		row := strconv.Itoa(i + 2)
		if err := f.SetCellStr(er.sheetWrite, "A"+row, workReview.Name); err != nil {
			return fmt.Errorf("excel write error: set name \"%v\" error: %v", workReview.Name, err)
		}
		if err := f.SetCellStr(er.sheetWrite, "B"+row, workReview.Group); err != nil {
			return fmt.Errorf("excel write error: set group \"%v\" error: %v", workReview.Group, err)
		}
		if err := f.SetCellInt(er.sheetWrite, "C"+row, workReview.TotalGrade); err != nil {
			return fmt.Errorf("excel write error: set total grade \"%v\" error: %v", workReview.TotalGrade, err)
		}
		if err := f.SetCellStr(er.sheetWrite, "D"+row, workReview.DB); err != nil {
			return fmt.Errorf("excel write error: set total grade \"%v\" error: %v", workReview.TotalGrade, err)
		}

		offset = 0
		for j, taskReview := range workReview.Tasks {
			task := taskReview.Task
			review := taskReview.Review

			// Omitting error cause it will mess the code while not having much affect on algorithm
			cell, _ := excelize.CoordinatesToCellName(initialCell+j+offset, i+2)
			firstCell := cell
			offset++
			if err := f.SetCellStr(er.sheetWrite, cell, task.QuestionText); err != nil {
				return fmt.Errorf("excel write error: task %v, error %v", task, err)
			}

			cell, _ = excelize.CoordinatesToCellName(initialCell+j+offset, i+2)
			offset++
			if err := f.SetCellStr(er.sheetWrite, cell, task.Answer); err != nil {
				return fmt.Errorf("excel write error: task %v, error %v", task, err)
			}

			cell, _ = excelize.CoordinatesToCellName(initialCell+j+offset, i+2)
			offset++
			corrAnss, err := json.Marshal(task.CorrectAnswers)
			if err != nil {
				return fmt.Errorf("excel write error: task %v, json marshal error: %v", task, err)
			}
			if err := f.SetCellStr(er.sheetWrite, cell, string(corrAnss)); err != nil {
				return fmt.Errorf("excel write error: correct answer %v, error %v", corrAnss, err)
			}

			cell, _ = excelize.CoordinatesToCellName(initialCell+j+offset, i+2)
			offset++
			if err := f.SetCellInt(er.sheetWrite, cell, review.Points); err != nil {
				return fmt.Errorf("excel write error: review %v, error %v", review, err)
			}

			cell, _ = excelize.CoordinatesToCellName(initialCell+j+offset, i+2)
			errorStr := fmt.Sprintf("%v", review.Error)
			if review.Error == nil {
				errorStr = ""
			}
			if err := f.SetCellStr(er.sheetWrite, cell, errorStr); err != nil {
				return fmt.Errorf("excel write error: review %v, error %v", review, err)
			}

			var color []string

			switch {
			case errors.Is(review.Error, mainDomain.ErrEmptyAnswer):
				color = []string{"D9D9D9"}
			case errors.Is(review.Error, mainDomain.ErrIncorrectAnswer):
				color = []string{"FF3300"}
			case errors.Is(review.Error, mainDomain.ErrSyntaxError):
				color = []string{"C65911"}
			case errors.Is(review.Error, mainDomain.ErrInternalError):
				color = []string{"33CCFF"}
			}

			style, err := f.NewStyle(&excelize.Style{
				Fill: excelize.Fill{Type: "pattern", Color: color, Pattern: 1},
			})
			if err != nil {
				return fmt.Errorf("excel write style error: review %v, error %v", review, err)
			}
			f.SetCellStyle(er.sheetWrite, firstCell, cell, style)

		}
	}
	return f.SaveAs(er.bookWrite)
}
