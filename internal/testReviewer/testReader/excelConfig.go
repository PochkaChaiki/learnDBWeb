package testReader

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Task struct {
	Question      string `json:"question"`
	Answer        string `json:"answer"`
	CorrectAnswer string `json:"correct_answer"`
}

type ExcelConfig struct {
	Name  []string `json:"name"`
	Group string   `json:"group"`
	Tasks []Task   `json:"tasks"`
}

func MustLoad(configPath string) *ExcelConfig {
	xlsx := new(ExcelConfig)
	err := cleanenv.ReadConfig(configPath, xlsx)
	if err != nil {
		panic(fmt.Errorf("read excel config error: %w", err))
	}
	return xlsx
}
