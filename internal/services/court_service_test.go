package services

import (
	"context"
	"testing"

	"github.com/adubr/reservar-club-back/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCourt(t *testing.T) {
	mockCourtRepo := new(MockCourtRepository)
	mockSportRepo := new(MockSportRepository)
	sport := domain.NewSport("sport-1", "Tennis")

	mockSportRepo.On("GetByID", mock.Anything, "sport-1").Return(sport, nil)
	mockCourtRepo.On("Create", mock.Anything, mock.MatchedBy(func(c *domain.Court) bool {
		return c.Name == "Court A" && c.SportID == "sport-1"
	})).Return(nil)

	service := NewCourtService(mockCourtRepo, mockSportRepo)
	court, err := service.CreateCourt(context.Background(), "Court A", "sport-1")

	assert.NoError(t, err)
	assert.NotNil(t, court)
	assert.Equal(t, "Court A", court.Name)
	assert.Equal(t, "sport-1", court.SportID)
}

func TestCreateCourtSportNotFound(t *testing.T) {
	mockCourtRepo := new(MockCourtRepository)
	mockSportRepo := new(MockSportRepository)

	mockSportRepo.On("GetByID", mock.Anything, "nonexistent").Return(nil, domain.ErrNotFound)

	service := NewCourtService(mockCourtRepo, mockSportRepo)
	court, err := service.CreateCourt(context.Background(), "Court A", "nonexistent")

	assert.Equal(t, domain.ErrSportNotFound, err)
	assert.Nil(t, court)
}

func TestGetCourtByID(t *testing.T) {
	mockCourtRepo := new(MockCourtRepository)
	mockSportRepo := new(MockSportRepository)
	expectedCourt := domain.NewCourt("court-1", "Court A", "sport-1")

	mockCourtRepo.On("GetByID", mock.Anything, "court-1").Return(expectedCourt, nil)

	service := NewCourtService(mockCourtRepo, mockSportRepo)
	court, err := service.GetCourtByID(context.Background(), "court-1")

	assert.NoError(t, err)
	assert.Equal(t, expectedCourt, court)
}

func TestListAllCourts(t *testing.T) {
	mockCourtRepo := new(MockCourtRepository)
	mockSportRepo := new(MockSportRepository)
	courts := []*domain.Court{
		domain.NewCourt("court-1", "Court A", "sport-1"),
		domain.NewCourt("court-2", "Court B", "sport-2"),
	}

	mockCourtRepo.On("ListAll", mock.Anything).Return(courts, nil)

	service := NewCourtService(mockCourtRepo, mockSportRepo)
	result, err := service.ListAllCourts(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestUpdateCourt(t *testing.T) {
	mockCourtRepo := new(MockCourtRepository)
	mockSportRepo := new(MockSportRepository)
	existingCourt := domain.NewCourt("court-1", "Court A", "sport-1")

	mockCourtRepo.On("GetByID", mock.Anything, "court-1").Return(existingCourt, nil)
	mockCourtRepo.On("Update", mock.Anything, mock.MatchedBy(func(c *domain.Court) bool {
		return c.Name == "Court Updated"
	})).Return(nil)

	service := NewCourtService(mockCourtRepo, mockSportRepo)
	court, err := service.UpdateCourt(context.Background(), "court-1", "Court Updated")

	assert.NoError(t, err)
	assert.Equal(t, "Court Updated", court.Name)
}

func TestDeleteCourt(t *testing.T) {
	mockCourtRepo := new(MockCourtRepository)
	mockSportRepo := new(MockSportRepository)

	mockCourtRepo.On("Delete", mock.Anything, "court-1").Return(nil)

	service := NewCourtService(mockCourtRepo, mockSportRepo)
	err := service.DeleteCourt(context.Background(), "court-1")

	assert.NoError(t, err)
}
