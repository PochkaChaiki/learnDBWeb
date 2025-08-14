package controller

import (
	"fmt"

	"github.com/pochkachaiki/learndb/internal/testreview"
	"github.com/pochkachaiki/learndb/internal/testreview/config"
	"github.com/pochkachaiki/learndb/internal/testreview/datasrc/excel"

	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	MAX_BYTES int64 = 1 << 20
)

type TestReviewerController struct {
	logger *slog.Logger
	tr     *testreview.TestReviewer
}

func NewTestReviewerController(logger *slog.Logger, tr *testreview.TestReviewer) *TestReviewerController {
	return &TestReviewerController{logger: logger, tr: tr}
}

func (c *TestReviewerController) ExcelHandler(w http.ResponseWriter, r *http.Request) {
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

	jsonConfig := new(config.ExcelConfig)
	if err := cleanenv.ParseJSON(jsonFile, jsonConfig); err != nil {
		c.logger.Error("parse json config error", slog.Any("error", err))
		http.Error(w, "Failed to read JSON config", http.StatusBadRequest)
		return
	}

	excelRepo := excel.New(jsonConfig.Sheet, jsonConfig)

	studentWorks, err := excelRepo.Read(excelFile)
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
