// cmd/url-shortener/main.go

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/url/redirect"
	"url-shortener/internal/http-server/handlers/url/remove"

	"url-shortener/internal/http-server/handlers/url/save"
	mwLogger "url-shortener/internal/http-server/middleware/logger"

	//"url-shortener/internal/lib/logger/handlers/slogpretty"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"
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
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	router.Use(middleware.URLFormat) // Парсер URLов поступающих запросов

	// По умолчанию middleware.Logger использует свой собственный внутренний логгер,
	// который желательно переопределить, чтобы использовался наш,
	// иначе могут возникнуть проблемы — например, со сбором логов.
	// Либо можно написать собственный middleware для логирования запросов. Так и сделаем

	// РАЗОБРАТЬСЯ!!!!
	//--------------------------------------------------------------------------------
	router.Post("/", save.New(log, storage))

	// router.Route("/url", func(r chi.Router) {
	// 	r.Use(middleware.BasicAuth("url-shortener", map[string]string{
	// 		cfg.HTTPServer.User: cfg.HTTPServer.Password,
	// 	}))

	// 	r.Post("/", save.New(log, storage))
	// 	// TODO: add DELETE /url/{id}
	// })

	// Подключаем редирект-хендлер.
	// Здесь формируем путь для обращения и именуем его параметр — {alias}.
	// В хендлере можно получить этот параметр по указанному имени
	router.Get("/{alias}", redirect.New(log, storage))
	// Это очень удобная и гибкая штука. Вы можете формировать и более сложные пути, например:
	//// router.Get("/v1/{user_id}/uid", redirect.New(log, storage))

	//прикручиваем ремувер
	router.Delete("/{alias}", remove.New(log, storage))

	// ЗАПУСК СЕРВЕРА
	log.Info("starting server", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	// ждем, пока в канал не придет сигнал с остановкой сервера
	<-done
	log.Info("stopping server")

	// TODO: move timeout to config
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	// TODO: close storage
	log.Info("server stopped")
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
