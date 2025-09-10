package api

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/pprof"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pochkachaiki/learndb/internal/api/controller"
	"github.com/pochkachaiki/learndb/internal/config"
	"github.com/pochkachaiki/learndb/internal/dbmgr"
	"github.com/pochkachaiki/learndb/internal/scriptmgr"
	"github.com/pochkachaiki/learndb/internal/sysdb/postgres"
	"github.com/pochkachaiki/learndb/internal/testreview"
)

type APIServer struct {
	cfg *config.Config

	logger *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) *APIServer {

	return &APIServer{
		cfg:    cfg,
		logger: logger,
	}
}

func (api *APIServer) Run(ctx context.Context) {

	manager := dbmgr.New()

	for name, connStr := range api.cfg.Databases {
		api.logger.Info("connecting to database", slog.String("db", name))
		err := manager.AddRepository(name, connStr)
		if err != nil {
			api.logger.Warn("cannot create db repository", slog.String("db", name), slog.Any("error", err))
		}
	}

	reviewer := testreview.New(manager, "olympics")
	reviewerController := controller.NewTestReviewerController(api.logger, reviewer)

	db, err := pgxpool.New(ctx, api.cfg.SystemDB)
	if err != nil {
		api.logger.Error("cannot create system db repository", slog.String("db", "postgres"), slog.Any("error", err))
	}

	sysRepo := postgres.NewRepository(db)

	scriptMan := scriptmgr.New(manager, sysRepo, 200)
	scriptController := controller.NewScriptManagerController(api.logger, scriptMan)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /upload", reviewerController.ExcelHandler)
	mux.HandleFunc("POST /scriptRunner/runScript", scriptController.RunScriptHandler)
	mux.HandleFunc("POST /scriptRunner/validateScript", scriptController.ValidateScriptHandler)
	mux.HandleFunc("GET /task", scriptController.GetRandomTask)

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	mux.Handle("/", http.FileServer(http.Dir("./static")))

	s := http.Server{
		Addr:    api.cfg.Address,
		Handler: mux,
	}

	api.logger.Info("start server", slog.String("address", api.cfg.Address))
	if err := s.ListenAndServe(); err != nil {
		api.logger.Error("server run error", slog.Any("error", err))
	}

}
