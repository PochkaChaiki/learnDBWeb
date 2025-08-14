package excel

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

	mainDomain "github.com/pochkachaiki/learndb/internal/domain"
	"github.com/pochkachaiki/learndb/internal/domain/answer"
	"github.com/pochkachaiki/learndb/internal/testreview/config"
	"github.com/pochkachaiki/learndb/internal/testreview/domain"

	"regexp"

	"github.com/xuri/excelize/v2"
)

const (
	scriptExpr = "(?i)(with[^;]*)?select(([\\'\\\"]\\s?;\\s?[\\'\\\"])|[^;])*"
)

var (
	ErrAnswerNotPresent = errors.New("answer is not present")
)

type ExcelRepository struct {
	bookRead   string
	sheetRead  string
	bookWrite  string
	sheetWrite string
	works      map[*domain.Work]*StudentWork
	scriptRe   *regexp.Regexp
	xlsx       *config.ExcelConfig
}

func NewFileRepo(bookread string, sheetread string, bookwrite string, sheetwrite string, xlsx *config.ExcelConfig) *ExcelRepository {
	scriptRe, err := regexp.Compile(scriptExpr)
	if err != nil {
		return nil
	}

	return &ExcelRepository{
		bookRead:   bookread,
		sheetRead:  sheetread,
		bookWrite:  bookwrite,
		sheetWrite: sheetwrite,
		scriptRe:   scriptRe,
		xlsx:       xlsx,
	}
}

func New(sheet string, xlsx *config.ExcelConfig) *ExcelRepository {
	scriptRe, err := regexp.Compile(scriptExpr)
	if err != nil {
		return nil
	}

	return &ExcelRepository{
		sheetRead:  sheet,
		sheetWrite: sheet,
		scriptRe:   scriptRe,
		xlsx:       xlsx,
	}
}

func (ex *ExcelRepository) retrieveScripts(ans string) []string {
	return ex.scriptRe.FindAllString(ans, -1)
}

func (ex *ExcelRepository) answerContains(ans string, corrAnsValue string) bool {
	res, _ := regexp.MatchString(fmt.Sprintf("(?i)(^(%s)|(%[1]s)$)", corrAnsValue), ans)
	return res
}

func (er *ExcelRepository) ReadFile() (domain.Works, error) {

	f, err := excelize.OpenFile(er.bookRead)
	if err != nil {
		return nil, fmt.Errorf("excel open error: %w", err)
	}

	defer f.Close()
	return er.read(f)
}

func (er *ExcelRepository) Read(r io.Reader) (domain.Works, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("excel open reader error: %w", err)
	}
	defer f.Close()

	return er.read(f)
}

func (er *ExcelRepository) read(f *excelize.File) (domain.Works, error) {

	rows, err := f.GetRows(er.sheetRead)
	if err != nil {
		return nil, fmt.Errorf("excel read error: %s", err)
	}

	er.works = make(map[*domain.Work]*StudentWork)

	workSlice := make(domain.Works, 0, len(rows))

	for i := 2; i <= len(rows); i++ {

		row := strconv.Itoa(i)

		dbInstall := 2
		var name string
		for _, col := range er.xlsx.Name {
			namePart, err := f.GetCellValue(er.sheetRead, col+row)
			if err != nil {
				return nil, fmt.Errorf("excel read error: %v; cell: %s", err, col+row)
			}
			name += namePart
		}
		group, err := f.GetCellValue(er.sheetRead, er.xlsx.Group+row)
		if err != nil {
			return nil, fmt.Errorf("excel read error: %v; cell: %s", err, er.xlsx.Group+row)
		}
		db, err := f.GetCellValue(er.sheetRead, er.xlsx.DB+row)
		if err != nil {
			return nil, fmt.Errorf("excel read error: %v; cell: %s", err, er.xlsx.DB+row)
		}
		// REMOVE IT AFTER CHECKING ----------------------------------------------------------------------------
		switch db {
		case "postgresql":
			db = "postgres"
		case "Мне помогали":
			db = "postgres"
			dbInstall = 1
		case "-":
			db = "postgres"
			dbInstall = 0
		default:
		}
		// -----------------------------------------------------------------------------------------------------
		sTasks := make([]Task, 0, len(er.xlsx.Tasks))
		wUnits := make([]*domain.CheckUnit, 0, len(er.xlsx.Tasks))
		for _, task := range er.xlsx.Tasks {

			var review domain.CheckResult

			// Reading question
			question, err := f.GetCellValue(er.sheetRead, task.Question+row)
			if err != nil {
				return nil, fmt.Errorf("excel read error: %v; cell: %s", err, task.Question+row)
			}

			// Read answer
			answerFromCell, err := f.GetCellValue(er.sheetRead, task.Answer+row)
			if err != nil {
				return nil, fmt.Errorf("excel read error: %v; cell: %s", err, task.Answer+row)
			}

			// Read correct answer
			corrAnsFromCell, err := f.GetCellValue(er.sheetRead, task.CorrectAnswer+row)
			if err != nil {
				return nil, fmt.Errorf("excel read error: %v; cell: %s", err, task.CorrectAnswer+row)
			}

			corrAnswers := make([]answer.CorrectAnswer, 0)
			if err = json.Unmarshal([]byte(corrAnsFromCell), &corrAnswers); err != nil {
				// Stub for cases when answers are kept in incorrect way
				corrAns := fmt.Sprintf("[{\"values\":[\"%s\"], \"points\": %d}]", corrAnsFromCell, task.Points)
				if err = json.Unmarshal([]byte(corrAns), &corrAnswers); err != nil {
					return nil, fmt.Errorf("excel read error: json unmarshall error: %v; try unmarshall: %s", err, corrAns)
				}
			}

			// Parse answer
			var script string
			sql := er.retrieveScripts(answerFromCell)
			if len(sql) == 0 {
				script = ""
			} else {
				script = sql[len(sql)-1]
			}
			if len(sql) > 1 {
				review.Points--
			}

			// Check corr answer
			if !er.answerContains(answerFromCell, corrAnsFromCell) {
				review.Error = ErrAnswerNotPresent
			}

			sTasks = append(sTasks, Task{
				Question:       question,
				Answer:         answerFromCell,
				CorrectAnswers: corrAnswers,
			})
			unit := new(domain.CheckUnit)
			unit.Script = script
			unit.CorrectAnswers = corrAnswers
			unit.Review = review
			wUnits = append(wUnits, unit)
		}

		sw := new(StudentWork)
		sw.Name = name
		sw.Group = group
		sw.Tasks = sTasks

		work := new(domain.Work)
		work.DB = db
		work.TotalGrade = dbInstall
		work.Units = wUnits

		er.works[work] = sw

		workSlice = append(workSlice, work)
	}

	return workSlice, nil

}

