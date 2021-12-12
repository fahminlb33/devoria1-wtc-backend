package database

import (
	"errors"
	"log"

	postgres "go.elastic.co/apm/module/apmgormv2/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/fahminlb33/devoria1-wtc-backend/domain/articles"
	"github.com/fahminlb33/devoria1-wtc-backend/domain/users"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/authentication"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/config"
)

func Initialize() (*gorm.DB, error) {
	log.Println("Opening connection to database...")
	db, err := gorm.Open(postgres.Open(config.GlobalConfig.Database.URI), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return db, err
	}

	return db, nil
}

func MigrateIfNeeded(db *gorm.DB) error {
	log.Println("Running database migration if necessary...")
	err := db.AutoMigrate(&users.User{}, &articles.Article{})
	if err != nil {
		return err
	}

	return nil
}

func SeedIfNeeded(db *gorm.DB) error {
	if db.Migrator().HasTable(&users.User{}) {
		if migrateErr := db.First(&users.User{}).Error; errors.Is(migrateErr, gorm.ErrRecordNotFound) {
			log.Println("Seeding database...")

			adminPassword, _ := authentication.HashPassword("fahmi")
			db.Create(&users.User{
				Email:     "fahminlb33@gmail.com",
				Password:  &adminPassword,
				FirstName: "Fahmi",
				LastName:  "Noor Fiqri",
				Role:      users.ADMIN,
			})

			log.Println("Seeding database... Done")
		}
	}

	return nil
}
