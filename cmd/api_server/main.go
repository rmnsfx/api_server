package main

import (
	"api_server/internal/config"
	"api_server/internal/storage/sqlite"
	"net/http"

	// "fmt"
	"log/slog"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
	log = log.With(slog.String("env", cfg.Env)) // К каждому сообщению будет добавляться поле с информацией о текущем окружении

	log.Info("Initializing server", slog.String("address", cfg.Address)) // Помимо сообщения выведем параметр с адресом
	log.Debug("Logger debug mode enabled")

	/* storage */
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to initialize storage", err)
	}
	
    storage.SaveGameLaunch("device")

	/* middleware */
    router := chi.NewRouter()  
  
    router.Use(middleware.RequestID) // Добавляет request_id в каждый запрос, для трейсинга
    router.Use(middleware.Logger) // Логирование всех запросов
    router.Use(middleware.Recoverer)  // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
    router.Use(middleware.URLFormat) // Парсер URLов поступающих запросов

    /* http server (chi) */
    log.Info("Server start... ", slog.String("address", cfg.Address) )

    srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

    if err := srv.ListenAndServe(); err != nil {
        log.Error("Failed to start server")
    }

    log.Info("Stopping server")


	// test

}
