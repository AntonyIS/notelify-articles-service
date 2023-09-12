package main

import (
	"flag"

	"github.com/AntonyIS/notelify-articles-service/config"
	"github.com/AntonyIS/notelify-articles-service/internal/adapters/app"
	"github.com/AntonyIS/notelify-articles-service/internal/adapters/logger"
	"github.com/AntonyIS/notelify-articles-service/internal/adapters/repository/postgres"
	"github.com/AntonyIS/notelify-articles-service/internal/core/services"
)

var env string

func init() {
	flag.StringVar(&env, "env", "dev", "The environment the application is running")
	flag.Parse()
	// logger.SetupLogger()
}

func main() {
	/*
	This is the microservice article service
	
	*/
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

		contentSVC := services.NewArticleManagementService(postgresDBRepo)
		app.InitGinRoutes(contentSVC, logger, *conf)
	}
	logger.PostLogMessage(err.Error())

}
