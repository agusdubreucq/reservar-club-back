package postgres

import (
	"context"
	"database/sql"

	"github.com/adubr/reservar-club-back/internal/domain"
)

type SportRepository struct {
	db *sql.DB
}

func NewSportRepository(db *sql.DB) *SportRepository {
	return &SportRepository{db: db}
}

func (r *SportRepository) Create(ctx context.Context, sport *domain.Sport) error {
	query := `
		INSERT INTO sports (id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query, sport.ID, sport.Name, sport.CreatedAt, sport.UpdatedAt)
	return err
}

func (r *SportRepository) GetByID(ctx context.Context, id string) (*domain.Sport, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM sports
		WHERE id = $1
	`

	sport := &domain.Sport{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&sport.ID,
		&sport.Name,
		&sport.CreatedAt,
		&sport.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return sport, nil
}

func (r *SportRepository) GetByName(ctx context.Context, name string) (*domain.Sport, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM sports
		WHERE name = $1
	`

	sport := &domain.Sport{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&sport.ID,
		&sport.Name,
		&sport.CreatedAt,
		&sport.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return sport, nil
}

func (r *SportRepository) ListAll(ctx context.Context) ([]*domain.Sport, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM sports
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sports []*domain.Sport
	for rows.Next() {
		sport := &domain.Sport{}
		if err := rows.Scan(&sport.ID, &sport.Name, &sport.CreatedAt, &sport.UpdatedAt); err != nil {
			return nil, err
		}
		sports = append(sports, sport)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sports, nil
}

func (r *SportRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM sports WHERE id = $1`

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

func (r *SportRepository) CountCourtsBySportID(ctx context.Context, sportID string) (int, error) {
	query := `SELECT COUNT(*) FROM courts WHERE sport_id = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, sportID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
