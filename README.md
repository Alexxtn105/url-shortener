REST API сервис — URL Shortener

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

Компилятор по адресу:
http://www.equation.com/servlet/equation.cmd?fa=fortran


-----------------------------------------------------------------------------------------
УСТАНОВКА ПАКЕТОВ:
ilyakaznacheev/cleanenv (для конфигурирования):
go get github.com/ilyakaznacheev/cleanenv

SQLite:
go get github.com/mattn/go-sqlite3

go-chi/chi (для работы с HTTP-сервером):
go get -u github.com/go-chi/chi/v5
go get github.com/go-chi/render