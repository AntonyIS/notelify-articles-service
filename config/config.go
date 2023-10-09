package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env                   string
	Port                  string
	ContentTable          string
	AWS_ACCESS_KEY        string
	AWS_SECRET_KEY        string
	AWS_DEFAULT_REGION    string
	RDSInstanceIdentifier string
	LoggerURL             string
	SECRET_KEY            string
	DatabaseName          string
	DatabaseUser          string
	DatabaseHost          string
	DatabasePort          string
	DatabasePassword      string
	Debugging             bool
	Testing               bool
}

func LoadEnv() error {
	err := godotenv.Load(".env")

	if err != nil {
		return err
	}
	return nil
}

func NewConfig(Env string) (*Config, error) {
	if Env == "dev" {
		err := godotenv.Load(".env")
		if err != nil {
			return nil, err
		}
	} else if Env == "test" {
		err := godotenv.Load("../../../.env")
		if err != nil {
			return nil, err
		}
	}

	var (
		AWS_ACCESS_KEY        = os.Getenv("AWS_ACCESS_KEY")
		AWS_SECRET_KEY        = os.Getenv("AWS_SECRET_KEY")
		AWS_DEFAULT_REGION    = os.Getenv("AWS_DEFAULT_REGION")
		RDSInstanceIdentifier = os.Getenv("RDSInstanceIdentifier")
		SECRET_KEY            = os.Getenv("SECRET_KEY")
		LoggerURL             = os.Getenv("LoggerURL")
		DatabaseUser          = os.Getenv("DatabaseUser")
		DatabasePassword      = os.Getenv("DatabasePassword")
		Port                  = "8001"
		ContentTable          = "Articles"
		DatabaseName          = "postgres"
		DatabasePort          = "5432"
		DatabaseHost          = "localhost"
		Testing               = false
		Debugging             = false
	)

	switch Env {
	case "prod":
		Testing = false
		Debugging = false
		DatabaseHost = os.Getenv("DatabaseHost")
		DatabaseName = "notlify_db_init"

	case "test_prod":
		Testing = true
		Debugging = true
		DatabaseHost = os.Getenv("DatabaseHost")
		DatabaseName = "notlify_db_init"
		ContentTable = "TestArticles"

	case "test":
		Testing = true
		Debugging = true
		ContentTable = "TestArticles"

	case "dev":
		Testing = true
		Debugging = true

	}

	config := Config{
		Env:                   Env,
		Port:                  Port,
		ContentTable:          ContentTable,
		AWS_ACCESS_KEY:        AWS_ACCESS_KEY,
		AWS_SECRET_KEY:        AWS_SECRET_KEY,
		RDSInstanceIdentifier: RDSInstanceIdentifier,
		SECRET_KEY:            SECRET_KEY,
		AWS_DEFAULT_REGION:    AWS_DEFAULT_REGION,
		LoggerURL:             LoggerURL,
		Debugging:             Debugging,
		Testing:               Testing,
		DatabaseName:          DatabaseName,
		DatabaseUser:          DatabaseUser,
		DatabaseHost:          DatabaseHost,
		DatabasePort:          DatabasePort,
		DatabasePassword:      DatabasePassword,
	}

	return &config, nil
}
