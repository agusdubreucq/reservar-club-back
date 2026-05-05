package services

import (
	"context"
	"testing"

	"github.com/adubr/reservar-club-back/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSport(t *testing.T) {
	mockRepo := new(MockSportRepository)

	mockRepo.On("GetByName", mock.Anything, "Tennis").Return(nil, domain.ErrNotFound)
	mockRepo.On("Creates", mock.Anything, mock.MatchedBy(func(s *domain.Sport) bool {
		return s.Name == "Tennis"
	})).Return(nil)

	service := NewSportService(mockRepo)
	sport, err := service.CreateSport(context.Background(), "Tennis")

	assert.NoError(t, err)
	assert.NotNil(t, sport)
	assert.Equal(t, "Tennis", sport.Name)
}

func TestCreateSportDuplicate(t *testing.T) {
	mockRepo := new(MockSportRepository)
	existingSport := domain.NewSport("sport-1", "Tennis")

	mockRepo.On("GetByName", mock.Anything, "Tennis").Return(existingSport, nil)

	service := NewSportService(mockRepo)
	sport, err := service.CreateSport(context.Background(), "Tennis")

	assert.Equal(t, domain.ErrDuplicateSportName, err)
	assert.Nil(t, sport)
}

func TestGetSportByID(t *testing.T) {
	mockRepo := new(MockSportRepository)
	expectedSport := domain.NewSport("sport-1", "Tennis")

	mockRepo.On("GetByID", mock.Anything, "sport-1").Return(expectedSport, nil)

	service := NewSportService(mockRepo)
	sport, err := service.GetSportByID(context.Background(), "sport-1")

	assert.NoError(t, err)
	assert.Equal(t, expectedSport, sport)
}

func TestListAllSports(t *testing.T) {
	mockRepo := new(MockSportRepository)
	sports := []*domain.Sport{
		domain.NewSport("sport-1", "Tennis"),
		domain.NewSport("sport-2", "Basketball"),
	}

	mockRepo.On("ListAll", mock.Anything).Return(sports, nil)

	service := NewSportService(mockRepo)
	result, err := service.ListAllSports(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestDeleteSport(t *testing.T) {
	mockRepo := new(MockSportRepository)

	mockRepo.On("CountCourtsBySportID", mock.Anything, "sport-1").Return(0, nil)
	mockRepo.On("Delete", mock.Anything, "sport-1").Return(nil)

	service := NewSportService(mockRepo)
	err := service.DeleteSport(context.Background(), "sport-1")

	assert.NoError(t, err)
}

func TestDeleteSportWithCourts(t *testing.T) {
	mockRepo := new(MockSportRepository)

	mockRepo.On("CountCourtsBySportID", mock.Anything, "sport-1").Return(3, nil)

	service := NewSportService(mockRepo)
	err := service.DeleteSport(context.Background(), "sport-1")

	assert.Equal(t, domain.ErrSportHasCourts, err)
	mockRepo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}
