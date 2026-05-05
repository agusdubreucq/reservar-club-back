package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSport(t *testing.T) {
	id := "sport-1"
	name := "Tennis"

	before := time.Now().UTC()
	sport := NewSport(id, name)
	after := time.Now().UTC()

	assert.Equal(t, id, sport.ID)
	assert.Equal(t, name, sport.Name)
	assert.True(t, sport.CreatedAt.After(before.Add(-time.Second)))
	assert.True(t, sport.CreatedAt.Before(after.Add(time.Second)))
	assert.Equal(t, sport.CreatedAt, sport.UpdatedAt)
}

func TestSportUpdate(t *testing.T) {
	sport := NewSport("sport-1", "Tennis")
	originalCreatedAt := sport.CreatedAt

	time.Sleep(10 * time.Millisecond)
	newName := "Basketball"
	before := time.Now().UTC()
	sport.Update(newName)
	after := time.Now().UTC()

	assert.Equal(t, newName, sport.Name)
	assert.Equal(t, originalCreatedAt, sport.CreatedAt)
	assert.True(t, sport.UpdatedAt.After(before.Add(-time.Second)))
	assert.True(t, sport.UpdatedAt.Before(after.Add(time.Second)))
	assert.True(t, sport.UpdatedAt.After(sport.CreatedAt))
}
