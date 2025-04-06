package service

import (
	"errors"
	"learnDB/internal/testReviewer/config"
	"learnDB/internal/testReviewer/domain"
	"testing"
)

func TestExcelRead(t *testing.T) {
	staticPath := "/home/pochka/projects/learnDB/static/"
	config := config.MustLoadConfig(staticPath + "excelConfig.json")

	tr := new(ExcelReader)

	res, err := tr.Read(staticPath+"testMod.xlsx", "Sheet1", config)
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

	tr := new(ExcelReader)
	sw, err := tr.Read(staticPath+"testMod.xlsx", "Sheet1", config)
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

	err = tr.Write(staticPath+"reviewedWorks.xlsx", "Sheet1", rws)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

}
