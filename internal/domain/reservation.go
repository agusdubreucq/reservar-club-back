package domain

import "time"

type ReservationStatus string

const (
	ReservationStatusActive    ReservationStatus = "active"
	ReservationStatusCancelled ReservationStatus = "cancelled"
)

type Reservation struct {
	ID        string            `db:"id"`
	UserID    string            `db:"user_id"`
	CourtID   string            `db:"court_id"`
	StartDate time.Time         `db:"start_date"`
	EndDate   time.Time         `db:"end_date"`
	Status    ReservationStatus `db:"status"`
	CreatedAt time.Time         `db:"created_at"`
	UpdatedAt time.Time         `db:"updated_at"`
}

func NewReservation(id, userID, courtID string, startDate, endDate time.Time) *Reservation {
	now := time.Now().UTC()
	return &Reservation{
		ID:        id,
		UserID:    userID,
		CourtID:   courtID,
		StartDate: startDate,
		EndDate:   endDate,
		Status:    ReservationStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (r *Reservation) Cancel() {
	r.Status = ReservationStatusCancelled
	r.UpdatedAt = time.Now().UTC()
}

func (r *Reservation) IsOverlapping(other *Reservation) bool {
	return r.StartDate.Before(other.EndDate) && r.EndDate.After(other.StartDate)
}

func (r *Reservation) IsValidDateRange() bool {
	return r.EndDate.After(r.StartDate)
}
