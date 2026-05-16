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
	mockScheduleRepo := new(MockClubScheduleRepository)
	mockPricingRepo := new(MockPricingRuleRepository)

	court := domain.NewCourt("court-1", "Court A", "sport-1")
	startDate := time.Date(2025, 1, 2, 10, 0, 0, 0, time.UTC)
	endDate := startDate.Add(2 * time.Hour)

	mockCourtRepo.On("GetByID", mock.Anything, "court-1").Return(court, nil)
	mockScheduleRepo.On("GetByDayOfWeek", mock.Anything, 4).Return(&domain.ClubSchedule{
		IsOpen:   true,
		OpenTime: "08:00",
		CloseTime: "22:00",
	}, nil)
	mockResRepo.On("ListByCourtAndDateRange", mock.Anything, "court-1", startDate, endDate).Return([]*domain.Reservation{}, nil)
	mockPricingRepo.On("ListBySportID", mock.Anything, "sport-1").Return([]*domain.PricingRule{}, nil)
	mockResRepo.On("Create", mock.Anything, mock.MatchedBy(func(r *domain.Reservation) bool {
		return r.UserID == "user-1" && r.CourtID == "court-1"
	})).Return(nil)

	pricingService := NewPricingService(mockPricingRepo)
	service := NewReservationService(mockResRepo, mockCourtRepo, mockScheduleRepo, pricingService)
	reservation, err := service.CreateReservation(context.Background(), "user-1", "court-1", startDate, endDate)

	assert.NoError(t, err)
	assert.NotNil(t, reservation)
	assert.Equal(t, "user-1", reservation.UserID)
	assert.Equal(t, "court-1", reservation.CourtID)
}

func TestCreateReservationInvalidDateRange(t *testing.T) {
	mockResRepo := new(MockReservationRepository)
	mockCourtRepo := new(MockCourtRepository)
	mockScheduleRepo := new(MockClubScheduleRepository)
	mockPricingRepo := new(MockPricingRuleRepository)

	startDate := time.Now().UTC()
	endDate := startDate.Add(-1 * time.Hour)

	pricingService := NewPricingService(mockPricingRepo)
	service := NewReservationService(mockResRepo, mockCourtRepo, mockScheduleRepo, pricingService)
	reservation, err := service.CreateReservation(context.Background(), "user-1", "court-1", startDate, endDate)

	assert.Equal(t, domain.ErrInvalidDateRange, err)
	assert.Nil(t, reservation)
}

func TestCreateReservationCourtNotFound(t *testing.T) {
	mockResRepo := new(MockReservationRepository)
	mockCourtRepo := new(MockCourtRepository)
	mockScheduleRepo := new(MockClubScheduleRepository)
	mockPricingRepo := new(MockPricingRuleRepository)

	startDate := time.Date(2025, 1, 2, 10, 0, 0, 0, time.UTC)
	endDate := startDate.Add(2 * time.Hour)

	mockCourtRepo.On("GetByID", mock.Anything, "nonexistent").Return(nil, domain.ErrNotFound)

	pricingService := NewPricingService(mockPricingRepo)
	service := NewReservationService(mockResRepo, mockCourtRepo, mockScheduleRepo, pricingService)
	reservation, err := service.CreateReservation(context.Background(), "user-1", "nonexistent", startDate, endDate)

	assert.Equal(t, domain.ErrCourtNotFound, err)
	assert.Nil(t, reservation)
}

func TestCreateReservationOverlap(t *testing.T) {
	mockResRepo := new(MockReservationRepository)
	mockCourtRepo := new(MockCourtRepository)
	mockScheduleRepo := new(MockClubScheduleRepository)
	mockPricingRepo := new(MockPricingRuleRepository)

	court := domain.NewCourt("court-1", "Court A", "sport-1")
	startDate := time.Date(2025, 1, 2, 10, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 2, 12, 0, 0, 0, time.UTC)

	existingRes := domain.NewReservation("res-1", "user-1", "court-1",
		time.Date(2025, 1, 2, 11, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 2, 13, 0, 0, 0, time.UTC), 0)

	mockCourtRepo.On("GetByID", mock.Anything, "court-1").Return(court, nil)
	mockScheduleRepo.On("GetByDayOfWeek", mock.Anything, 4).Return(&domain.ClubSchedule{
		IsOpen:   true,
		OpenTime: "08:00",
		CloseTime: "22:00",
	}, nil)
	mockResRepo.On("ListByCourtAndDateRange", mock.Anything, "court-1", startDate, endDate).Return([]*domain.Reservation{existingRes}, nil)

	pricingService := NewPricingService(mockPricingRepo)
	service := NewReservationService(mockResRepo, mockCourtRepo, mockScheduleRepo, pricingService)
	reservation, err := service.CreateReservation(context.Background(), "user-1", "court-1", startDate, endDate)

	assert.Equal(t, domain.ErrReservationOverlap, err)
	assert.Nil(t, reservation)
}

