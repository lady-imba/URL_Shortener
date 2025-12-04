package main

import (
	"URL_SHORTENER/internal/config"
	"URL_SHORTENER/internal/lib/logger/sl"
	"URL_SHORTENER/internal/storage/sqlite"
	"log/slog"
	"os"
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

	_ = storage
	
	log.Info("Storage init")

	//TODO: init storage

	//TODO: init router

	//TODO: run server
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