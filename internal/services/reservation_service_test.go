package services

import (
	"context"
	"testing"
	"time"

	"github.com/adubr/reservar-club-back/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateReservation(t *testing.T) {
	mockResRepo := new(MockReservationRepository)
	mockCourtRepo := new(MockCourtRepository)
	court := domain.NewCourt("court-1", "Court A", "sport-1")
	startDate := time.Now().UTC().Add(24 * time.Hour)
	endDate := startDate.Add(2 * time.Hour)

	mockCourtRepo.On("GetByID", mock.Anything, "court-1").Return(court, nil)
	mockResRepo.On("ListByCourtAndDateRange", mock.Anything, "court-1", startDate, endDate).Return([]*domain.Reservation{}, nil)
	mockResRepo.On("Create", mock.Anything, mock.MatchedBy(func(r *domain.Reservation) bool {
		return r.UserID == "user-1" && r.CourtID == "court-1"
	})).Return(nil)

	service := NewReservationService(mockResRepo, mockCourtRepo)
	reservation, err := service.CreateReservation(context.Background(), "user-1", "court-1", startDate, endDate)

	assert.NoError(t, err)
	assert.NotNil(t, reservation)
	assert.Equal(t, "user-1", reservation.UserID)
	assert.Equal(t, "court-1", reservation.CourtID)
}

func TestCreateReservationInvalidDateRange(t *testing.T) {
	mockResRepo := new(MockReservationRepository)
	mockCourtRepo := new(MockCourtRepository)
	startDate := time.Now().UTC()
	endDate := startDate.Add(-1 * time.Hour)

	service := NewReservationService(mockResRepo, mockCourtRepo)
	reservation, err := service.CreateReservation(context.Background(), "user-1", "court-1", startDate, endDate)

	assert.Equal(t, domain.ErrInvalidDateRange, err)
	assert.Nil(t, reservation)
}

func TestCreateReservationCourtNotFound(t *testing.T) {
	mockResRepo := new(MockReservationRepository)
	mockCourtRepo := new(MockCourtRepository)
	startDate := time.Now().UTC().Add(24 * time.Hour)
	endDate := startDate.Add(2 * time.Hour)

	mockCourtRepo.On("GetByID", mock.Anything, "nonexistent").Return(nil, domain.ErrNotFound)

	service := NewReservationService(mockResRepo, mockCourtRepo)
	reservation, err := service.CreateReservation(context.Background(), "user-1", "nonexistent", startDate, endDate)

	assert.Equal(t, domain.ErrCourtNotFound, err)
	assert.Nil(t, reservation)
}

func TestCreateReservationOverlap(t *testing.T) {
	mockResRepo := new(MockReservationRepository)
	mockCourtRepo := new(MockCourtRepository)
	court := domain.NewCourt("court-1", "Court A", "sport-1")
	startDate := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	existingRes := domain.NewReservation("res-1", "user-1", "court-1",
		time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 13, 0, 0, 0, time.UTC))

	mockCourtRepo.On("GetByID", mock.Anything, "court-1").Return(court, nil)
	mockResRepo.On("ListByCourtAndDateRange", mock.Anything, "court-1", startDate, endDate).Return([]*domain.Reservation{existingRes}, nil)

	service := NewReservationService(mockResRepo, mockCourtRepo)
	reservation, err := service.CreateReservation(context.Background(), "user-1", "court-1", startDate, endDate)

	assert.Equal(t, domain.ErrReservationOverlap, err)
	assert.Nil(t, reservation)
}

func TestGetReservationByID(t *testing.T) {
	mockResRepo := new(MockReservationRepository)
	mockCourtRepo := new(MockCourtRepository)
	expectedRes := domain.NewReservation("res-1", "user-1", "court-1",
		time.Now().UTC().Add(24*time.Hour),
		time.Now().UTC().Add(26*time.Hour))

	mockResRepo.On("GetByID", mock.Anything, "res-1").Return(expectedRes, nil)

	service := NewReservationService(mockResRepo, mockCourtRepo)
	reservation, err := service.GetReservationByID(context.Background(), "res-1")

	assert.NoError(t, err)
	assert.Equal(t, expectedRes, reservation)
}

func TestListUserReservations(t *testing.T) {
	mockResRepo := new(MockReservationRepository)
	mockCourtRepo := new(MockCourtRepository)
	startDate := time.Now().UTC().Add(24 * time.Hour)
	reservations := []*domain.Reservation{
		domain.NewReservation("res-1", "user-1", "court-1", startDate, startDate.Add(2*time.Hour)),
		domain.NewReservation("res-2", "user-1", "court-2", startDate.Add(48*time.Hour), startDate.Add(50*time.Hour)),
	}

	mockResRepo.On("ListByUserID", mock.Anything, "user-1").Return(reservations, nil)

	service := NewReservationService(mockResRepo, mockCourtRepo)
	result, err := service.ListUserReservations(context.Background(), "user-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestCancelReservation(t *testing.T) {
	mockResRepo := new(MockReservationRepository)
	mockCourtRepo := new(MockCourtRepository)
	startDate := time.Now().UTC().Add(24 * time.Hour)
	existingRes := domain.NewReservation("res-1", "user-1", "court-1", startDate, startDate.Add(2*time.Hour))

	mockResRepo.On("GetByID", mock.Anything, "res-1").Return(existingRes, nil)
	mockResRepo.On("Cancel", mock.Anything, "res-1").Return(nil)

	service := NewReservationService(mockResRepo, mockCourtRepo)
	reservation, err := service.CancelReservation(context.Background(), "res-1")

	assert.NoError(t, err)
	assert.Equal(t, domain.ReservationStatusCancelled, reservation.Status)
}

func TestCancelReservationNotFound(t *testing.T) {
	mockResRepo := new(MockReservationRepository)
	mockCourtRepo := new(MockCourtRepository)

	mockResRepo.On("GetByID", mock.Anything, "nonexistent").Return(nil, domain.ErrNotFound)

	service := NewReservationService(mockResRepo, mockCourtRepo)
	reservation, err := service.CancelReservation(context.Background(), "nonexistent")

	assert.Equal(t, domain.ErrReservationNotFound, err)
	assert.Nil(t, reservation)
}
