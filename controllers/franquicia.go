package controllers

import (
	"clubhub-hotel-management/internal/domain"
	"clubhub-hotel-management/internal/franquicia"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type Franquicia struct {
	service franquicia.Service
}

func NewUser(service franquicia.Service) *Franquicia {
	return &Franquicia{service: service}
}

func (f *Franquicia) Create(c *gin.Context) {
	var req domain.Franquicia
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.Name == "" || req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and URL are required"})
		return
	}

	if _, err := url.ParseRequestURI(req.URL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	err := f.service.CreateFranquicia(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "ok"})
}
