package main

import (
	"api_server/internal/config"
	"fmt"
	"log/slog"
	"os"
)

const (
    envLocal = "local"
    envDev   = "dev"
    envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
    var log *slog.Logger

    switch env {
    case envLocal:
        log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    case envDev:
        log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    case envProd:
        log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
    }

    return log
}

func main() {

	/* cleanenv */
	cfg := config.LoadConfig()
	// fmt.Println(cfg)

	/* slog */
	log := setupLogger(cfg.Env)
    log = log.With(slog.String("env", cfg.Env)) // к каждому сообщению будет добавляться поле с информацией о текущем окружении

    log.Info("Initializing server", slog.String("address", cfg.Address)) // Помимо сообщения выведем параметр с адресом
    log.Debug("Logger debug mode enabled")

	// storage

	// chi

	// test

}