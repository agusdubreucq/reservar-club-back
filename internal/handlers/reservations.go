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

type ReservationHandler struct {
	reservationService *services.ReservationService
}

func NewReservationHandler(reservationService *services.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservationService: reservationService}
}

type CreateReservationRequest struct {
	CourtID   string    `json:"court_id" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	userID := c.GetString("userID")

	var req CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	reservation, err := h.reservationService.CreateReservation(c.Request.Context(), userID, req.CourtID, req.StartDate, req.EndDate)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidDateRange) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date range or duration (minimum 1 hour, multiple of 30 min, must start/end at :00 or :30)"})
			return
		}
		if errors.Is(err, domain.ErrCourtNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "court not found"})
			return
		}
		if errors.Is(err, domain.ErrReservationOverlap) {
			c.JSON(http.StatusConflict, gin.H{"error": "reservation overlaps with existing reservation"})
			return
		}
		if errors.Is(err, domain.ErrReservationOutsideSchedule) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "reservation is outside club schedule"})
			return
		}
		log.Printf("failed to create reservation: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, reservation)
}

func (h *ReservationHandler) GetReservation(c *gin.Context) {
	id := c.Param("id")

	reservation, err := h.reservationService.GetReservationByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("failed to get reservation: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "reservation not found"})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

func (h *ReservationHandler) ListUserReservations(c *gin.Context) {
	userID := c.GetString("userID")

	reservations, err := h.reservationService.ListUserReservations(c.Request.Context(), userID)
	if err != nil {
		log.Printf("failed to list reservations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if reservations == nil {
		reservations = []*domain.Reservation{}
	}

	c.JSON(http.StatusOK, reservations)
}

func (h *ReservationHandler) CancelReservation(c *gin.Context) {
	id := c.Param("id")

	reservation, err := h.reservationService.CancelReservation(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrReservationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "reservation not found"})
			return
		}
		log.Printf("failed to cancel reservation: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, reservation)
}
