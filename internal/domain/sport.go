package domain

import "time"

type Sport struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewSport(id, name string) *Sport {
	now := time.Now().UTC()
	return &Sport{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (s *Sport) Update(name string) {
	s.Name = name
	s.UpdatedAt = time.Now().UTC()
}
