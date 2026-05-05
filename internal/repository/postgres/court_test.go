package postgres

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adubr/reservar-club-back/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCourtRepositoryCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewCourtRepository(db)
	court := domain.NewCourt("court-1", "Court A", "sport-1")

	mock.ExpectExec("INSERT INTO courts").
		WithArgs(court.ID, court.Name, court.SportID, court.CreatedAt, court.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), court)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCourtRepositoryGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewCourtRepository(db)
	court := domain.NewCourt("court-1", "Court A", "sport-1")

	rows := sqlmock.NewRows([]string{"id", "name", "sport_id", "created_at", "updated_at"}).
		AddRow(court.ID, court.Name, court.SportID, court.CreatedAt, court.UpdatedAt)

	mock.ExpectQuery("SELECT id, name, sport_id, created_at, updated_at FROM courts WHERE id = ").
		WithArgs(court.ID).
		WillReturnRows(rows)

	result, err := repo.GetByID(context.Background(), court.ID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, court.ID, result.ID)
	assert.Equal(t, court.Name, result.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCourtRepositoryListBySportID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewCourtRepository(db)
	court1 := domain.NewCourt("court-1", "Court A", "sport-1")
	court2 := domain.NewCourt("court-2", "Court B", "sport-1")

	rows := sqlmock.NewRows([]string{"id", "name", "sport_id", "created_at", "updated_at"}).
		AddRow(court1.ID, court1.Name, court1.SportID, court1.CreatedAt, court1.UpdatedAt).
		AddRow(court2.ID, court2.Name, court2.SportID, court2.CreatedAt, court2.UpdatedAt)

	mock.ExpectQuery("SELECT id, name, sport_id, created_at, updated_at FROM courts WHERE sport_id = ").
		WithArgs("sport-1").
		WillReturnRows(rows)

	results, err := repo.ListBySportID(context.Background(), "sport-1")

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCourtRepositoryListAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewCourtRepository(db)
	court1 := domain.NewCourt("court-1", "Court A", "sport-1")
	court2 := domain.NewCourt("court-2", "Court B", "sport-2")

	rows := sqlmock.NewRows([]string{"id", "name", "sport_id", "created_at", "updated_at"}).
		AddRow(court1.ID, court1.Name, court1.SportID, court1.CreatedAt, court1.UpdatedAt).
		AddRow(court2.ID, court2.Name, court2.SportID, court2.CreatedAt, court2.UpdatedAt)

	mock.ExpectQuery("SELECT id, name, sport_id, created_at, updated_at FROM courts ORDER BY name").
		WillReturnRows(rows)

	results, err := repo.ListAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCourtRepositoryUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewCourtRepository(db)
	court := domain.NewCourt("court-1", "Court A", "sport-1")
	court.Name = "Court Updated"

	mock.ExpectExec("UPDATE courts SET name = ").
		WithArgs(court.Name, court.SportID, court.UpdatedAt, court.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(context.Background(), court)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCourtRepositoryDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewCourtRepository(db)

	mock.ExpectExec("DELETE FROM courts WHERE id = ").
		WithArgs("court-1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(context.Background(), "court-1")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCourtRepositoryDeleteNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewCourtRepository(db)

	mock.ExpectExec("DELETE FROM courts WHERE id = ").
		WithArgs("nonexistent").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Delete(context.Background(), "nonexistent")

	assert.Equal(t, domain.ErrNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
