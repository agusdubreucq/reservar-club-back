package domain

import "time"

type ReservationStatus string

const (
	ReservationStatusActive    ReservationStatus = "active"
	ReservationStatusCancelled ReservationStatus = "cancelled"
)

type Reservation struct {
	ID         string            `db:"id"`
	UserID     string            `db:"user_id"`
	CourtID    string            `db:"court_id"`
	StartDate  time.Time         `db:"start_date"`
	EndDate    time.Time         `db:"end_date"`
	Status     ReservationStatus `db:"status"`
	TotalPrice float64           `db:"total_price"`
	CreatedAt  time.Time         `db:"created_at"`
	UpdatedAt  time.Time         `db:"updated_at"`
}

func NewReservation(id, userID, courtID string, startDate, endDate time.Time, totalPrice float64) *Reservation {
	now := time.Now().UTC()
	return &Reservation{
		ID:         id,
		UserID:     userID,
		CourtID:    courtID,
		StartDate:  startDate,
		EndDate:    endDate,
		Status:     ReservationStatusActive,
		TotalPrice: totalPrice,
		CreatedAt:  now,
		UpdatedAt:  now,
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

func (r *Reservation) IsValidDuration() bool {
	duration := r.EndDate.Sub(r.StartDate)
	oneHour := 60 * time.Minute
	thirtyMin := 30 * time.Minute

	if duration < oneHour {
		return false
	}

	if duration%thirtyMin != 0 {
		return false
	}

	return true
}

func (r *Reservation) HasValidSlotStartTimes() bool {
	startMin := r.StartDate.Minute()
	endMin := r.EndDate.Minute()

	if startMin != 0 && startMin != 30 {
		return false
	}

	if endMin != 0 && endMin != 30 {
		return false
	}

	return true
}
