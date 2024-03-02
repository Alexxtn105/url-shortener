// cmd/url-shortener/main.go

package main

import (
	"fmt"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// для запуска (в bash):
// CONFIG_PATH="./config/local.yaml" go run  "./cmd/url-shortener/main.go"
func main() {
	//получаем объект конфига
	//cfg := config.MustLoad()			// с использованием переменной окружения CONFIG_PATH
	cfg := config.MustLoadFetchFlag() // ...или с использованием параметра командной строки
	//	fmt.Printf("Конфигурация загружена успешно: %s\n", cfg)
	fmt.Println("Конфигурация загружена успешно")

	// Получаем логгер.
	// Будем использовать slog, потому что это очень гибкий пакет,
	// и конкретная реализация может быть разной.
	// Мы можем написать собственный хендлер (обработчик логов, который определяет, что происходит с записями),
	// обернуть в него привычный логгер (например, zap или logrus)
	// либо использовать дефолтные варианты, которые предоставляются вместе с пакетом.
	// Из коробки в slog есть два вида хендлеров.
	// Для локальной разработки нам подойдет TextHandler,
	// а для деплоя лучше использовать JSONHandler,
	// чтобы агрегатор логов (Kibana, Grafana, Loki и другие) мог его распарсить.
	// Кроме того, важно учесть уровень логирования —
	// это минимальный уровень сообщений, которые будут выводиться.
	// К примеру, если мы установим уровень Info, то Debug-сообщения не увидим.
	// Поэтому для локальной разработки и Dev-окружения лучше использовать уровень Debug,
	// а для продакшена — Info.

	//создаем логгер
	log := setupLogger(cfg.Env)
	//добавим параметр env с помощью метода log.With
	log = log.With(slog.String("env", cfg.Env)) // к каждому сообщению будет добавляться поле с информацией о текущем окружении

	log.Info("initializing server", slog.String("address", cfg.Address)) // Помимо сообщения выведем параметр с адресом
	log.Debug("logger debug mode enabled")

	//создаем объект Storage
	storage, err := sqlite.NewStorage(cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
	}

	log.Info("storage created")
	fmt.Println(storage)

	//----------------------создаем http-сервер

	//создаем объект роутера
	router := chi.NewRouter()

	router.Use(middleware.RequestID) // Добавляет request_id в каждый запрос, для трейсинга
	router.Use(middleware.Logger)    // Логирование всех запросов. Желательно написать собственный
	router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	router.Use(middleware.URLFormat) // Парсер URLов поступающих запросов

	// По умолчанию middleware.Logger использует свой собственный внутренний логгер,
	// который желательно переопределить, чтобы использовался наш,
	// иначе могут возникнуть проблемы — например, со сбором логов.
	// Либо можно написать собственный middleware для логирования запросов. Так и сделаем

}

// setupLogger создает логгер в зависимости от окружения с разными параметрами — TextHandler / JSONHandler и уровень LevelDebug / LevelInfo
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
