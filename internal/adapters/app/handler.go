package app

import (
	"fmt"

	appConfig "github.com/AntonyIS/notlify-content-svc/config"
	"github.com/AntonyIS/notlify-content-svc/internal/adapters/logger"
	"github.com/AntonyIS/notlify-content-svc/internal/core/ports"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitGinRoutes(svc ports.ContentService, logger logger.LoggerType, conf appConfig.Config) {
	gin.SetMode(gin.DebugMode)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	handler := NewGinHandler(svc, conf.SECRET_KEY)

	contentsRoutes := router.Group("/v1/contents")

	{
		contentsRoutes.GET("/", handler.ReadContents)
		contentsRoutes.GET("/:id", handler.ReadContent)
		contentsRoutes.PUT("/:id", handler.UpdateContent)
		contentsRoutes.DELETE("/:id", handler.DeleteContent)
		contentsRoutes.POST("/", handler.CreateContent)
		contentsRoutes.DELETE("/delete/all", handler.DeleteAllContent)
		contentsRoutes.GET("/users/:creator_id", handler.ReadCreatorContents)
	}
	logger.PostLogMessage(fmt.Sprintf("Server running on port :%s", conf.Port))
	router.Run(fmt.Sprintf(":%s", conf.Port))
}
