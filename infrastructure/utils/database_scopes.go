package utils

import (
	"gorm.io/gorm"
)

func Pagination(page int, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

func Like(field []string, value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, v := range field {
			db = db.Or(v+" LIKE ?", "%"+value+"%")
		}

		return db
	}
}
