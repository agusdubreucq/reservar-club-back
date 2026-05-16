package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/adubr/reservar-club-back/internal/domain"
	"github.com/adubr/reservar-club-back/internal/services"
	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	scheduleService *services.ScheduleService
}

func NewScheduleHandler(scheduleService *services.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{scheduleService: scheduleService}
}

type UpsertScheduleRequest struct {
	OpenTime string `json:"open_time" binding:"required"`
	CloseTime string `json:"close_time" binding:"required"`
	IsOpen   bool   `json:"is_open"`
}

func (h *ScheduleHandler) GetSchedule(c *gin.Context) {
	schedule, err := h.scheduleService.GetSchedule(c.Request.Context())
	if err != nil {
		log.Printf("failed to get schedule: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if schedule == nil {
		schedule = []*domain.ClubSchedule{}
	}

	c.JSON(http.StatusOK, schedule)
}

func (h *ScheduleHandler) UpsertSchedule(c *gin.Context) {
	dayStr := c.Param("day")
	day, err := strconv.Atoi(dayStr)
	if err != nil || day < 0 || day > 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid day of week (0-6)"})
		return
	}

	var req UpsertScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	schedule, err := h.scheduleService.UpsertSchedule(
		c.Request.Context(),
		day,
		req.OpenTime,
		req.CloseTime,
		req.IsOpen,
	)
	if err != nil {
		log.Printf("failed to upsert schedule: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, schedule)
}
