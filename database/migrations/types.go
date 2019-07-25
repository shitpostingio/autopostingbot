package migrations

import (
	"github.com/jinzhu/gorm"
	entities "gitlab.com/shitposting/datalibrary/entities/autopostingbot"
)

// CreateTypes will create the table for the entity `Type`
func CreateTypes(db *gorm.DB) error {
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
	return db.CreateTable(&entities.Type{}).Error
}

// MigrateTypes performs the migration for the `types` table
func MigrateTypes(db *gorm.DB) error {
	return db.AutoMigrate(&entities.Type{}).Error
}

// DropTypes drops the `types` table
func DropTypes(db *gorm.DB) error {
	return db.DropTableIfExists(&entities.Type{}).Error
}
