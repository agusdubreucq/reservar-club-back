package services

import (
	"context"
	"testing"

	"github.com/adubr/reservar-club-back/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetOrCreateUserByGoogle_UserExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	existingUser := domain.NewUser("user-1", "John Doe", "john@example.com", "google-123")

	mockRepo.On("GetByGoogleID", mock.Anything, "google-123").
		Return(existingUser, nil)

	service := NewAuthService(mockRepo)
	user, err := service.GetOrCreateUserByGoogle(context.Background(), "john@example.com", "google-123", "John Doe")

	assert.NoError(t, err)
	assert.Equal(t, existingUser, user)
	mockRepo.AssertCalled(t, "GetByGoogleID", mock.Anything, "google-123")
	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestGetOrCreateUserByGoogle_UserNotExists(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetByGoogleID", mock.Anything, "google-123").
		Return(nil, domain.ErrNotFound)

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Email == "john@example.com" && u.GoogleID == "google-123" && u.Name == "John Doe"
	})).Return(nil)

	service := NewAuthService(mockRepo)
	user, err := service.GetOrCreateUserByGoogle(context.Background(), "john@example.com", "google-123", "John Doe")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, "google-123", user.GoogleID)
	mockRepo.AssertCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestGetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	expectedUser := domain.NewUser("user-1", "John Doe", "john@example.com", "google-123")

	mockRepo.On("GetByID", mock.Anything, "user-1").
		Return(expectedUser, nil)

	service := NewAuthService(mockRepo)
	user, err := service.GetUserByID(context.Background(), "user-1")

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertCalled(t, "GetByID", mock.Anything, "user-1")
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetByID", mock.Anything, "nonexistent").
		Return(nil, domain.ErrNotFound)

	service := NewAuthService(mockRepo)
	user, err := service.GetUserByID(context.Background(), "nonexistent")

	assert.Equal(t, domain.ErrNotFound, err)
	assert.Nil(t, user)
}
