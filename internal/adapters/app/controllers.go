package app

import (
	"net/http"

	"github.com/AntonyIS/notlify-content-svc/internal/core/domain"
	"github.com/AntonyIS/notlify-content-svc/internal/core/ports"
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
	content, err := h.svc.CreateArticle(res)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, content)
}

func (h handler) GetArticleByID(ctx *gin.Context) {
	id := ctx.Param("id")
	article, err := h.svc.GetArticleByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, article)
}

func (h handler) GetArticles(ctx *gin.Context) {
	articles, err := h.svc.GetArticles()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, *articles)
}

func (h handler) GetArticlesByAuthor(ctx *gin.Context) {
	id := ctx.Param("id")
	articles, err := h.svc.GetArticlesByAuthor(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, articles)
}

func (h handler) GetArticlesByTag(ctx *gin.Context) {
	tag := ctx.Param("tag")
	articles, err := h.svc.GetArticlesByTag(tag)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, articles)
}

func (h handler) UpdateArticle(ctx *gin.Context) {
	var res *domain.Article
	if err := ctx.ShouldBindJSON(&res); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	article, err := h.svc.UpdateArticle(res)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, article)
}

func (h handler) DeleteArticle(ctx *gin.Context) {
	id := ctx.Param("id")
	err := h.svc.DeleteArticle(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}

func (h handler) DeleteArticleAll(ctx *gin.Context) {
	err := h.svc.DeleteArticleAll()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}
