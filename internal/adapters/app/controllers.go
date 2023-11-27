package app

import (
	"net/http"

	"github.com/AntonyIS/notelify-articles-service/internal/core/domain"
	"github.com/AntonyIS/notelify-articles-service/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type GinHandler interface {
	CreateArticle(ctx *gin.Context)
	GetArticleByID(ctx *gin.Context)
	GetArticles(ctx *gin.Context)
	GetArticlesByAuthor(ctx *gin.Context)
	GetArticlesByTag(ctx *gin.Context)
	UpdateArticle(ctx *gin.Context)
	DeleteArticle(ctx *gin.Context)
	DeleteArticleAll(ctx *gin.Context)
}

type handler struct {
	svc       ports.ArticleService
	secretKey string
}

func NewGinHandler(svc ports.ArticleService, secretKey string) GinHandler {
	routerHandler := handler{
		svc:       svc,
		secretKey: secretKey,
	}
	return routerHandler
}

func (h handler) CreateArticle(ctx *gin.Context) {
	var res *domain.Article
	if err := ctx.ShouldBindJSON(&res); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response, err := h.svc.CreateArticle(res)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, response)
}

func (h handler) GetArticleByID(ctx *gin.Context) {
	id := ctx.Param("article_id")
	response, err := h.svc.GetArticleByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h handler) GetArticles(ctx *gin.Context) {
	response, err := h.svc.GetArticles()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, *response)
}

func (h handler) GetArticlesByAuthor(ctx *gin.Context) {
	id := ctx.Param("author_id")
	response, err := h.svc.GetArticlesByAuthor(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (h handler) GetArticlesByTag(ctx *gin.Context) {
	tag := ctx.Param("tag_name")
	response, err := h.svc.GetArticlesByTag(tag)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (h handler) UpdateArticle(ctx *gin.Context) {
	article_id := ctx.Param("article_id")

	var res *domain.Article
	if err := ctx.ShouldBindJSON(&res); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	response, err := h.svc.UpdateArticle(article_id, res)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (h handler) DeleteArticle(ctx *gin.Context) {
	article_id := ctx.Param("article_id")
	err := h.svc.DeleteArticle(article_id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}

func (h handler) DeleteArticleAll(ctx *gin.Context) {
	err := h.svc.DeleteArticleAll()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}