func TestGetReservationByID(t *testing.T) {
	mockResRepo := new(MockReservationRepository)
	mockCourtRepo := new(MockCourtRepository)
	mockScheduleRepo := new(MockClubScheduleRepository)
	mockPricingRepo := new(MockPricingRuleRepository)

	expectedRes := domain.NewReservation("res-1", "user-1", "court-1",
		time.Now().UTC().Add(24*time.Hour),
		time.Now().UTC().Add(26*time.Hour), 50.0)

	mockResRepo.On("GetByID", mock.Anything, "res-1").Return(expectedRes, nil)

	pricingService := NewPricingService(mockPricingRepo)
	service := NewReservationService(mockResRepo, mockCourtRepo, mockScheduleRepo, pricingService)
	reservation, err := service.GetReservationByID(context.Background(), "res-1")

	assert.NoError(t, err)
	assert.Equal(t, expectedRes, reservation)
}

func TestListUserReservations(t *testing.T) {
	mockResRepo := new(MockReservationRepository)
	mockCourtRepo := new(MockCourtRepository)
	mockScheduleRepo := new(MockClubScheduleRepository)
	mockPricingRepo := new(MockPricingRuleRepository)

	startDate := time.Now().UTC().Add(24 * time.Hour)
	reservations := []*domain.Reservation{
		domain.NewReservation("res-1", "user-1", "court-1", startDate, startDate.Add(2*time.Hour), 100.0),
		domain.NewReservation("res-2", "user-1", "court-2", startDate.Add(48*time.Hour), startDate.Add(50*time.Hour), 120.0),
	}

	mockResRepo.On("ListByUserID", mock.Anything, "user-1").Return(reservations, nil)

	pricingService := NewPricingService(mockPricingRepo)
	service := NewReservationService(mockResRepo, mockCourtRepo, mockScheduleRepo, pricingService)
	result, err := service.ListUserReservations(context.Background(), "user-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestCancelReservation(t *testing.T) {
	mockResRepo := new(MockReservationRepository)
	mockCourtRepo := new(MockCourtRepository)
	mockScheduleRepo := new(MockClubScheduleRepository)
	mockPricingRepo := new(MockPricingRuleRepository)

	startDate := time.Now().UTC().Add(24 * time.Hour)
	existingRes := domain.NewReservation("res-1", "user-1", "court-1", startDate, startDate.Add(2*time.Hour), 75.0)

	mockResRepo.On("GetByID", mock.Anything, "res-1").Return(existingRes, nil)
	mockResRepo.On("Cancel", mock.Anything, "res-1").Return(nil)

	pricingService := NewPricingService(mockPricingRepo)
	service := NewReservationService(mockResRepo, mockCourtRepo, mockScheduleRepo, pricingService)
	reservation, err := service.CancelReservation(context.Background(), "res-1")

	assert.NoError(t, err)
	assert.Equal(t, domain.ReservationStatusCancelled, reservation.Status)
}

func TestCancelReservationNotFound(t *testing.T) {
	mockResRepo := new(MockReservationRepository)
	mockCourtRepo := new(MockCourtRepository)
	mockScheduleRepo := new(MockClubScheduleRepository)
	mockPricingRepo := new(MockPricingRuleRepository)

	mockResRepo.On("GetByID", mock.Anything, "nonexistent").Return(nil, domain.ErrNotFound)

	pricingService := NewPricingService(mockPricingRepo)
	service := NewReservationService(mockResRepo, mockCourtRepo, mockScheduleRepo, pricingService)
	reservation, err := service.CancelReservation(context.Background(), "nonexistent")

	assert.Equal(t, domain.ErrReservationNotFound, err)
	assert.Nil(t, reservation)
}
