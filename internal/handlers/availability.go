package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/adubr/reservar-club-back/internal/domain"
	"github.com/adubr/reservar-club-back/internal/services"
	"github.com/gin-gonic/gin"
)

type AvailabilityHandler struct {
	reservationService *services.ReservationService
	scheduleService    *services.ScheduleService
	courtService       *services.CourtService
}

func NewAvailabilityHandler(
	reservationService *services.ReservationService,
	scheduleService *services.ScheduleService,
	courtService *services.CourtService,
) *AvailabilityHandler {
	return &AvailabilityHandler{
		reservationService: reservationService,
		scheduleService:    scheduleService,
		courtService:       courtService,
	}
}

type TimeSlot struct {
	Start       string                `json:"start"`
	End         string                `json:"end"`
	Available   bool                  `json:"available"`
	Reservation *ReservationSlotInfo  `json:"reservation"`
}

type ReservationSlotInfo struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type CourtAvailability struct {
	ID    string      `json:"id"`
	Name  string      `json:"name"`
	Slots []TimeSlot  `json:"slots"`
}

type AvailabilityResponse struct {
	Date     string                 `json:"date"`
	SportID  string                 `json:"sport_id"`
	Schedule ScheduleInfo           `json:"schedule"`
	Courts   []CourtAvailability    `json:"courts"`
}

type ScheduleInfo struct {
	Open  string `json:"open"`
	Close string `json:"close"`
	IsOpen bool `json:"is_open"`
}

func (h *AvailabilityHandler) GetAvailability(c *gin.Context) {
	dateStr := c.Query("date")
	sportID := c.Query("sport_id")

	if dateStr == "" || sportID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date and sport_id are required"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format (use YYYY-MM-DD)"})
		return
	}

	dayOfWeek := int(date.Weekday())
	schedule, err := h.scheduleService.GetScheduleForDay(c.Request.Context(), dayOfWeek)
	if err != nil && err != domain.ErrNotFound {
		log.Printf("failed to get schedule: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if schedule == nil {
		schedule = &domain.ClubSchedule{IsOpen: true, OpenTime: "00:00", CloseTime: "23:59"}
	}

	if !schedule.IsOpen {
		c.JSON(http.StatusOK, AvailabilityResponse{
			Date:    dateStr,
			SportID: sportID,
			Schedule: ScheduleInfo{
				Open:   schedule.OpenTime,
				Close:  schedule.CloseTime,
				IsOpen: false,
			},
			Courts: []CourtAvailability{},
		})
		return
	}

	courts, err := h.courtService.ListCourtsBySport(c.Request.Context(), sportID)
	if err != nil {
		log.Printf("failed to get courts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var courtAvailabilities []CourtAvailability
	for _, court := range courts {
		slots := h.generateTimeSlots(c.Request.Context(), court, date, schedule)
		courtAvailabilities = append(courtAvailabilities, CourtAvailability{
			ID:    court.ID,
			Name:  court.Name,
			Slots: slots,
		})
	}

	c.JSON(http.StatusOK, AvailabilityResponse{
		Date:    dateStr,
		SportID: sportID,
		Schedule: ScheduleInfo{
			Open:   schedule.OpenTime,
			Close:  schedule.CloseTime,
			IsOpen: true,
		},
		Courts: courtAvailabilities,
	})
}

func (h *AvailabilityHandler) generateTimeSlots(ctx context.Context, court *domain.Court, date time.Time, schedule *domain.ClubSchedule) []TimeSlot {
	openTime := parseTime(schedule.OpenTime)
	closeTime := parseTime(schedule.CloseTime)

	dayStart := time.Date(date.Year(), date.Month(), date.Day(), openTime.Hour(), openTime.Minute(), 0, 0, date.Location())
	dayEnd := time.Date(date.Year(), date.Month(), date.Day(), closeTime.Hour(), closeTime.Minute(), 0, 0, date.Location())

	var slots []TimeSlot
	current := dayStart

	for current.Before(dayEnd) {
		slotEnd := current.Add(30 * time.Minute)
		if slotEnd.After(dayEnd) {
			break
		}

		available, reservation := h.isSlotAvailable(ctx, court.ID, current, slotEnd)

		slots = append(slots, TimeSlot{
			Start:       current.Format("15:04"),
			End:         slotEnd.Format("15:04"),
			Available:   available,
			Reservation: reservation,
		})

		current = slotEnd
	}

	return slots
}

func (h *AvailabilityHandler) isSlotAvailable(ctx context.Context, courtID string, start, end time.Time) (bool, *ReservationSlotInfo) {
	reservations, err := h.reservationService.GetReservationsForCourtInRange(ctx, courtID, start, end)
	if err != nil {
		return true, nil
	}

	for _, res := range reservations {
		if res.Status == domain.ReservationStatusActive {
			return false, &ReservationSlotInfo{
				ID:     res.ID,
				UserID: res.UserID,
			}
		}
	}

	return true, nil
}

func parseTime(timeStr string) time.Time {
	t, _ := time.Parse("15:04", timeStr)
	return t
}
