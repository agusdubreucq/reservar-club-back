package services

import (
	"context"

	"github.com/adubr/reservar-club-back/internal/domain"
	"github.com/adubr/reservar-club-back/internal/repository"
	"github.com/google/uuid"
)

type ScheduleService struct {
	scheduleRepo repository.ClubScheduleRepository
}

func NewScheduleService(scheduleRepo repository.ClubScheduleRepository) *ScheduleService {
	return &ScheduleService{
		scheduleRepo: scheduleRepo,
	}
}

func (s *ScheduleService) GetSchedule(ctx context.Context) ([]*domain.ClubSchedule, error) {
	return s.scheduleRepo.ListAll(ctx)
}

func (s *ScheduleService) GetScheduleForDay(ctx context.Context, dayOfWeek int) (*domain.ClubSchedule, error) {
	return s.scheduleRepo.GetByDayOfWeek(ctx, dayOfWeek)
}

func (s *ScheduleService) UpsertSchedule(ctx context.Context, dayOfWeek int, openTime, closeTime string, isOpen bool) (*domain.ClubSchedule, error) {
	schedule := domain.NewClubSchedule(uuid.New().String(), dayOfWeek, openTime, closeTime, isOpen)
	return s.scheduleRepo.Upsert(ctx, schedule)
}
