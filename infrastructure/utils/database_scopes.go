package utils

import (
	"gorm.io/gorm"
)

func Pagination(page int, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

func Like(field string, value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" LIKE ?", "%"+value+"%")
	}
}
