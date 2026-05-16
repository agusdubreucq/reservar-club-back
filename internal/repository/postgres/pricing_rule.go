package postgres

import (
	"context"
	"database/sql"

	"github.com/adubr/reservar-club-back/internal/domain"
)

type PricingRuleRepository struct {
	db *sql.DB
}

func NewPricingRuleRepository(db *sql.DB) *PricingRuleRepository {
	return &PricingRuleRepository{db: db}
}

func (r *PricingRuleRepository) Create(ctx context.Context, rule *domain.PricingRule) error {
	query := `
		INSERT INTO pricing_rules (id, sport_id, day_of_week, start_hour, end_hour, price_per_hour, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		rule.ID,
		rule.SportID,
		rule.DayOfWeek,
		rule.StartHour,
		rule.EndHour,
		rule.PricePerHour,
		rule.CreatedAt,
		rule.UpdatedAt,
	)
	return err
}

func (r *PricingRuleRepository) GetByID(ctx context.Context, id string) (*domain.PricingRule, error) {
	query := `
		SELECT id, sport_id, day_of_week, start_hour, end_hour, price_per_hour, created_at, updated_at
		FROM pricing_rules
		WHERE id = $1
	`

	rule := &domain.PricingRule{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rule.ID,
		&rule.SportID,
		&rule.DayOfWeek,
		&rule.StartHour,
		&rule.EndHour,
		&rule.PricePerHour,
		&rule.CreatedAt,
		&rule.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return rule, nil
}

func (r *PricingRuleRepository) ListBySportID(ctx context.Context, sportID string) ([]*domain.PricingRule, error) {
	query := `
		SELECT id, sport_id, day_of_week, start_hour, end_hour, price_per_hour, created_at, updated_at
		FROM pricing_rules
		WHERE sport_id = $1
		ORDER BY day_of_week, start_hour
	`

	rows, err := r.db.QueryContext(ctx, query, sportID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*domain.PricingRule
	for rows.Next() {
		rule := &domain.PricingRule{}
		if err := rows.Scan(
			&rule.ID,
			&rule.SportID,
			&rule.DayOfWeek,
			&rule.StartHour,
			&rule.EndHour,
			&rule.PricePerHour,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		); err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}

func (r *PricingRuleRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM pricing_rules WHERE id = $1`

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
