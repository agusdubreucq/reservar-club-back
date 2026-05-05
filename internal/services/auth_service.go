package services

import (
	"context"
	"errors"

	"github.com/adubr/reservar-club-back/internal/domain"
	"github.com/adubr/reservar-club-back/internal/repository"
	"github.com/google/uuid"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) GetOrCreateUserByGoogle(ctx context.Context, email, googleID, name string) (*domain.User, error) {
	user, err := s.userRepo.GetByGoogleID(ctx, googleID)
	if err == nil && user != nil {
		return user, nil
	}

	if !errors.Is(err, domain.ErrNotFound) && err != nil {
		return nil, err
	}

	user = domain.NewUser(uuid.New().String(), name, email, googleID)
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}
