package services

import (
	"context"

	"github.com/adubr/reservar-club-back/internal/domain"
	"github.com/adubr/reservar-club-back/internal/repository"
	"github.com/google/uuid"
)

type CourtService struct {
	courtRepo repository.CourtRepository
	sportRepo repository.SportRepository
}

func NewCourtService(courtRepo repository.CourtRepository, sportRepo repository.SportRepository) *CourtService {
	return &CourtService{courtRepo: courtRepo, sportRepo: sportRepo}
}

func (s *CourtService) CreateCourt(ctx context.Context, name, sportID string) (*domain.Court, error) {
	sport, err := s.sportRepo.GetByID(ctx, sportID)
	if err != nil {
		return nil, domain.ErrSportNotFound
	}

	if sport == nil {
		return nil, domain.ErrSportNotFound
	}

	court := domain.NewCourt(uuid.New().String(), name, sportID)
	if err := s.courtRepo.Create(ctx, court); err != nil {
		return nil, err
	}

	return court, nil
}

func (s *CourtService) GetCourtByID(ctx context.Context, id string) (*domain.Court, error) {
	return s.courtRepo.GetByID(ctx, id)
}

func (s *CourtService) ListAllCourts(ctx context.Context) ([]*domain.Court, error) {
	return s.courtRepo.ListAll(ctx)
}

func (s *CourtService) ListCourtsBySport(ctx context.Context, sportID string) ([]*domain.Court, error) {
	return s.courtRepo.ListBySportID(ctx, sportID)
}

func (s *CourtService) UpdateCourt(ctx context.Context, id, name string) (*domain.Court, error) {
	court, err := s.courtRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrCourtNotFound
	}

	if court == nil {
		return nil, domain.ErrCourtNotFound
	}

	court.Update(name)
	if err := s.courtRepo.Update(ctx, court); err != nil {
		return nil, err
	}

	return court, nil
}

func (s *CourtService) DeleteCourt(ctx context.Context, id string) error {
	return s.courtRepo.Delete(ctx, id)
}
