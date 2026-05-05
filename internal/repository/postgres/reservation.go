package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/adubr/reservar-club-back/internal/domain"
)

type ReservationRepository struct {
	db *sql.DB
}

func NewReservationRepository(db *sql.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (r *ReservationRepository) Create(ctx context.Context, reservation *domain.Reservation) error {
	query := `
		INSERT INTO reservations (id, user_id, court_id, start_date, end_date, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		reservation.ID,
		reservation.UserID,
		reservation.CourtID,
		reservation.StartDate,
		reservation.EndDate,
		reservation.Status,
		reservation.CreatedAt,
		reservation.UpdatedAt,
	)
	return err
}

func (r *ReservationRepository) GetByID(ctx context.Context, id string) (*domain.Reservation, error) {
	query := `
		SELECT id, user_id, court_id, start_date, end_date, status, created_at, updated_at
		FROM reservations
		WHERE id = $1
	`

	reservation := &domain.Reservation{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&reservation.ID,
		&reservation.UserID,
		&reservation.CourtID,
		&reservation.StartDate,
		&reservation.EndDate,
		&reservation.Status,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (r *ReservationRepository) ListByUserID(ctx context.Context, userID string) ([]*domain.Reservation, error) {
	query := `
		SELECT id, user_id, court_id, start_date, end_date, status, created_at, updated_at
		FROM reservations
		WHERE user_id = $1
		ORDER BY start_date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*domain.Reservation
	for rows.Next() {
		reservation := &domain.Reservation{}
		if err := rows.Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.CourtID,
			&reservation.StartDate,
			&reservation.EndDate,
			&reservation.Status,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
		); err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}

func (r *ReservationRepository) ListByCourtAndDateRange(ctx context.Context, courtID string, startDate, endDate time.Time) ([]*domain.Reservation, error) {
	query := `
		SELECT id, user_id, court_id, start_date, end_date, status, created_at, updated_at
		FROM reservations
		WHERE court_id = $1
		AND start_date < $3
		AND end_date > $2
		ORDER BY start_date
	`

	rows, err := r.db.QueryContext(ctx, query, courtID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*domain.Reservation
	for rows.Next() {
		reservation := &domain.Reservation{}
		if err := rows.Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.CourtID,
			&reservation.StartDate,
			&reservation.EndDate,
			&reservation.Status,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
		); err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}

func (r *ReservationRepository) Cancel(ctx context.Context, id string) error {
	query := `
		UPDATE reservations
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, domain.ReservationStatusCancelled, time.Now().UTC(), id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}
