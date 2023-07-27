package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env                string
	Port               string
	UserTable          string
	AWS_ACCESS_KEY     string
	AWS_SECRET_KEY     string
	AWS_DEFAULT_REGION string
	LoggerURL          string
	SECRET_KEY         string
	DatabaseName       string
	DatabaseUser       string
	DatabaseHost       string
	DatabasePort       string
	DatabasePassword   string
	Debugging          bool
	Testing            bool
}

func NewConfig(Env string) (*Config, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return nil, err
	}
	var (
		AWS_ACCESS_KEY     = os.Getenv("AWS_ACCESS_KEY")
		AWS_SECRET_KEY     = os.Getenv("AWS_SECRET_KEY")
		AWS_DEFAULT_REGION = os.Getenv("AWS_DEFAULT_REGION")
		LoggerURL          = os.Getenv("LoggerURL")
		SECRET_KEY         = os.Getenv("SECRET_KEY")
		Port               = "8080"
		UserTable          = "UsersTable"
		DatabaseName       = "Notlify"
		DatabaseUser       = os.Getenv("DatabaseUser")
		DatabasePort       = "5432"
		DatabaseHost       = ""
		DatabasePassword   = os.Getenv("DatabasePassword")
		Testing            = false
		Debugging          = false
	)

	switch Env {
	case "testing":
		Testing = true
		Debugging = true

	case "dev":
		Testing = true
		Debugging = true
		DatabaseHost = "localhost"

	case "prod":
		Testing = false
		Debugging = false
		DatabaseHost = os.Getenv("DatabaseHost")
	}

	config := Config{
		Env:                Env,
		Port:               Port,
		UserTable:          UserTable,
		AWS_ACCESS_KEY:     AWS_ACCESS_KEY,
		AWS_SECRET_KEY:     AWS_SECRET_KEY,
		SECRET_KEY:         SECRET_KEY,
		AWS_DEFAULT_REGION: AWS_DEFAULT_REGION,
		LoggerURL:          LoggerURL,
		Debugging:          Debugging,
		Testing:            Testing,
		DatabaseName:       DatabaseName,
		DatabaseUser:       DatabaseUser,
		DatabaseHost:       DatabaseHost,
		DatabasePort:       DatabasePort,
		DatabasePassword:   DatabasePassword,
	}

	return &config, nil
}
