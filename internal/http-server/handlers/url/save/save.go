// internal/http-server/handlers/url/save/save.go
package save

import (
	"errors"
	"io"
	"net/http"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	resp "url-shortener/internal/lib/api/response" // для краткости даем короткий алиас пакету
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"
)

// структура запроса
type Request struct {
	URL   string `json:"url" validate:"required,url"` // эта строчка для валидации, об этом будет ниже.
	Alias string `json:"alias,omitempty"`
}

// структура ответа
type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// TODO: move to config when needed
const aliasLength = 6

// интерфейс сохранения полученной URL-строки
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver
type URLSaver interface {
	SaveURL(URL, alias string) (int64, error)
}

// Тесты:
// Mockery generation fo SaveURL:
// ./internal/http-server/handlers/url/save/save.go

// New Конструктор обработчика запросов
func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		// Добавляем к текущему объекту логгера поля op и request_id
		// Они могут очень упростить нам жизнь в будущем
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Создаем объект запроса и анмаршаллим в него запрос
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// Такую ошибку встретим, если получили запрос с пустым телом
			// Обработаем её отдельно
			log.Error("request body is empty")

			//	render.JSON(w, r, resp.Response{
			//		Status: resp.StatusError,
			//		Error:  "empty request",
			//	})
			//переписал так:
			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			// render.JSON(w, r, resp.Response{
			// 	Status: resp.StatusError,
			// 	Error:  "failed to decode request",
			// })
			render.JSON(w, r, resp.Error("empty request"))

			return
		}

		// Лучше больше логов, чем меньше - лишнее мы легко сможем почистить,
		// при необходимости. А вот недостающую информацию мы уже не получим.
		log.Info("request body decoded", slog.Any("req", req))

		// Создаем объект валидатора
		// и передаем в него структуру, которую нужно провалидировать
		if err := validator.New().Struct(req); err != nil {
			// Приводим ошибку к типу ошибки валидации
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.Error(validateErr.Error()))

			return
		}

		// Alias проверяем вручную. Если он пустой — генерируем случайный:
		alias := req.Alias
		if alias == "" {
			// используем собственный генератор случайных строк
			alias = random.NewRandomString(aliasLength)
		}

		// Осталось только сохранить URL и Alias,
		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			// отдельно обрабатываем ситуацию, когда запись с таким alias уже существует
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}

		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		// а после — вернуть ответ с сообщением об успехе.
		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
