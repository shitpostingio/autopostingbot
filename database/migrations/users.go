package migrations

import (
	"github.com/jinzhu/gorm"
	entities "gitlab.com/shitposting/datalibrary/entities/autopostingbot" // We need SQLite to perform migrations
)

// CreateUsers creates the table for the entity `User`
func CreateUsers(db *gorm.DB) error {
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
	return db.CreateTable(&entities.User{}).Error
}

// MigrateUsers performs the migration for the `users` table
func MigrateUsers(db *gorm.DB) error {
	return db.AutoMigrate(&entities.User{}).Error
}

// DropUsers drops the `users` table
func DropUsers(db *gorm.DB) error {
	return db.DropTableIfExists(&entities.User{}).Error
}
