package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/http-server/handlers/url/save/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string // Имя теста
		alias     string // Отправляемый alias
		url       string // Отправляемый URL
		respError string // Какую ошибку мы должны получить?
		mockError error  // Ошибку, которую вернёт мок
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
			// Тут поля respError и mockError оставляем пустыми,
			// т.к. это успешный запрос
		},
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://google.com",
		},
		{
			name:      "Empty URL",
			url:       "",
			alias:     "some_alias",
			respError: "field URL is a required field",
		},
		{
			name:      "Invalid URL",
			url:       "some invalid URL",
			alias:     "some_alias",
			respError: "field URL is not a valid URL",
		},
		{
			name:      "SaveURL Error",
			alias:     "test_alias",
			url:       "https://google.com",
			respError: "failed to add url",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Создаем объект мока стораджа
			urlSaverMock := mocks.NewURLSaver(t)

			// Если ожидается успешный ответ, значит к моку точно будет вызов
			// Либо даже если в ответе ожидаем ошибку,
			// но мок должен ответить с ошибкой, к нему тоже будет запрос:
			if tc.respError == "" || tc.mockError != nil {
				// Сообщаем моку, какой к нему будет запрос, и что надо вернуть
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockError).
					Once() // Запрос будет ровно один
			}
			// Создаем наш хэндлер
			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			// Формируем тело запроса
			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			// Создаем ResponseRecorder для записи ответа хэндлера
			rr := httptest.NewRecorder()

			// Обрабатываем запрос, записывая ответ в рекордер
			handler.ServeHTTP(rr, req)

			// Проверяем, что статус ответа корректный
			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp save.Response

			// Анмаршаллим тело, и проверяем что при этом не возникло ошибок
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			// Проверяем наличие требуемой ошибки в ответе
			require.Equal(t, tc.respError, resp.Error)

			// Другие проверки
		})
	}
}
