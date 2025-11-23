package main

import (
	"log/slog"
	"net/http"
	"os"

	"api/internal/config"
	"api/internal/http-server/handlers/courses/get"
	"api/internal/http-server/handlers/courses/save"
	"api/internal/lib/logger/sl"
	"api/internal/storage/postgresql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Start app", slog.String("env", cfg.Env))

	psql, err := postgresql.New(cfg.Db)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/courses", func(r chi.Router) {
		r.Post("/", save.New(log, psql))
		r.Get("/", get.GetAll(log, psql))
		r.Get("/{id}", get.GetByID(log, psql))
		// r.Put("/{id}", )
		// r.Delete("/{id}", )
	})

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Info("starting http server", slog.String("address", cfg.Address))

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("failed to start server")
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
