package domain

import "time"

type User struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	GoogleID  string    `db:"google_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewUser(id, name, email, googleID string) *User {
	now := time.Now().UTC()
	return &User{
		ID:        id,
		Name:      name,
		Email:     email,
		GoogleID:  googleID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u *User) Update(name string) {
	u.Name = name
	u.UpdatedAt = time.Now().UTC()
}
