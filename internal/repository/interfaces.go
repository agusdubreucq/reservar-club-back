package repository

import (
	"context"
	"time"

	"github.com/adubr/reservar-club-back/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByGoogleID(ctx context.Context, googleID string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

type SportRepository interface {
	Create(ctx context.Context, sport *domain.Sport) error
	GetByID(ctx context.Context, id string) (*domain.Sport, error)
	GetByName(ctx context.Context, name string) (*domain.Sport, error)
	ListAll(ctx context.Context) ([]*domain.Sport, error)
	Delete(ctx context.Context, id string) error
	CountCourtsBySportID(ctx context.Context, sportID string) (int, error)
}

type CourtRepository interface {
	Create(ctx context.Context, court *domain.Court) error
	GetByID(ctx context.Context, id string) (*domain.Court, error)
	ListBySportID(ctx context.Context, sportID string) ([]*domain.Court, error)
	ListAll(ctx context.Context) ([]*domain.Court, error)
	Update(ctx context.Context, court *domain.Court) error
	Delete(ctx context.Context, id string) error
}

type ReservationRepository interface {
	Create(ctx context.Context, reservation *domain.Reservation) error
	GetByID(ctx context.Context, id string) (*domain.Reservation, error)
	ListByUserID(ctx context.Context, userID string) ([]*domain.Reservation, error)
	ListByCourtAndDateRange(ctx context.Context, courtID string, startDate, endDate time.Time) ([]*domain.Reservation, error)
	Cancel(ctx context.Context, id string) error
}
