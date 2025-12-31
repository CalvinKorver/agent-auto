package db

import (
	"fmt"
	"log"

	"carbuyer/internal/db/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(dbURL string) (*Database, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")

	return &Database{DB: db}, nil
}

func (d *Database) AutoMigrate() error {
	log.Println("Running database migrations...")

	// Handle migration from old schema (make/model as strings) to new schema (foreign keys)
	if d.DB.Migrator().HasTable(&models.UserPreferences{}) {
		// Check if old string columns exist
		hasOldMake := d.DB.Migrator().HasColumn(&models.UserPreferences{}, "make")
		hasOldModel := d.DB.Migrator().HasColumn(&models.UserPreferences{}, "model")

		if hasOldMake || hasOldModel {
			var count int64
			d.DB.Model(&models.UserPreferences{}).Count(&count)
			if count > 0 {
				log.Printf("Found %d existing user_preferences with old schema. Deleting old records...", count)
				log.Println("Note: Users will need to recreate their preferences through the API")
				// Delete old preferences since we can't migrate string values to foreign keys
				// without the makes/models tables being populated first
				if err := d.DB.Exec("DELETE FROM user_preferences").Error; err != nil {
					log.Printf("Warning: Could not delete old preferences: %v", err)
				}
			}

			// Drop old columns
			if hasOldMake {
				if err := d.DB.Migrator().DropColumn(&models.UserPreferences{}, "make"); err != nil {
					log.Printf("Warning: Could not drop old 'make' column: %v", err)
				}
			}
			if hasOldModel {
				if err := d.DB.Migrator().DropColumn(&models.UserPreferences{}, "model"); err != nil {
					log.Printf("Warning: Could not drop old 'model' column: %v", err)
				}
			}
		}
	}

	// Run AutoMigrate - it will create new columns and tables
	err := d.DB.AutoMigrate(
		&models.User{},
		&models.Make{},
		&models.Model{},
		&models.VehicleTrim{},
		&models.UserPreferences{},
		&models.Dealer{},
		&models.Thread{},
		&models.Message{},
		&models.TrackedOffer{},
		&models.GmailToken{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
