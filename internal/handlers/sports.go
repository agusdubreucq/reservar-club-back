package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/adubr/reservar-club-back/internal/domain"
	"github.com/adubr/reservar-club-back/internal/services"
	"github.com/gin-gonic/gin"
)

type SportHandler struct {
	sportService *services.SportService
}

func NewSportHandler(sportService *services.SportService) *SportHandler {
	return &SportHandler{sportService: sportService}
}

type CreateSportRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *SportHandler) CreateSport(c *gin.Context) {
	var req CreateSportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	sport, err := h.sportService.CreateSport(c.Request.Context(), req.Name)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateSportName) {
			c.JSON(http.StatusConflict, gin.H{"error": "sport name already exists"})
			return
		}
		log.Printf("failed to create sport: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, sport)
}

func (h *SportHandler) ListSports(c *gin.Context) {
	sports, err := h.sportService.ListAllSports(c.Request.Context())
	if err != nil {
		log.Printf("failed to list sports: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if sports == nil {
		sports = []*domain.Sport{}
	}

	c.JSON(http.StatusOK, sports)
}

func (h *SportHandler) GetSport(c *gin.Context) {
	id := c.Param("id")

	sport, err := h.sportService.GetSportByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("failed to get sport: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "sport not found"})
		return
	}

	c.JSON(http.StatusOK, sport)
}

func (h *SportHandler) DeleteSport(c *gin.Context) {
	id := c.Param("id")

	err := h.sportService.DeleteSport(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrSportHasCourts) {
			c.JSON(http.StatusConflict, gin.H{"error": "cannot delete sport with associated courts"})
			return
		}
		log.Printf("failed to delete sport: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
