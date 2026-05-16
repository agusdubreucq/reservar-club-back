package domain

import "time"

type ClubSchedule struct {
	ID        string    `db:"id"`
	DayOfWeek int       `db:"day_of_week"`
	OpenTime  string    `db:"open_time"`
	CloseTime string    `db:"close_time"`
	IsOpen    bool      `db:"is_open"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewClubSchedule(id string, dayOfWeek int, openTime, closeTime string, isOpen bool) *ClubSchedule {
	now := time.Now().UTC()
	return &ClubSchedule{
		ID:        id,
		DayOfWeek: dayOfWeek,
		OpenTime:  openTime,
		CloseTime: closeTime,
		IsOpen:    isOpen,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (cs *ClubSchedule) Update(openTime, closeTime string, isOpen bool) {
	cs.OpenTime = openTime
	cs.CloseTime = closeTime
	cs.IsOpen = isOpen
	cs.UpdatedAt = time.Now().UTC()
}
