package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	id := "user-1"
	name := "John Doe"
	email := "john@example.com"
	googleID := "google-123"

	before := time.Now().UTC()
	user := NewUser(id, name, email, googleID)
	after := time.Now().UTC()

	assert.Equal(t, id, user.ID)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, googleID, user.GoogleID)
	assert.True(t, user.CreatedAt.After(before.Add(-time.Second)))
	assert.True(t, user.CreatedAt.Before(after.Add(time.Second)))
	assert.Equal(t, user.CreatedAt, user.UpdatedAt)
}

func TestUserUpdate(t *testing.T) {
	user := NewUser("user-1", "John Doe", "john@example.com", "google-123")
	originalCreatedAt := user.CreatedAt

	time.Sleep(10 * time.Millisecond)
	newName := "Jane Doe"
	before := time.Now().UTC()
	user.Update(newName)
	after := time.Now().UTC()

	assert.Equal(t, newName, user.Name)
	assert.Equal(t, originalCreatedAt, user.CreatedAt)
	assert.True(t, user.UpdatedAt.After(before.Add(-time.Second)))
	assert.True(t, user.UpdatedAt.Before(after.Add(time.Second)))
	assert.True(t, user.UpdatedAt.After(user.CreatedAt))
}
