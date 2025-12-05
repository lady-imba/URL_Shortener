package main

import (
	"URL_SHORTENER/internal/config"
	"URL_SHORTENER/internal/http-server/handlers/redirect"
	"URL_SHORTENER/internal/http-server/handlers/url/delete"
	"URL_SHORTENER/internal/http-server/handlers/url/save"
	"URL_SHORTENER/internal/http-server/middleware/logger"
	"URL_SHORTENER/internal/lib/logger/sl"
	"URL_SHORTENER/internal/storage/sqlite"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)

func main(){
	config := config.MustLoad()

	log := setupLogger(config.Env)

	log.Info("Starting url-shrtener", slog.String("env", config.Env))

	storage, err := sqlite.NewStorage(config.StoragePath)
	if err != nil {
		log.Error("Failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router){
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			config.HTTPServer.User: config.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", delete.New(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", config.Address))

	server := &http.Server{
		Addr: config.Address,
		Handler: router,
		ReadTimeout: config.HTTPServer.Timeout,
		WriteTimeout: config.HTTPServer.Timeout,
		IdleTimeout: config.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("Failed to start server", sl.Err(err))
	}

	log.Error("server stopped")
}

func setupLogger(env string)*slog.Logger{
	var log *slog.Logger

	switch env{
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout,&slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout,&slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout,&slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log 
}