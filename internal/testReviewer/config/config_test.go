package config_test

import (
	"learnDB/internal/testReviewer/config"
	"reflect"
	"testing"
)

func TestExcelConfigLoad(t *testing.T) {
	configPath := "/home/pochka/projects/learnDB/static/excelConfig.json"
	result := config.MustLoadConfig(configPath)
	want := &config.ExcelConfig{
		Name:  []string{"A", "B"},
		Group: "C",
		Tasks: []config.Task{
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
