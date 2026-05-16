package services

import (
	"context"
	"time"

	"github.com/adubr/reservar-club-back/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {
	args := m.Called(ctx, googleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

type MockSportRepository struct {
	mock.Mock
}

func (m *MockSportRepository) Create(ctx context.Context, sport *domain.Sport) error {
	args := m.Called(ctx, sport)
	return args.Error(0)
}

func (m *MockSportRepository) GetByID(ctx context.Context, id string) (*domain.Sport, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Sport), args.Error(1)
}

func (m *MockSportRepository) GetByName(ctx context.Context, name string) (*domain.Sport, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Sport), args.Error(1)
}

func (m *MockSportRepository) ListAll(ctx context.Context) ([]*domain.Sport, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Sport), args.Error(1)
}

func (m *MockSportRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSportRepository) CountCourtsBySportID(ctx context.Context, sportID string) (int, error) {
	args := m.Called(ctx, sportID)
	return args.Int(0), args.Error(1)
}

type MockCourtRepository struct {
	mock.Mock
}

func (m *MockCourtRepository) Create(ctx context.Context, court *domain.Court) error {
	args := m.Called(ctx, court)
	return args.Error(0)
}

func (m *MockCourtRepository) GetByID(ctx context.Context, id string) (*domain.Court, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Court), args.Error(1)
}

func (m *MockCourtRepository) ListBySportID(ctx context.Context, sportID string) ([]*domain.Court, error) {
	args := m.Called(ctx, sportID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Court), args.Error(1)
}

func (m *MockCourtRepository) ListAll(ctx context.Context) ([]*domain.Court, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Court), args.Error(1)
}

func (m *MockCourtRepository) Update(ctx context.Context, court *domain.Court) error {
	args := m.Called(ctx, court)
	return args.Error(0)
}

func (m *MockCourtRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockReservationRepository struct {
	mock.Mock
}

func (m *MockReservationRepository) Create(ctx context.Context, reservation *domain.Reservation) error {
	args := m.Called(ctx, reservation)
	return args.Error(0)
}

func (m *MockReservationRepository) GetByID(ctx context.Context, id string) (*domain.Reservation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Reservation), args.Error(1)
}

func (m *MockReservationRepository) ListByUserID(ctx context.Context, userID string) ([]*domain.Reservation, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Reservation), args.Error(1)
}

func (m *MockReservationRepository) ListByCourtAndDateRange(ctx context.Context, courtID string, startDate, endDate time.Time) ([]*domain.Reservation, error) {
	args := m.Called(ctx, courtID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Reservation), args.Error(1)
}

func (m *MockReservationRepository) Cancel(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockClubScheduleRepository struct {
	mock.Mock
}

func (m *MockClubScheduleRepository) GetByDayOfWeek(ctx context.Context, dayOfWeek int) (*domain.ClubSchedule, error) {
	args := m.Called(ctx, dayOfWeek)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ClubSchedule), args.Error(1)
}

func (m *MockClubScheduleRepository) ListAll(ctx context.Context) ([]*domain.ClubSchedule, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ClubSchedule), args.Error(1)
}

func (m *MockClubScheduleRepository) Upsert(ctx context.Context, schedule *domain.ClubSchedule) (*domain.ClubSchedule, error) {
	args := m.Called(ctx, schedule)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ClubSchedule), args.Error(1)
}

type MockPricingRuleRepository struct {
	mock.Mock
}

func (m *MockPricingRuleRepository) Create(ctx context.Context, rule *domain.PricingRule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *MockPricingRuleRepository) GetByID(ctx context.Context, id string) (*domain.PricingRule, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PricingRule), args.Error(1)
}

func (m *MockPricingRuleRepository) ListBySportID(ctx context.Context, sportID string) ([]*domain.PricingRule, error) {
	args := m.Called(ctx, sportID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PricingRule), args.Error(1)
}

func (m *MockPricingRuleRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
