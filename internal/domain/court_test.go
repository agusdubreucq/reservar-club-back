package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCourt(t *testing.T) {
	id := "court-1"
	name := "Court A"
	sportID := "sport-1"

	before := time.Now().UTC()
	court := NewCourt(id, name, sportID)
	after := time.Now().UTC()

	assert.Equal(t, id, court.ID)
	assert.Equal(t, name, court.Name)
	assert.Equal(t, sportID, court.SportID)
	assert.True(t, court.CreatedAt.After(before.Add(-time.Second)))
	assert.True(t, court.CreatedAt.Before(after.Add(time.Second)))
	assert.Equal(t, court.CreatedAt, court.UpdatedAt)
}

func TestCourtUpdate(t *testing.T) {
	court := NewCourt("court-1", "Court A", "sport-1")
	originalCreatedAt := court.CreatedAt

	time.Sleep(10 * time.Millisecond)
	newName := "Court B"
	before := time.Now().UTC()
	court.Update(newName)
	after := time.Now().UTC()

	assert.Equal(t, newName, court.Name)
	assert.Equal(t, originalCreatedAt, court.CreatedAt)
	assert.True(t, court.UpdatedAt.After(before.Add(-time.Second)))
	assert.True(t, court.UpdatedAt.Before(after.Add(time.Second)))
	assert.True(t, court.UpdatedAt.After(court.CreatedAt))
}
