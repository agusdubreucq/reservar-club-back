package handlers

import (
	"context"
	"testing"
	"time"

	"github.com/adubr/reservar-club-back/internal/domain"
	"github.com/adubr/reservar-club-back/internal/services"
	"github.com/adubr/reservar-club-back/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepo) GetByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {
	args := m.Called(ctx, googleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepo) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func TestTokenManager(t *testing.T) {
	tm := jwt.NewTokenManager("test-secret")
	token, err := tm.GenerateToken("user-1", "john@example.com", 24*time.Hour)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := tm.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "user-1", claims.UserID)
	assert.Equal(t, "john@example.com", claims.Email)
}

func TestTokenManagerInvalidToken(t *testing.T) {
	tm := jwt.NewTokenManager("test-secret")

	claims, err := tm.ValidateToken("invalid-token")

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestTokenManagerExpiredToken(t *testing.T) {
	tm := jwt.NewTokenManager("test-secret")
	token, _ := tm.GenerateToken("user-1", "john@example.com", -time.Second)

	time.Sleep(100 * time.Millisecond)
	claims, err := tm.ValidateToken(token)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestAuthServiceIntegration(t *testing.T) {
	mockRepo := new(mockUserRepo)
	mockRepo.On("GetByGoogleID", mock.Anything, "google-123").
		Return(nil, domain.ErrNotFound)
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Email == "john@example.com"
	})).Return(nil)

	service := services.NewAuthService(mockRepo)
	user, err := service.GetOrCreateUserByGoogle(context.Background(), "john@example.com", "google-123", "John Doe")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "john@example.com", user.Email)
}
