package app

import (
	"fmt"
	"net/http"

	"github.com/AntonyIS/notlify-content-svc/internal/core/domain"
	"github.com/AntonyIS/notlify-content-svc/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type GinHandler interface {
	CreateContent(ctx *gin.Context)
	ReadContent(ctx *gin.Context)
	ReadCreatorContents(ctx *gin.Context)
	ReadContents(ctx *gin.Context)
	UpdateContent(ctx *gin.Context)
	DeleteContent(ctx *gin.Context)
}

type handler struct {
	svc       ports.ContentService
	secretKey string
}

func NewGinHandler(svc ports.ContentService, secretKey string) GinHandler {
	routerHandler := handler{
		svc:       svc,
		secretKey: secretKey,
	}

	return routerHandler
}

func (h handler) CreateContent(ctx *gin.Context) {
	var res *domain.Content
	if err := ctx.ShouldBindJSON(&res); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Println(res)
	content, err := h.svc.CreateContent(res)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, content)
}

func (h handler) ReadContent(ctx *gin.Context) {
	id := ctx.Param("id")
	content, err := h.svc.ReadContent(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, content)
}

func (h handler) ReadCreatorContents(ctx *gin.Context) {
	creator_id := ctx.Param("creator_id")
	contents, err := h.svc.ReadCreatorContents(creator_id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, contents)
}

func (h handler) ReadContents(ctx *gin.Context) {
	contents, err := h.svc.ReadContents()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, contents)
}

func (h handler) UpdateContent(ctx *gin.Context) {
	id := ctx.Param("id")
	_, err := h.svc.ReadContent(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}

	var res *domain.Content
	if err := ctx.ShouldBindJSON(&res); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	res.ContentId = id
	content, err := h.svc.UpdateContent(res)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, content)
}

func (h handler) DeleteContent(ctx *gin.Context) {
	id := ctx.Param("id")
	message, err := h.svc.DeleteContent(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
