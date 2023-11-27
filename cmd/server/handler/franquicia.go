package handler

import (
	"clubhub-hotel-management/internal/domain"
	"clubhub-hotel-management/internal/franquicia"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (f *Franquicia) GetByLocation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		city := ctx.Query("city")
		country := ctx.Query("country")

		franquicias, err := f.service.GetByLocation(ctx, city, country)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, franquicias)
	}
}

func (f *Franquicia) GetFranquiciasByDateRange() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startDate := ctx.Query("start")
		endDate := ctx.Query("end")

		franquicias, err := f.service.GetByDateRange(ctx, startDate, endDate)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, franquicias)
	}
}

func (f *Franquicia) GetFranquiciaByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		franquicia, err := f.service.GetFranquiciaByID(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, franquicia)
	}
}

func (h *Franquicia) GetFranquiciasByName() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Query("name")

		franquicias, err := h.service.GetByFranchiseName(ctx, name)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, franquicias)
	}

}
func (f *Franquicia) GetAllFranquicias() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		franquicias, err := f.service.GetAllFranquicias(ctx)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, franquicias)
	}
}

func (f *Franquicia) UpdateFranquicia() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "id is empty"})
			return
		}

		var req domain.FranquiciaRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fr := domain.Franquicia{
			ID:       objID,
			Name:     req.Name,
			URL:      req.URL,
			Location: req.Location,
		}

		if err := f.service.UpdateFranquicia(ctx, fr); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Franquicia actualizada correctamente"})
	}
}
