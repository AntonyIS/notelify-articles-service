package app

import (
	"fmt"
	"time"

	appConfig "github.com/AntonyIS/notelify-articles-service/config"
	"github.com/AntonyIS/notelify-articles-service/internal/core/ports"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitGinRoutes(svc ports.ArticleService, logger ports.Logger, conf appConfig.Config) {
	gin.SetMode(gin.DebugMode)

	router := gin.Default()
	router.Use(ginRequestLogger(logger))
	if conf.Env == "prod" {
		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://notelify-client-service:3000", "http://localhost:3000"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		}))

	} else {
		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		}))

	}

	handler := NewGinHandler(svc, conf.SECRET_KEY)

	articleRoutes := router.Group("/v1/articles")
	{
		articleRoutes.POST("/", handler.CreateArticle)
		articleRoutes.GET("/:article_id", handler.GetArticleByID)
		articleRoutes.GET("/", handler.GetArticles)
		articleRoutes.GET("/author/:author_id", handler.GetArticlesByAuthor)
		articleRoutes.GET("/tag/:tag_name", handler.GetArticlesByTag)
		articleRoutes.PUT("/:article_id", handler.UpdateArticle)
		articleRoutes.DELETE("/:article_id", handler.DeleteArticle)
		articleRoutes.DELETE("/", handler.DeleteArticleAll)
	}

	logger.Info(fmt.Sprintf("Server running on port :%s", conf.Port))
	router.Run(fmt.Sprintf(":%s", conf.Port))
}

func ginRequestLogger(logger ports.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		logger.Info(fmt.Sprintf("%s %s %s %d %s %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.Proto,
			c.Writer.Status(),
			latency.String(),
			c.ClientIP(),
		))
	}
}
