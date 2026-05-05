package postgres

import (
	"context"
	"database/sql"

	"github.com/adubr/reservar-club-back/internal/domain"
)

type CourtRepository struct {
	db *sql.DB
}

func NewCourtRepository(db *sql.DB) *CourtRepository {
	return &CourtRepository{db: db}
}

func (r *CourtRepository) Create(ctx context.Context, court *domain.Court) error {
	query := `
		INSERT INTO courts (id, name, sport_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query, court.ID, court.Name, court.SportID, court.CreatedAt, court.UpdatedAt)
	return err
}

func (r *CourtRepository) GetByID(ctx context.Context, id string) (*domain.Court, error) {
	query := `
		SELECT id, name, sport_id, created_at, updated_at
		FROM courts
		WHERE id = $1
	`

	court := &domain.Court{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&court.ID,
		&court.Name,
		&court.SportID,
		&court.CreatedAt,
		&court.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return court, nil
}

func (r *CourtRepository) ListBySportID(ctx context.Context, sportID string) ([]*domain.Court, error) {
	query := `
		SELECT id, name, sport_id, created_at, updated_at
		FROM courts
		WHERE sport_id = $1
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query, sportID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courts []*domain.Court
	for rows.Next() {
		court := &domain.Court{}
		if err := rows.Scan(&court.ID, &court.Name, &court.SportID, &court.CreatedAt, &court.UpdatedAt); err != nil {
			return nil, err
		}
		courts = append(courts, court)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courts, nil
}

func (r *CourtRepository) ListAll(ctx context.Context) ([]*domain.Court, error) {
	query := `
		SELECT id, name, sport_id, created_at, updated_at
		FROM courts
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courts []*domain.Court
	for rows.Next() {
		court := &domain.Court{}
		if err := rows.Scan(&court.ID, &court.Name, &court.SportID, &court.CreatedAt, &court.UpdatedAt); err != nil {
			return nil, err
		}
		courts = append(courts, court)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courts, nil
}

func (r *CourtRepository) Update(ctx context.Context, court *domain.Court) error {
	query := `
		UPDATE courts
		SET name = $1, sport_id = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx, query, court.Name, court.SportID, court.UpdatedAt, court.ID)
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

func (r *CourtRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM courts WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
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
