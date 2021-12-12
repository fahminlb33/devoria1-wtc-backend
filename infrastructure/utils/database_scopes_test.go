package utils_test

import (
	"testing"

	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/utils"
	"github.com/fahminlb33/devoria1-wtc-backend/mocks"
)

type TestingEntity struct {
	ID   int
	Name string
}

func TestPaginationScope(t *testing.T) {
	db, gormdb, mock := mocks.SetupGormMock(t)
	defer db.Close()

	rows := mock.NewRows([]string{"id", "name"}).AddRow(1, "test")
	mock.ExpectQuery("LIMIT 1").WillReturnRows(rows)

	var entity TestingEntity
	if err := gormdb.Scopes(utils.Pagination(0, 1)).Find(&entity).Error; err != nil {
		t.Errorf("Failed to query to gorm db, got error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLikeScope(t *testing.T) {
	db, gormdb, mock := mocks.SetupGormMock(t)
	defer db.Close()

	rows := mock.NewRows([]string{"id", "name"}).AddRow(1, "test")
	mock.ExpectQuery("LIKE ?").WithArgs("%meong%").WillReturnRows(rows)

	var entity TestingEntity
	if err := gormdb.Scopes(utils.Like([]string{"content"}, "meong")).Find(&entity).Error; err != nil {
		t.Errorf("Failed to query to gorm db, got error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
