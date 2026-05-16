package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewReservation(t *testing.T) {
	id := "res-1"
	userID := "user-1"
	courtID := "court-1"
	startDate := time.Now().UTC().Add(24 * time.Hour)
	endDate := startDate.Add(2 * time.Hour)
	price := 100.0

	res := NewReservation(id, userID, courtID, startDate, endDate, price)

	assert.Equal(t, id, res.ID)
	assert.Equal(t, userID, res.UserID)
	assert.Equal(t, courtID, res.CourtID)
	assert.Equal(t, startDate, res.StartDate)
	assert.Equal(t, endDate, res.EndDate)
	assert.Equal(t, ReservationStatusActive, res.Status)
	assert.Equal(t, price, res.TotalPrice)
	assert.Equal(t, res.CreatedAt, res.UpdatedAt)
}

func TestReservationCancel(t *testing.T) {
	res := NewReservation("res-1", "user-1", "court-1",
		time.Now().UTC().Add(24*time.Hour),
		time.Now().UTC().Add(26*time.Hour), 50.0)

	originalCreatedAt := res.CreatedAt

	time.Sleep(10 * time.Millisecond)
	before := time.Now().UTC()
	res.Cancel()
	after := time.Now().UTC()

	assert.Equal(t, ReservationStatusCancelled, res.Status)
	assert.Equal(t, originalCreatedAt, res.CreatedAt)
	assert.True(t, res.UpdatedAt.After(before.Add(-time.Second)))
	assert.True(t, res.UpdatedAt.Before(after.Add(time.Second)))
}

func TestReservationIsOverlapping(t *testing.T) {
	startDate := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	res1 := NewReservation("res-1", "user-1", "court-1", startDate, endDate, 0)

	tests := []struct {
		name      string
		startDate time.Time
		endDate   time.Time
		expected  bool
	}{
		{
			name:      "no overlap - before",
			startDate: time.Date(2025, 1, 1, 8, 0, 0, 0, time.UTC),
			endDate:   time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			expected:  false,
		},
		{
			name:      "no overlap - after",
			startDate: time.Date(2025, 1, 1, 13, 0, 0, 0, time.UTC),
			endDate:   time.Date(2025, 1, 1, 14, 0, 0, 0, time.UTC),
			expected:  false,
		},
		{
			name:      "overlap - inside",
			startDate: time.Date(2025, 1, 1, 10, 30, 0, 0, time.UTC),
			endDate:   time.Date(2025, 1, 1, 11, 30, 0, 0, time.UTC),
			expected:  true,
		},
		{
			name:      "overlap - partial start",
			startDate: time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			endDate:   time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC),
			expected:  true,
		},
		{
			name:      "overlap - partial end",
			startDate: time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC),
			endDate:   time.Date(2025, 1, 1, 13, 0, 0, 0, time.UTC),
			expected:  true,
		},
		{
			name:      "overlap - covers all",
			startDate: time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
			endDate:   time.Date(2025, 1, 1, 13, 0, 0, 0, time.UTC),
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res2 := NewReservation("res-2", "user-1", "court-1", tt.startDate, tt.endDate, 0)
			assert.Equal(t, tt.expected, res1.IsOverlapping(res2))
		})
	}
}

func TestReservationIsValidDateRange(t *testing.T) {
	tests := []struct {
		name      string
		startDate time.Time
		endDate   time.Time
		expected  bool
	}{
		{
			name:      "valid - end after start",
			startDate: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
			endDate:   time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
			expected:  true,
		},
		{
			name:      "invalid - end before start",
			startDate: time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
			endDate:   time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
			expected:  false,
		},
		{
			name:      "invalid - same time",
			startDate: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
			endDate:   time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := NewReservation("res-1", "user-1", "court-1", tt.startDate, tt.endDate, 0)
			assert.Equal(t, tt.expected, res.IsValidDateRange())
		})
	}
}
