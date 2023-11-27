package cmd

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
}

func RunService() {
	// Read application environment and load configurations
	conf, err := config.NewConfig(env)
	if err != nil {
		panic(err)
	}
	// Initialise console and file logger
	consoleFileLogger := logger.NewLogger()
	// Initialize the logging service
	loggerSvc := services.NewLoggerService(&consoleFileLogger)
	// // Postgres Clien
	databaseRepo, err := postgres.NewPostgresClient(*conf, loggerSvc)
	if err != nil {
		loggerSvc.Error(err.Error())
		panic(err)
	}
	// Initialize the article service
	articleService := services.NewArticleManagementService(databaseRepo)
	// Run HTTP Server
	app.InitGinRoutes(articleService, loggerSvc, *conf)
	loggerSvc.Close()
}
