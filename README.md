# REST API сервис — URL Shortener

## Используюемые библиотеки:

- go-chi/chi              — для обработки HTTP-запросов,
- slog                    — для логирования,
- stretchr/testify        — для покрытия проекта тестами,
- ilyakaznacheev/cleanenv — для конфигурирования,
- SQLite                  — для хранения данных, СУБД.


## ВНИМАНИЕ!!!
Для успешной компиляции sqlite3 должен быть установлен компилятор gcc и установлен флаг:
```bash
go env -w CGO_ENABLED=1
```

Компилятор gcc по адресу:
http://www.equation.com/servlet/equation.cmd?fa=fortran

ЕСТЬ АЛЬТЕРНАТИВА SQLITE3:
https://gitlab.com/cznic/sqlite
Для его установки:
```bash
go get modernc.org/sqlite
```

## УСТАНОВКА ПАКЕТОВ:
ilyakaznacheev/cleanenv (для конфигурирования):
```bash
go get github.com/ilyakaznacheev/cleanenv
```

SQLite:
```bash
go get github.com/mattn/go-sqlite3
```


Маршрутизатор go-chi/chi:
```bash
go get -u github.com/go-chi/chi/v5
go get github.com/go-chi/render
go get github.com/go-chi/cors
```

testify (моки):
```bash
go get github.com/stretchr/testify
go get github.com/stretchr/testify/mock
go get github.com/stretchr/testify/require
```

Библиотеки, которые очень упрощают написание тестов:
httpexpect — для тестирования REST API,
```bash
go get github.com/brianvoe/gofakeit/v6
```

gofakeit — для генерации случайных данных разного формата (имена, имейлы, номера телефонов, URL и другое).
```bash
go get github.com/gavv/httpexpect/v2
```
-----------------------------------------------------------------------------------------



## ЗАПУСК СЕРВИСА:
```bash
go run ./cmd/url-shortener/main.go --config=./config/local.yaml
```

Для запуска (в bash, с использованием переменной окружения):
```bash
CONFIG_PATH="./config/local.yaml" go run  "./cmd/url-shortener/main.go"
```

ЗАПУСК ТЕСТОВ:
```bash
go test ./tests -count=1 -v
```
-----------------------------------------------------------------------------------------
Пример POST-запроса:
```http request
localhost:8082/url
```
тело (body) JSON:
```json
{
  "URL": "https://ya.ru",
  "Alias": "ya"
}
```
Пример GET-запроса:
```http request
localhost:8082/ViSq4r
```

-----------------------------------------------------------------------------------------
## ПРИМЕР РУЧНОЙ УСТАНОВКИ ТЕГА
```bash
git tag v0.0.8 && git push origin v0.0.8
```


