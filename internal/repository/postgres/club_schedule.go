package postgres

import (
	"context"
	"database/sql"

	"github.com/adubr/reservar-club-back/internal/domain"
)

type ClubScheduleRepository struct {
	db *sql.DB
}

func NewClubScheduleRepository(db *sql.DB) *ClubScheduleRepository {
	return &ClubScheduleRepository{db: db}
}

func (r *ClubScheduleRepository) GetByDayOfWeek(ctx context.Context, dayOfWeek int) (*domain.ClubSchedule, error) {
	query := `
		SELECT id, day_of_week, open_time, close_time, is_open, created_at, updated_at
		FROM club_schedules
		WHERE day_of_week = $1
	`

	schedule := &domain.ClubSchedule{}
	err := r.db.QueryRowContext(ctx, query, dayOfWeek).Scan(
		&schedule.ID,
		&schedule.DayOfWeek,
		&schedule.OpenTime,
		&schedule.CloseTime,
		&schedule.IsOpen,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return schedule, nil
}

func (r *ClubScheduleRepository) ListAll(ctx context.Context) ([]*domain.ClubSchedule, error) {
	query := `
		SELECT id, day_of_week, open_time, close_time, is_open, created_at, updated_at
		FROM club_schedules
		ORDER BY day_of_week
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*domain.ClubSchedule
	for rows.Next() {
		schedule := &domain.ClubSchedule{}
		if err := rows.Scan(
			&schedule.ID,
			&schedule.DayOfWeek,
			&schedule.OpenTime,
			&schedule.CloseTime,
			&schedule.IsOpen,
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
		); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *ClubScheduleRepository) Upsert(ctx context.Context, schedule *domain.ClubSchedule) (*domain.ClubSchedule, error) {
	query := `
		INSERT INTO club_schedules (id, day_of_week, open_time, close_time, is_open, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (day_of_week) DO UPDATE
		SET open_time = $3, close_time = $4, is_open = $5, updated_at = $7
		RETURNING id, day_of_week, open_time, close_time, is_open, created_at, updated_at
	`

	result := &domain.ClubSchedule{}
	err := r.db.QueryRowContext(
		ctx,
		query,
		schedule.ID,
		schedule.DayOfWeek,
		schedule.OpenTime,
		schedule.CloseTime,
		schedule.IsOpen,
		schedule.CreatedAt,
		schedule.UpdatedAt,
	).Scan(
		&result.ID,
		&result.DayOfWeek,
		&result.OpenTime,
		&result.CloseTime,
		&result.IsOpen,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
