package utils_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestingEntity struct {
	ID   int
	Name string
}

func SetupGormMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	dbgorm, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Errorf("Failed to open gorm v2 db, got error: %v", err)
	}

	return db, dbgorm, mock
}

func TestPaginationScope(t *testing.T) {
	db, gormdb, mock := SetupGormMock(t)
	defer db.Close()

	rows := mock.NewRows([]string{"id", "name"}).AddRow(1, "test")
	mock.ExpectQuery("LIMIT 1").WillReturnRows(rows)

	var entity TestingEntity
	if err := gormdb.Scopes(utils.Pagination(0, 1)).Find(&entity).Error; err != nil {
		t.Errorf("Failed to query to gorm db, got error: %v", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLikeScope(t *testing.T) {
	db, gormdb, mock := SetupGormMock(t)
	defer db.Close()

	rows := mock.NewRows([]string{"id", "name"}).AddRow(1, "test")
	mock.ExpectQuery("LIKE ?").WithArgs("%meong%").WillReturnRows(rows)

	var entity TestingEntity
	if err := gormdb.Scopes(utils.Like("content", "meong")).Find(&entity).Error; err != nil {
		t.Errorf("Failed to query to gorm db, got error: %v", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
