package controller

import (
	"encoding/json"
	"learnDB/internal/scriptManager"
	"log/slog"
	"net/http"
)

type ScriptManagerDTO struct {
	Sql        string `json:"sql"`
	DBMS       string `json:"dbms"`
	SchemaName string `json:"schema_name"`
	TaskName   string `json:"task_name"`
}

type ScriptManagerController struct {
	logger *slog.Logger
	sm     *scriptManager.ScriptManager
}

func NewScriptManagerController(log *slog.Logger, sm *scriptManager.ScriptManager) *ScriptManagerController {
	return &ScriptManagerController{
		logger: log,
		sm:     sm,
	}
}

func (c *ScriptManagerController) RunScriptHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	dto := new(ScriptManagerDTO)

	if err := json.NewDecoder(r.Body).Decode(dto); err != nil {
		c.logger.Error("parse json error", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"msg":     "InternalServerError",
		})
		return
	}

	opts := c.sm.GetOpts()

	opts.Sql = dto.Sql
	opts.DBMS = dto.DBMS
	opts.SchemaName = dto.SchemaName

	res := c.sm.RunScript(opts)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (c *ScriptManagerController) ValidateScriptHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	dto := new(ScriptManagerDTO)

	if err := json.NewDecoder(r.Body).Decode(dto); err != nil {
		c.logger.Error("parse json error", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"msg":     "InternalServerError",
		})
		return
	}

	opts := c.sm.GetOpts()

	opts.Sql = dto.Sql
	opts.DBMS = dto.DBMS
	opts.SchemaName = dto.SchemaName

	right, comment, err := c.sm.ValidateScript(opts)

	if err != nil {
		c.logger.Error("validate process error", slog.Any("error", err), slog.Any("params", opts))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"msg":     "InternalServerError",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"right_answer": right,
		"comment":      comment,
	})

}

func (c *ScriptManagerController) GetRandomTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	task, err := c.sm.GetRandomTask()

	if err != nil {
		c.logger.Error("fetching random task error", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"msg":     "InternalServerError",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"task": map[string]interface{}{
			"name":     task.Name,
			"question": task.Question,
		},
	})
}
