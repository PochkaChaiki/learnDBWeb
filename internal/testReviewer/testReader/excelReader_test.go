package testReader

import (
	"reflect"
	"testing"
)

func TestExcelConfigLoad(t *testing.T) {
	configPath := "/home/pochka/projects/learnDB/static/excelConfig.json"
	result := MustLoad(configPath)
	want := &ExcelConfig{
		Name:  []string{"A", "B"},
		Group: "C",
		Tasks: []Task{
			{
				Question:      "L",
				Answer:        "M",
				CorrectAnswer: "N",
			},
			{
				Question:      "O",
				Answer:        "P",
				CorrectAnswer: "Q",
			},
			{
				Question:      "R",
				Answer:        "S",
				CorrectAnswer: "T",
			},
		},
	}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("got %v, want: %v", result, want)
	}
}

func TestExcelRead(t *testing.T) {
	staticPath := "/home/pochka/projects/learnDB/static"
	config := MustLoad(staticPath + "/excelConfig.json")

	res, err := Read(staticPath+"/testMod.xlsx", "Sheet1", config)
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
