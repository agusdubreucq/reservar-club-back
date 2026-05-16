package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/adubr/reservar-club-back/internal/domain"
	"github.com/adubr/reservar-club-back/internal/services"
	"github.com/gin-gonic/gin"
)

type PricingHandler struct {
	pricingService *services.PricingService
}

func NewPricingHandler(pricingService *services.PricingService) *PricingHandler {
	return &PricingHandler{pricingService: pricingService}
}

type EstimatePriceRequest struct {
	SportID   string    `form:"sport_id" binding:"required"`
	StartDate time.Time `form:"start_date" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	EndDate   time.Time `form:"end_date" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
}

type EstimatePriceResponse struct {
	TotalPrice float64 `json:"total_price"`
}

type CreatePricingRuleRequest struct {
	SportID      string  `json:"sport_id" binding:"required"`
	DayOfWeek    *int    `json:"day_of_week"`
	StartHour    int     `json:"start_hour" binding:"required"`
	EndHour      int     `json:"end_hour" binding:"required"`
	PricePerHour float64 `json:"price_per_hour" binding:"required"`
}

func (h *PricingHandler) EstimatePrice(c *gin.Context) {
	var req EstimatePriceRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	price, err := h.pricingService.CalculatePrice(c.Request.Context(), req.SportID, req.StartDate, req.EndDate)
	if err != nil {
		log.Printf("failed to estimate price: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, EstimatePriceResponse{TotalPrice: price})
}

func (h *PricingHandler) CreatePricingRule(c *gin.Context) {
	var req CreatePricingRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	rule, err := h.pricingService.CreatePricingRule(
		c.Request.Context(),
		req.SportID,
		req.DayOfWeek,
		req.StartHour,
		req.EndHour,
		req.PricePerHour,
	)
	if err != nil {
		log.Printf("failed to create pricing rule: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, rule)
}

func (h *PricingHandler) ListPricingRules(c *gin.Context) {
	sportID := c.Query("sport_id")
	if sportID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sport_id is required"})
		return
	}

	rules, err := h.pricingService.ListPricingRules(c.Request.Context(), sportID)
	if err != nil {
		log.Printf("failed to list pricing rules: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if rules == nil {
		rules = []*domain.PricingRule{}
	}

	c.JSON(http.StatusOK, rules)
}

func (h *PricingHandler) DeletePricingRule(c *gin.Context) {
	id := c.Param("id")

	err := h.pricingService.DeletePricingRule(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "pricing rule not found"})
			return
		}
		log.Printf("failed to delete pricing rule: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
