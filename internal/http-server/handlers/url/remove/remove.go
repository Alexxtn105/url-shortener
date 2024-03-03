// internal/http-server/handlers/remove

package remove

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// URLRemover is an interface for removing url by alias.
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLRemover
type URLRemover interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlRemover URLRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.remove.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Роутер chi позволяет делать вот такие финты -
		// получать GET-параметры по их именам.
		// Имена определяются при добавлении хэндлера в роутер, это будет ниже.
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("not found"))
			return
		}

		// Находим URL по алиасу в БД
		err := urlRemover.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			// Не нашли URL, сообщаем об этом клиенту
			log.Info("url not found", "alias", alias)
			render.JSON(w, r, resp.Error("not found"))
			return
		}
		if err != nil {
			// Не удалось осуществить поиск
			log.Error("failed to delete url", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("delete url by alias", slog.String("alias", alias))

	}

}
