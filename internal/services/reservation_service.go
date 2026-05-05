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
}

func NewReservationService(reservationRepo repository.ReservationRepository, courtRepo repository.CourtRepository) *ReservationService {
	return &ReservationService{
		reservationRepo: reservationRepo,
		courtRepo:       courtRepo,
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

	existingReservations, err := s.reservationRepo.ListByCourtAndDateRange(ctx, courtID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	reservation := domain.NewReservation(uuid.New().String(), userID, courtID, startDate, endDate)

	for _, existing := range existingReservations {
		if existing.Status == domain.ReservationStatusActive && reservation.IsOverlapping(existing) {
			return nil, domain.ErrReservationOverlap
		}
	}

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
