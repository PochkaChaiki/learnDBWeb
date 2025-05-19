package config_test

import (
	"learnDB/internal/testReviewer/config"
	"reflect"
	"testing"
)

func TestExcelConfigLoad(t *testing.T) {
	configPath := "/home/pochka/projects/learnDB/excelFiles/testExcelConfig.json"
	result := config.MustLoadConfig(configPath)
	want := &config.ExcelConfig{
		Name:  []string{"A", "B"},
		Group: "C",
		DB:    "postgres",
		Tasks: []config.Task{
			{
				Question:      "L",
				Answer:        "M",
				CorrectAnswer: "N",
				Points:        3,
			},
			{
				Question:      "O",
				Answer:        "P",
				CorrectAnswer: "Q",
				Points:        3,
			},
			{
				Question:      "R",
				Answer:        "S",
				CorrectAnswer: "T",
				Points:        3,
			},
		},
	}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("got %v, want: %v", result, want)
	}
}
