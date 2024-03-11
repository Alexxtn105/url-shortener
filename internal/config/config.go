// internal/config/config.go
package config

// Установка cleanenv:
// go get -u github.com/ilyakaznacheev/cleanenv

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config Структура конфигурации
// Здесь используются cледующие struct-теги (для анмаршалинга):
// yaml — имя соответствующего параметра в Yaml-файле,
// env-default — дефолтное значение,
// env-required — делает параметры обязательными. Если такой параметр не указан, мы будем получать ошибку.
type Config struct {
	Env         string `yaml:"env" env-default:"development"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
	Clients     ClientConfig `yaml:"clients"`
	appSecret   string       `yaml:"app_secret" env-required:"true" env:"APP_SECRET"` // секретный ключ, с помощью которого приложение будет проверять JWT-токены
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

type Client struct {
	Address      string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retriesCount"`
	//Insecure bool `yaml:"insecure"`  //может пригоидиться в будущем. Можно запускать приложение в двух режимах
}

type ClientConfig struct {
	SSO Client `yaml:"sso"`
}

// MustLoad формирует стурктуру Config из переменной окружения CONFIG_PATH
// ПРИМЕЧАНИЕ: Приставка Must в имени функции обычно говорит,
// что функция вместо возврата ошибки аварийно завершает работу приложения
// — например, будет паниковать. Таким подходом злоупотреблять не стоит,
// но иногда это бывает удобно. Например, если ваше приложение при запуске упадет
// с паникой из-за кривого или отсутствующего конфиг-файла, это нормально.
// А вот в бизнес-логике такого лучше не допускать
// Также обращаю внимание, что путь до конфиг-файла я получаю из переменной окружения CONFIG_PATH,
// дефолтный путь не предусмотрен. Чтобы передать значение такой переменной,
// можно запустить приложение следующей командой:
// CONFIG_PATH=./config/local.yaml ./your-app
func MustLoad() *Config {
	// получаем путь до конфиг-файла из ENV-переменной CONFIG_PATH
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	//проверяем существование конфиг-файла
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}

	//читаем конфиг файл в структуру
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}

func MustLoadFetchFlag() *Config {
	// получаем путь до конфиг-файла из ENV-переменной CONFIG_PATH
	configPath := fetchConfigPath()

	if configPath == "" {
		log.Fatal("config path is empty")
	}

	//проверяем существование конфиг-файла
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}

	//читаем конфиг файл в структуру
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
