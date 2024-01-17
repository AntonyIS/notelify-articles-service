package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ENV               string
	SERVER_PORT       string
	ARTICLE_TABLE     string
	LOGGER_URL        string
	SECRET_KEY        string
	POSTGRES_DB       string
	POSTGRES_USER     string
	POSTGRES_HOST     string
	POSTGRES_PORT     string
	POSTGRES_PASSWORD string
	DEBUG             bool
	TEST              bool
}

func NewConfig() (*Config, error) {
	ENV := os.Getenv("ENV")

	if ENV == "test" {
		err := godotenv.Load("../../../.env")
		if err != nil {
			return nil, err
		}
	} else {
		err := godotenv.Load(".env")
		if err != nil {
			return nil, err
		}
		ENV = os.Getenv("ENV")
	}

	var (
		SECRET_KEY        = os.Getenv("SECRET_KEY")
		POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
		POSTGRES_USER     = "postgres"
		POSTGRES_DB       = "postgres"
		POSTGRES_HOST     = "postgres"
		POSTGRES_PORT     = "5432"
		SERVER_PORT       = "8001"
		ARTICLE_TABLE     = "Articles"
		LOGGER_URL        = "http://localhost:8003/v1/logger"
		DEBUG             = false
		TEST              = false
	)

	switch ENV {
	case "production":
		TEST = false
		DEBUG = false

	case "production_test":
		TEST = true
		DEBUG = true
		ARTICLE_TABLE = "TestArticles"

	case "developement_test":
		TEST = true
		DEBUG = true
		ARTICLE_TABLE = "TestArticles"

	case "development":
		TEST = true
		DEBUG = true
		POSTGRES_HOST = "localhost"
		ARTICLE_TABLE = "DevArticles"

	case "docker":
		TEST = true
		DEBUG = true
		ARTICLE_TABLE = "DockerArticles"
		LOGGER_URL = "/v1/logger/"

	case "docker_test":
		TEST = true
		DEBUG = true
		ARTICLE_TABLE = "DockerArticles"
	}

	config := Config{
		ENV:               ENV,
		SERVER_PORT:       SERVER_PORT,
		ARTICLE_TABLE:     ARTICLE_TABLE,
		SECRET_KEY:        SECRET_KEY,
		LOGGER_URL:        LOGGER_URL,
		DEBUG:             DEBUG,
		TEST:              TEST,
		POSTGRES_DB:       POSTGRES_DB,
		POSTGRES_USER:     POSTGRES_USER,
		POSTGRES_HOST:     POSTGRES_HOST,
		POSTGRES_PORT:     POSTGRES_PORT,
		POSTGRES_PASSWORD: POSTGRES_PASSWORD,
	}

	return &config, nil
}
