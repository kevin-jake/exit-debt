package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"pay-your-dues/internal/models"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(dsn string) (*Database, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(
		&models.User{},
		&models.UserSettings{},
		&models.Contact{},
		&models.UserContact{},
		&models.DebtList{},
		&models.DebtItem{},
		&models.Notification{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	// Add unique constraint on user_contacts (user_id, email)
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_user_contacts_user_email 
		ON user_contacts(user_id, email) 
		WHERE email IS NOT NULL
	`).Error; err != nil {
		return nil, fmt.Errorf("failed to create unique index on user_contacts: %v", err)
	}

	log.Println("Database connected and migrated successfully")

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
} 