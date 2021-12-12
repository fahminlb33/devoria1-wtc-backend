package users

import (
	"gorm.io/gorm"
)

type UserRole string

const (
	ADMIN       UserRole = "ADMIN"
	CONTRIBUTOR UserRole = "CONTRIBUTOR"
)

type User struct {
	gorm.Model
	Email     string   `gorm:"size:255;not null;uniqueIndex:idx_email"`
	Password  *string  `gorm:"size:255;not null"`
	FirstName string   `gorm:"size:255;not null"`
	LastName  string   `gorm:"size:255;not null"`
	Role      UserRole `gorm:"size:255;not null"`
}
