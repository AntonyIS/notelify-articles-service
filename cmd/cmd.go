package cmd

import (
	"github.com/AntonyIS/notelify-articles-service/config"
	"github.com/AntonyIS/notelify-articles-service/internal/adapters/app"
	"github.com/AntonyIS/notelify-articles-service/internal/adapters/repository/postgres"
	"github.com/AntonyIS/notelify-articles-service/internal/core/domain"
	"github.com/AntonyIS/notelify-articles-service/internal/core/services"
)

func RunService() {
	conf, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	newLoggerService := services.NewLoggingManagementService(conf.LOGGER_URL)

	databaseRepo, err := postgres.NewPostgresClient(*conf)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "INFOR",
			Service:  "users",
			Message:  err.Error(),
		}
		newLoggerService.LogInfo(logEntry)
		panic(err)
	}

	articleService := services.NewArticleManagementService(databaseRepo, newLoggerService)
	app.InitGinRoutes(articleService, newLoggerService, *conf)

}
