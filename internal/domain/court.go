package domain

import "time"

type Court struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	SportID   string    `db:"sport_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewCourt(id, name, sportID string) *Court {
	now := time.Now().UTC()
	return &Court{
		ID:        id,
		Name:      name,
		SportID:   sportID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (c *Court) Update(name string) {
	c.Name = name
	c.UpdatedAt = time.Now().UTC()
}
