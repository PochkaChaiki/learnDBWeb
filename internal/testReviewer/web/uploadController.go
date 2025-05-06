package web

import (
	"learnDB/internal/testReviewer"
	"learnDB/internal/testReviewer/config"
	"learnDB/internal/testReviewer/repository/excel"

	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"

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

func (c *Controller) prepareFileFromResponse(w *http.ResponseWriter, r *http.Request) (tempExcelPath string, tempOutExcelPath string, ok bool) {
	ok = false
	excelFile, excelHeader, err := r.FormFile("excelFile")
	if err != nil {
		c.logger.Error("get excel file error", slog.Any("error", err))
		http.Error(*w, "Excel file is required", http.StatusInternalServerError)
		return
	}
	defer excelFile.Close()

	tempExcelPath = "./temp_" + excelHeader.Filename
	tempOutExcelPath = "./temp_out_" + excelHeader.Filename

	outFile, err := os.Create(tempExcelPath)
	if err != nil {
		c.logger.Error("create temp excel file error", slog.Any("error", err))
		http.Error(*w, "Failed to save Excel file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	if _, err = io.Copy(outFile, excelFile); err != nil {
		c.logger.Error("save excel file error", slog.Any("error", err))
		http.Error(*w, "Failed to save Excel file", http.StatusInternalServerError)
		return
	}
	ok = true
	return
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

	tempExcelPath, tempOutExcelPath, ok := c.prepareFileFromResponse(&w, r)
	if !ok {
		return
	}

	defer os.Remove(tempExcelPath)

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

	exrepo := excel.New(tempExcelPath, jsonConfig.Sheet, tempOutExcelPath, jsonConfig.Sheet)

	studentWorks, err := exrepo.Read(&jsonConfig)
	if err != nil {
		c.logger.Error("cannot read excel", slog.Any("error", err))
		http.Error(w, "Failed to process excel", http.StatusInternalServerError)
		return
	}

	reviewedWorks := c.tr.CheckTest(studentWorks)

	exrepo.Write(reviewedWorks)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Disposition", "attachment; filename=NameOfFile")

	sendFile, err := os.Open(tempOutExcelPath)
	if err != nil {
		c.logger.Error("cannot open excel file to send", slog.Any("error", err))
		http.Error(w, "Failed to process excel", http.StatusInternalServerError)
		return
	}
	defer sendFile.Close()

	if _, err = io.Copy(w, sendFile); err != nil {
		c.logger.Error("cannot send excel file", slog.Any("error", err))
		http.Error(w, "Failed to process excel", http.StatusInternalServerError)
		return
	}

}
