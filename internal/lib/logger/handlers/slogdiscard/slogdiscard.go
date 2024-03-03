// internal/lib/logger/handlers/slogdiscard/slogdiscard.go

// DiscardHandler. Это пакет расширения логгера.
// В таком виде логгер будет игнорировать все сообщения, которые мы в него отправляем,
// это понадобится в тестах.
// имплементируем в нем интерфейс slog.Handler

package slogdiscard

import (
	"context"

	"log/slog"
)

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())

}

// DiscardHandler. Это пакет расширения логгера для тестов
type DiscardHandler struct{}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	//просто игнорируем запись журнала
	return nil
}

func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	// Возвращает тот же обработчик, так как нет атрибутов для сохранения
	return h
}

func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	// Возвращает тот же обработчик, так как нет группы для сохранения
	return h
}

func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	// Всегда возвращает false, так как запись журнала игнорируется
	return false
}
