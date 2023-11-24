package handler

import (
	"clubhub-hotel-management/internal/domain"
	"clubhub-hotel-management/internal/franquicia"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Franquicia struct {
	service franquicia.Service
}

func NewUser(service franquicia.Service) *Franquicia {
	return &Franquicia{service: service}
}

func (f *Franquicia) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req domain.FranquiciaRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		if req.URL == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "name and URL are required"})
			return
		}

		franquicia := &domain.Franquicia{
			URL: req.URL,
		}

		err := f.service.CreateFranquicia(ctx, franquicia)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"status": "ok"})
	}
}

func (h *Franquicia) GetAllFranquicias() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		// defer cancel()

		franquicias, err := h.service.GetAllFranquicias(ctx)
		if err != nil {
			ctx.Error(err) // Gin manejará el error
			return
		}

		ctx.JSON(http.StatusOK, franquicias)
	}
}
