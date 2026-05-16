package services

import (
	"context"
	"time"

	"github.com/adubr/reservar-club-back/internal/domain"
	"github.com/adubr/reservar-club-back/internal/repository"
	"github.com/google/uuid"
)

type ReservationService struct {
	reservationRepo repository.ReservationRepository
	courtRepo       repository.CourtRepository
	scheduleRepo    repository.ClubScheduleRepository
	pricingService  *PricingService
}

func NewReservationService(
	reservationRepo repository.ReservationRepository,
	courtRepo repository.CourtRepository,
	scheduleRepo repository.ClubScheduleRepository,
	pricingService *PricingService,
) *ReservationService {
	return &ReservationService{
		reservationRepo: reservationRepo,
		courtRepo:       courtRepo,
		scheduleRepo:    scheduleRepo,
		pricingService:  pricingService,
	}
}

func (s *ReservationService) CreateReservation(ctx context.Context, userID, courtID string, startDate, endDate time.Time) (*domain.Reservation, error) {
	if !endDate.After(startDate) {
		return nil, domain.ErrInvalidDateRange
	}

	court, err := s.courtRepo.GetByID(ctx, courtID)
	if err != nil {
		return nil, domain.ErrCourtNotFound
	}

	if court == nil {
		return nil, domain.ErrCourtNotFound
	}

	reservation := domain.NewReservation(uuid.New().String(), userID, courtID, startDate, endDate, 0)

	if !reservation.IsValidDuration() {
		return nil, domain.ErrInvalidDateRange
	}

	if !reservation.HasValidSlotStartTimes() {
		return nil, domain.ErrInvalidDateRange
	}

	schedule, err := s.scheduleRepo.GetByDayOfWeek(ctx, int(startDate.Weekday()))
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}

	if schedule != nil && !schedule.IsOpen {
		return nil, domain.ErrReservationOutsideSchedule
	}

	if schedule != nil {
		if !s.isWithinSchedule(startDate, endDate, schedule) {
			return nil, domain.ErrReservationOutsideSchedule
		}
	}

	existingReservations, err := s.reservationRepo.ListByCourtAndDateRange(ctx, courtID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	for _, existing := range existingReservations {
		if existing.Status == domain.ReservationStatusActive && reservation.IsOverlapping(existing) {
			return nil, domain.ErrReservationOverlap
		}
	}

	price, err := s.pricingService.CalculatePrice(ctx, court.SportID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	reservation.TotalPrice = price

	if err := s.reservationRepo.Create(ctx, reservation); err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *ReservationService) GetReservationByID(ctx context.Context, id string) (*domain.Reservation, error) {
	return s.reservationRepo.GetByID(ctx, id)
}

func (s *ReservationService) ListUserReservations(ctx context.Context, userID string) ([]*domain.Reservation, error) {
	return s.reservationRepo.ListByUserID(ctx, userID)
}

func (s *ReservationService) CancelReservation(ctx context.Context, id string) (*domain.Reservation, error) {
	reservation, err := s.reservationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrReservationNotFound
	}

	if reservation == nil {
		return nil, domain.ErrReservationNotFound
	}

	reservation.Cancel()
	if err := s.reservationRepo.Cancel(ctx, id); err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *ReservationService) isWithinSchedule(startDate, endDate time.Time, schedule *domain.ClubSchedule) bool {
	openParts := timeStringToMinutes(schedule.OpenTime)
	closeParts := timeStringToMinutes(schedule.CloseTime)

	startMin := startDate.Hour()*60 + startDate.Minute()
	endMin := endDate.Hour()*60 + endDate.Minute()

	if startMin < openParts || startMin >= closeParts {
		return false
	}

	if endMin <= openParts || endMin > closeParts {
		return false
	}

	return true
}

func (s *ReservationService) GetReservationsForCourtInRange(ctx context.Context, courtID string, start, end time.Time) ([]*domain.Reservation, error) {
	return s.reservationRepo.ListByCourtAndDateRange(ctx, courtID, start, end)
}

func timeStringToMinutes(timeStr string) int {
	n, _ := time.Parse("15:04", timeStr)
	return n.Hour()*60 + n.Minute()
}
