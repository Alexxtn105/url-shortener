REST API сервис — URL Shortener
-----------------------------------------------------------------------------------------
Для проекта используются библиотеки:

go-chi/chi              — для обработки HTTP-запросов,
slog                    — для логирования,
stretchr/testify        — для покрытия проекта тестами,
ilyakaznacheev/cleanenv — для конфигурирования,
SQLite                  — для хранения данных, СУБД.


-----------------------------------------------------------------------------------------
ВНИМАНИЕ!!!
Для успешной компиляции sqlite3 должен быть установлен компилятор gcc и установлен флаг:
go env -w CGO_ENABLED=1

Компилятор gcc по адресу:
http://www.equation.com/servlet/equation.cmd?fa=fortran


-----------------------------------------------------------------------------------------
УСТАНОВКА ПАКЕТОВ:
1. ilyakaznacheev/cleanenv (для конфигурирования):
go get github.com/ilyakaznacheev/cleanenv

2. SQLite:
go get github.com/mattn/go-sqlite3

3. go-chi/chi (для работы с HTTP-сервером):
go get -u github.com/go-chi/chi/v5
go get github.com/go-chi/render

4. testify:
go get github.com/stretchr/testify
go get github.com/stretchr/testify/mock
go get github.com/stretchr/testify/require

5. Библиотеки, которые очень упрощают написание тестов:
5.1 httpexpect — для тестирования REST API,
go get github.com/brianvoe/gofakeit/v6

5.2 gofakeit — для генерации случайных данных разного формата (имена, имейлы, номера телефонов, URL и другое).
go get github.com/gavv/httpexpect/v2

ЗАПУСК СЕРВИСА:
go run ./cmd/url-shortener/main.go --config=./config/local.yaml

ЗАПУСК ТЕСТОВ:
go test ./tests -count=1 -v

Пример POST-запроса
localhost:8082/?url=https://ya.ru&alias=yaru