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

// @Summary Create a new Franquicia
// @Description Create a new Franquicia with the given details
// @Tags franquicia
// @Accept  json
// @Produce  json
// @Param   FranquiciaRequest  body  domain.FranquiciaRequest  true  "Franquicia Request"
// @Success 201  {object}  map[string]interface{}
// @Failure 400,500  {object}  map[string]interface{}
// @Router /franquicia [post]
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

// @Summary Get Franquicias by Location
// @Description Retrieves franquicias based on given location parameters
// @Tags franquicia
// @Accept  json
// @Produce  json
// @Param   city     query     string     true     "City"
// @Param   country  query     string     true     "Country"
// @Success 200 {array} domain.Franquicia
// @Failure 400,404,500 {object} map[string]interface{}
// @Router /franquicia/location [get]
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

// @Summary Get Franquicias by Date Range
// @Description Retrieves franquicias within a specified date range
// @Tags franquicia
// @Accept  json
// @Produce  json
// @Param   start    query     string     true     "Start Date"
// @Param   end      query     string     true     "End Date"
// @Success 200 {array} domain.Franquicia
// @Failure 400,404,500 {object} map[string]interface{}
// @Router /franquicia/daterange [get]
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

// @Summary Get Franquicia by ID
// @Description Retrieves a franquicia by its ID
// @Tags franquicia
// @Accept  json
// @Produce  json
// @Param   id       path      string     true     "Franquicia ID"
// @Success 200 {object} domain.Franquicia
// @Failure 400,404,500 {object} map[string]interface{}
// @Router /franquicia/{id} [get]
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

// @Summary Get Franquicias by Name
// @Description Retrieves franquicias by their name
// @Tags franquicia
// @Accept  json
// @Produce  json
// @Param   name     query     string     true     "Franquicia Name"
// @Success 200 {array} domain.Franquicia
// @Failure 400,404,500 {object} map[string]interface{}
// @Router /franquicia/name [get]
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

// @Summary Get All Franquicias
// @Description Retrieves all franquicias
// @Tags franquicia
// @Accept  json
// @Produce  json
// @Success 200 {array} domain.Franquicia
// @Failure 500 {object} map[string]interface{}
// @Router /franquicias [get]
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

// @Summary Update Franquicia
// @Description Updates a franquicia by its ID
// @Tags franquicia
// @Accept  json
// @Produce  json
// @Param   id       path      string     true     "Franquicia ID"
// @Param   FranquiciaRequest  body      domain.FranquiciaRequest  true  "Franquicia Update Request"
// @Success 200 {object} map[string]string
// @Failure 400,404,500 {object} map[string]interface{}
// @Router /franquicia/{id} [put]
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
