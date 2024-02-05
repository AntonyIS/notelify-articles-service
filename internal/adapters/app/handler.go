package app

import (
	"fmt"
	"log"
	"time"

	appConfig "github.com/AntonyIS/notelify-articles-service/config"
	"github.com/AntonyIS/notelify-articles-service/internal/core/domain"
	"github.com/AntonyIS/notelify-articles-service/internal/core/ports"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitGinRoutes(svc ports.ArticleService, logger ports.LoggingService, conf appConfig.Config) {
	gin.SetMode(gin.DebugMode)

	router := gin.Default()
	router.Use(ginRequestLogger(logger))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	handler := NewGinHandler(svc, conf.SECRET_KEY, logger)

	articleRoutes := router.Group("/posts/v1")
	{
		articleRoutes.GET("/healthcheck", handler.HealthCheck)
		articleRoutes.POST("/", handler.CreateArticle)
		articleRoutes.GET("/:post_id", handler.GetArticleByID)
		articleRoutes.GET("/", handler.GetArticles)
		articleRoutes.GET("/author/:author_id", handler.GetArticlesByAuthor)
		articleRoutes.GET("/tag/:tag_name", handler.GetArticlesByTag)
		articleRoutes.PUT("/:post_id", handler.UpdateArticle)
		articleRoutes.DELETE("/:post_id", handler.DeleteArticle)
		articleRoutes.DELETE("/", handler.DeleteArticleAll)
	}
	logEntry := domain.LogMessage{
		LogLevel: "INFO",
		Service:  "articles",
		Message:  fmt.Sprintf("Server running on port 0.0.0.0:%s", conf.SERVER_PORT),
	}
	logger.LogError(logEntry)

	log.Printf("Server running on port 0.0.0.0:%s", conf.SERVER_PORT)
	router.Run(fmt.Sprintf(":%s", conf.SERVER_PORT))
}

func ginRequestLogger(logger ports.LoggingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		logEntry := domain.LogMessage{
			LogLevel: "INFO",
			Service:  "articles",
			Message: fmt.Sprintf("%s %s %s %d %s %s",
				c.Request.Method,
				c.Request.URL.Path,
				c.Request.Proto,
				c.Writer.Status(),
				latency.String(),
				c.ClientIP(),
			),
		}
		logger.LogError(logEntry)
		// logger.Info(fmt.Sprintf("%s %s %s %d %s %s",
		// 	c.Request.Method,
		// 	c.Request.URL.Path,
		// 	c.Request.Proto,
		// 	c.Writer.Status(),
		// 	latency.String(),
		// 	c.ClientIP(),
		// ))
	}
}
