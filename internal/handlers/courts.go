package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/adubr/reservar-club-back/internal/domain"
	"github.com/adubr/reservar-club-back/internal/services"
	"github.com/gin-gonic/gin"
)

type CourtHandler struct {
	courtService *services.CourtService
}

func NewCourtHandler(courtService *services.CourtService) *CourtHandler {
	return &CourtHandler{courtService: courtService}
}

type CreateCourtRequest struct {
	Name    string `json:"name" binding:"required"`
	SportID string `json:"sport_id" binding:"required"`
}

type UpdateCourtRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *CourtHandler) CreateCourt(c *gin.Context) {
	var req CreateCourtRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	court, err := h.courtService.CreateCourt(c.Request.Context(), req.Name, req.SportID)
	if err != nil {
		if errors.Is(err, domain.ErrSportNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "sport not found"})
			return
		}
		log.Printf("failed to create court: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, court)
}

func (h *CourtHandler) ListCourts(c *gin.Context) {
	courts, err := h.courtService.ListAllCourts(c.Request.Context())
	if err != nil {
		log.Printf("failed to list courts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if courts == nil {
		courts = []*domain.Court{}
	}

	c.JSON(http.StatusOK, courts)
}

func (h *CourtHandler) GetCourt(c *gin.Context) {
	id := c.Param("id")

	court, err := h.courtService.GetCourtByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("failed to get court: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "court not found"})
		return
	}

	c.JSON(http.StatusOK, court)
}

func (h *CourtHandler) ListCourtsBySport(c *gin.Context) {
	sportID := c.Param("sport_id")

	courts, err := h.courtService.ListCourtsBySport(c.Request.Context(), sportID)
	if err != nil {
		log.Printf("failed to list courts by sport: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if courts == nil {
		courts = []*domain.Court{}
	}

	c.JSON(http.StatusOK, courts)
}

func (h *CourtHandler) UpdateCourt(c *gin.Context) {
	id := c.Param("id")

	var req UpdateCourtRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	court, err := h.courtService.UpdateCourt(c.Request.Context(), id, req.Name)
	if err != nil {
		if errors.Is(err, domain.ErrCourtNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "court not found"})
			return
		}
		log.Printf("failed to update court: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, court)
}

func (h *CourtHandler) DeleteCourt(c *gin.Context) {
	id := c.Param("id")

	err := h.courtService.DeleteCourt(c.Request.Context(), id)
	if err != nil {
		log.Printf("failed to delete court: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
