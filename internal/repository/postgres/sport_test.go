package postgres

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adubr/reservar-club-back/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestSportRepositoryCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewSportRepository(db)
	sport := domain.NewSport("sport-1", "Tennis")

	mock.ExpectExec("INSERT INTO sports").
		WithArgs(sport.ID, sport.Name, sport.CreatedAt, sport.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), sport)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSportRepositoryGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewSportRepository(db)
	sport := domain.NewSport("sport-1", "Tennis")

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow(sport.ID, sport.Name, sport.CreatedAt, sport.UpdatedAt)

	mock.ExpectQuery("SELECT id, name, created_at, updated_at FROM sports WHERE id = ").
		WithArgs(sport.ID).
		WillReturnRows(rows)

	result, err := repo.GetByID(context.Background(), sport.ID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, sport.ID, result.ID)
	assert.Equal(t, sport.Name, result.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSportRepositoryGetByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewSportRepository(db)
	sport := domain.NewSport("sport-1", "Tennis")

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow(sport.ID, sport.Name, sport.CreatedAt, sport.UpdatedAt)

	mock.ExpectQuery("SELECT id, name, created_at, updated_at FROM sports WHERE name = ").
		WithArgs(sport.Name).
		WillReturnRows(rows)

	result, err := repo.GetByName(context.Background(), sport.Name)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, sport.Name, result.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSportRepositoryListAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewSportRepository(db)
	sport1 := domain.NewSport("sport-1", "Tennis")
	sport2 := domain.NewSport("sport-2", "Basketball")

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow(sport1.ID, sport1.Name, sport1.CreatedAt, sport1.UpdatedAt).
		AddRow(sport2.ID, sport2.Name, sport2.CreatedAt, sport2.UpdatedAt)

	mock.ExpectQuery("SELECT id, name, created_at, updated_at FROM sports ORDER BY name").
		WillReturnRows(rows)

	results, err := repo.ListAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSportRepositoryDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewSportRepository(db)

	mock.ExpectExec("DELETE FROM sports WHERE id = ").
		WithArgs("sport-1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(context.Background(), "sport-1")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSportRepositoryCountCourtsBySportID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewSportRepository(db)

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(3)

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM courts WHERE sport_id = ").
		WithArgs("sport-1").
		WillReturnRows(rows)

	count, err := repo.CountCourtsBySportID(context.Background(), "sport-1")

	assert.NoError(t, err)
	assert.Equal(t, 3, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}