func (er *ExcelRepository) WriteFile(reviewedWorks domain.Works) error {
	f := excelize.NewFile()
	defer f.Close()

	if err := er.write(f, reviewedWorks); err != nil {
		return err
	}

	return f.SaveAs(er.bookWrite)
}

func (er *ExcelRepository) Write(w io.Writer, reviewedWorks domain.Works) error {
	f := excelize.NewFile()
	defer f.Close()

	if err := er.write(f, reviewedWorks); err != nil {
		return err
	}

	return f.Write(w)
}

func (er *ExcelRepository) write(f *excelize.File, reviewedWorks domain.Works) error {
	const initialCell int = 5

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
	for i := 0; i < len(reviewedWorks[0].Units); i++ {
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
		if err := f.SetCellStr(er.sheetWrite, cell, "ContainsCorrectAnswer"+num); err != nil {
			return fmt.Errorf("excel write error: set header name \"%v\" error: %v", "ContainsCorrectAnswer"+num, err)
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
	for i, work := range reviewedWorks {
		studentWork := er.works[work]

		row := strconv.Itoa(i + 2)
		if err := f.SetCellStr(er.sheetWrite, "A"+row, studentWork.Name); err != nil {
			return fmt.Errorf("excel write error: set name \"%v\" error: %v", studentWork.Name, err)
		}
		if err := f.SetCellStr(er.sheetWrite, "B"+row, studentWork.Group); err != nil {
			return fmt.Errorf("excel write error: set group \"%v\" error: %v", studentWork.Group, err)
		}
		if err := f.SetCellInt(er.sheetWrite, "C"+row, work.TotalGrade); err != nil {
			return fmt.Errorf("excel write error: set total grade \"%v\" error: %v", work.TotalGrade, err)
		}
		if err := f.SetCellStr(er.sheetWrite, "D"+row, work.DB); err != nil {
			return fmt.Errorf("excel write error: set total grade \"%v\" error: %v", work.DB, err)
		}

		offset = 0
		for j, taskReview := range studentWork.Tasks {
			task := taskReview
			review := work.Units[j].Review

			// Omitting error cause it will mess the code while not having much affect on algorithm
			// Write Question
			cell, _ := excelize.CoordinatesToCellName(initialCell+j+offset, i+2)
			firstCell := cell
			offset++
			if err := f.SetCellStr(er.sheetWrite, cell, task.Question); err != nil {
				return fmt.Errorf("excel write error: task %v, error %v", task, err)
			}

			// Write Answer
			cell, _ = excelize.CoordinatesToCellName(initialCell+j+offset, i+2)
			offset++
			if err := f.SetCellStr(er.sheetWrite, cell, task.Answer); err != nil {
				return fmt.Errorf("excel write error: task %v, error %v", task, err)
			}

			// Write Correct Answer
			cell, _ = excelize.CoordinatesToCellName(initialCell+j+offset, i+2)
			offset++
			corrAnss, err := json.Marshal(task.CorrectAnswers)
			if err != nil {
				return fmt.Errorf("excel write error: task %v, json marshal error: %v", task, err)
			}
			if err := f.SetCellStr(er.sheetWrite, cell, string(corrAnss)); err != nil {
				return fmt.Errorf("excel write error: correct answer %v, error %v", corrAnss, err)
			}

			// Write Does Correct Answer is present at Student's Answer
			cell, _ = excelize.CoordinatesToCellName(initialCell+j+offset, i+2)
			offset++
			contains := !errors.Is(review.Error, ErrAnswerNotPresent)
			if err := f.SetCellBool(er.sheetWrite, cell, contains); err != nil {
				return fmt.Errorf("excel write error: review %v, error %v", review, err)
			}

			// Write Points
			cell, _ = excelize.CoordinatesToCellName(initialCell+j+offset, i+2)
			offset++
			if err := f.SetCellInt(er.sheetWrite, cell, review.Points); err != nil {
				return fmt.Errorf("excel write error: review %v, error %v", review, err)
			}

			// Write Errors
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
	return nil
}
