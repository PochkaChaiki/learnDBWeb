package web

import (
	"fmt"
	"learnDB/internal/testReviewer"
	"learnDB/internal/testReviewer/config"
	"learnDB/internal/testReviewer/repository/excel"

	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	MAX_BYTES int64 = 1 << 20
)

type Controller struct {
	logger *slog.Logger
	tr     *testReviewer.TestReviewer
}

func NewController(logger *slog.Logger, tr *testReviewer.TestReviewer) *Controller {
	return &Controller{logger: logger, tr: tr}
}

func (c *Controller) ExcelHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(MAX_BYTES); err != nil {
		c.logger.Error("parse files error", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"msg":     "InternalServerError",
		})
		return
	}

	excelFile, excelHeader, err := r.FormFile("excelFile")
	if err != nil {
		c.logger.Error("get excel file error", slog.Any("error", err))
		http.Error(w, "Excel file is required", http.StatusInternalServerError)
		return
	}
	defer excelFile.Close()

	jsonFile, _, err := r.FormFile("jsonConfig")
	if err != nil {
		c.logger.Error("get json file error", slog.Any("error", err))
		http.Error(w, "JSON config is required", http.StatusBadRequest)
		return
	}
	defer jsonFile.Close()

	var jsonConfig config.ExcelConfig
	if err := cleanenv.ParseJSON(jsonFile, &jsonConfig); err != nil {
		c.logger.Error("parse json config error", slog.Any("error", err))
		http.Error(w, "Failed to read JSON config", http.StatusBadRequest)
		return
	}

	excelRepo := excel.New(jsonConfig.Sheet)

	studentWorks, err := excelRepo.Read(excelFile, &jsonConfig)
	if err != nil {
		c.logger.Error("cannot read excel", slog.Any("error", err))
		http.Error(w, "Failed to process excel", http.StatusInternalServerError)
		return
	}

	c.tr.CheckTest(studentWorks)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", excelHeader.Filename))
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	if err := excelRepo.Write(w, studentWorks); err != nil {
		c.logger.Error("cannot write excel", slog.Any("error", err))
		http.Error(w, "Failed to process excel", http.StatusInternalServerError)
		return
	}

}
