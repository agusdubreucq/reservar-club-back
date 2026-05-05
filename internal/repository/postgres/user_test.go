package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adubr/reservar-club-back/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestUserRepositoryCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	user := domain.NewUser("user-1", "John Doe", "john@example.com", "google-123")

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.ID, user.Name, user.Email, user.GoogleID, user.CreatedAt, user.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), user)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	user := domain.NewUser("user-1", "John Doe", "john@example.com", "google-123")

	rows := sqlmock.NewRows([]string{"id", "name", "email", "google_id", "created_at", "updated_at"}).
		AddRow(user.ID, user.Name, user.Email, user.GoogleID, user.CreatedAt, user.UpdatedAt)

	mock.ExpectQuery("SELECT id, name, email, google_id, created_at, updated_at FROM users WHERE id = ").
		WithArgs(user.ID).
		WillReturnRows(rows)

	result, err := repo.GetByID(context.Background(), user.ID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Name, result.Name)
	assert.Equal(t, user.Email, result.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryGetByIDNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectQuery("SELECT id, name, email, google_id, created_at, updated_at FROM users WHERE id = ").
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	result, err := repo.GetByID(context.Background(), "nonexistent")

	assert.Equal(t, domain.ErrNotFound, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryGetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	user := domain.NewUser("user-1", "John Doe", "john@example.com", "google-123")

	rows := sqlmock.NewRows([]string{"id", "name", "email", "google_id", "created_at", "updated_at"}).
		AddRow(user.ID, user.Name, user.Email, user.GoogleID, user.CreatedAt, user.UpdatedAt)

	mock.ExpectQuery("SELECT id, name, email, google_id, created_at, updated_at FROM users WHERE email = ").
		WithArgs(user.Email).
		WillReturnRows(rows)

	result, err := repo.GetByEmail(context.Background(), user.Email)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.Email, result.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryGetByGoogleID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	user := domain.NewUser("user-1", "John Doe", "john@example.com", "google-123")

	rows := sqlmock.NewRows([]string{"id", "name", "email", "google_id", "created_at", "updated_at"}).
		AddRow(user.ID, user.Name, user.Email, user.GoogleID, user.CreatedAt, user.UpdatedAt)

	mock.ExpectQuery("SELECT id, name, email, google_id, created_at, updated_at FROM users WHERE google_id = ").
		WithArgs(user.GoogleID).
		WillReturnRows(rows)

	result, err := repo.GetByGoogleID(context.Background(), user.GoogleID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.GoogleID, result.GoogleID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	user := domain.NewUser("user-1", "John Doe", "john@example.com", "google-123")
	user.Name = "Jane Doe"
	user.UpdatedAt = time.Now().UTC()

	mock.ExpectExec("UPDATE users SET name = ").
		WithArgs(user.Name, user.Email, user.GoogleID, user.UpdatedAt, user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(context.Background(), user)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryUpdateNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	user := domain.NewUser("user-1", "John Doe", "john@example.com", "google-123")

	mock.ExpectExec("UPDATE users SET name = ").
		WithArgs(user.Name, user.Email, user.GoogleID, user.UpdatedAt, user.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Update(context.Background(), user)

	assert.Equal(t, domain.ErrNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
