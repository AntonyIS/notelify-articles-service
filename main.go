package main

import (
	"flag"

	"github.com/AntonyIS/notlify-content-svc/config"
	"github.com/AntonyIS/notlify-content-svc/internal/adapters/app"
	"github.com/AntonyIS/notlify-content-svc/internal/adapters/logger"
	"github.com/AntonyIS/notlify-content-svc/internal/adapters/repository/postgres"
	"github.com/AntonyIS/notlify-content-svc/internal/core/services"
)

var env string

func init() {
	flag.StringVar(&env, "env", "dev", "The environment the application is running")
	flag.Parse()
	// logger.SetupLogger()
}

func main() {
	conf, err := config.NewConfig(env)
	if err != nil {
		panic(err)
	}
	// Logger service
	logger := logger.NewLoggerService(conf.LoggerURL)
	// // Postgres Client
	postgresDBRepo, err := postgres.NewPostgresClient(*conf, logger)
	if err != nil {
		logger.PostLogMessage(err.Error())
		panic(err)
	} else {
		// // User service
		contentSVC := services.NewContentManagementService(postgresDBRepo)
		// // Initialize HTTP server
		app.InitGinRoutes(contentSVC, logger, *conf)
	}
	logger.PostLogMessage(err.Error())

}
