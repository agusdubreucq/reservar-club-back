package domain

import "time"

type PricingRule struct {
	ID           string     `db:"id"`
	SportID      string     `db:"sport_id"`
	DayOfWeek    *int       `db:"day_of_week"`
	StartHour    int        `db:"start_hour"`
	EndHour      int        `db:"end_hour"`
	PricePerHour float64    `db:"price_per_hour"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}

func NewPricingRule(id, sportID string, dayOfWeek *int, startHour, endHour int, pricePerHour float64) *PricingRule {
	now := time.Now().UTC()
	return &PricingRule{
		ID:           id,
		SportID:      sportID,
		DayOfWeek:    dayOfWeek,
		StartHour:    startHour,
		EndHour:      endHour,
		PricePerHour: pricePerHour,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (pr *PricingRule) Update(startHour, endHour int, pricePerHour float64) {
	pr.StartHour = startHour
	pr.EndHour = endHour
	pr.PricePerHour = pricePerHour
	pr.UpdatedAt = time.Now().UTC()
}
