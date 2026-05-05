package services

import (
	"context"

	"github.com/adubr/reservar-club-back/internal/domain"
	"github.com/adubr/reservar-club-back/internal/repository"
	"github.com/google/uuid"
)

type SportService struct {
	sportRepo repository.SportRepository
}

func NewSportService(sportRepo repository.SportRepository) *SportService {
	return &SportService{sportRepo: sportRepo}
}

func (s *SportService) CreateSport(ctx context.Context, name string) (*domain.Sport, error) {
	existing, err := s.sportRepo.GetByName(ctx, name)
	if err == nil && existing != nil {
		return nil, domain.ErrDuplicateSportName
	}

	sport := domain.NewSport(uuid.New().String(), name)
	if err := s.sportRepo.Create(ctx, sport); err != nil {
		return nil, err
	}

	return sport, nil
}

func (s *SportService) GetSportByID(ctx context.Context, id string) (*domain.Sport, error) {
	return s.sportRepo.GetByID(ctx, id)
}

func (s *SportService) ListAllSports(ctx context.Context) ([]*domain.Sport, error) {
	return s.sportRepo.ListAll(ctx)
}

func (s *SportService) DeleteSport(ctx context.Context, id string) error {
	count, err := s.sportRepo.CountCourtsBySportID(ctx, id)
	if err != nil {
		return err
	}

	if count > 0 {
		return domain.ErrSportHasCourts
	}

	return s.sportRepo.Delete(ctx, id)
}
