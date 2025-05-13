package excel

import (
	"learnDB/internal/testReviewer/config"
	"testing"
)

func TestExcelRead(t *testing.T) {
	staticPath := "/home/pochka/projects/learnDB/excelFiles/"
	config := config.MustLoadConfig(staticPath + "excelConfig.json")

	tr := &ExcelRepository{
		bookRead:  staticPath + "test.xlsx",
		sheetRead: "Sheet1",
	}

	res, err := tr.ReadFile(config)
	if err != nil {
		t.Errorf("excelReader read error: %v", err)
	}
	if res == nil {
		t.Errorf("got %v, want to nil", res)
	}
	if len(res) == 0 {
		t.Fatalf("len(res) == 0, res: %v", res)
	}
}

func TestExcelWrite(t *testing.T) {
	staticPath := "/home/pochka/projects/learnDB/excelFiles/"
	config := config.MustLoadConfig(staticPath + "excelConfig.json")

	tr := &ExcelRepository{
		bookRead:   staticPath + "test.xlsx",
		sheetRead:  "Sheet1",
		bookWrite:  staticPath + "reviewedWorks.xlsx",
		sheetWrite: "Sheet1",
	}
	studentWorks, err := tr.ReadFile(config)
	if err != nil {
		t.Fatalf("excelReader read error: %v", err)
	}

	// reviewedWorks := make(domain.ReviewedWorks, 0, len(studentWorks))
	// for _, work := range studentWorks {
	// 	taskReviews := make([]domain.TaskReview, 0, len(work.Tasks))
	// 	for _, task := range work.Tasks {
	// 		taskReviews = append(taskReviews, domain.TaskReview{
	// 			Task: task,
	// 			Review: domain.CheckResult{
	// 				Points: 3,
	// 				Error:  nil,
	// 			},
	// 		})
	// 	}
	// 	reviewedWorks = append(reviewedWorks, domain.WorkReview{
	// 		Name:       work.Name,
	// 		Group:      work.Group,
	// 		DB:         work.DB,
	// 		TotalGrade: 12,
	// 		Tasks:      taskReviews,
	// 	})
	// }

	err = tr.WriteFile(studentWorks)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

}
