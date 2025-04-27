package excel

import (
	"errors"
	"learnDB/internal/testReviewer/config"
	"learnDB/internal/testReviewer/domain"
	"testing"
)

func TestExcelRead(t *testing.T) {
	staticPath := "/home/pochka/projects/learnDB/static/"
	config := config.MustLoadConfig(staticPath + "excelConfig.json")

	tr := &ExcelRepository{
		bookRead:  staticPath + "testMod.xlsx",
		sheetRead: "Sheet1",
	}

	res, err := tr.Read(config)
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
	staticPath := "/home/pochka/projects/learnDB/static/"
	config := config.MustLoadConfig(staticPath + "excelConfig.json")

	tr := &ExcelRepository{
		bookRead:   staticPath + "testMod.xlsx",
		sheetRead:  "Sheet1",
		bookWrite:  staticPath + "reviewedWorks.xlsx",
		sheetWrite: "Sheet1",
	}
	sw, err := tr.Read(config)
	if err != nil {
		t.Fatalf("excelReader read error: %v", err)
	}

	rws := make(domain.ReviewedWorks, len(sw))
	for i := range rws {
		rws[i].StudentWork = sw[i]
		rws[i].TotalGrade = 10
		wr := make(domain.WorkReview, len(rws[i].Tasks))
		for j := range rws[i].Tasks {
			wr[j].Points = 3
			wr[j].Error = errors.New("test error")
		}
		rws[i].WorkReview = wr
	}

	err = tr.Write(rws)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

}
