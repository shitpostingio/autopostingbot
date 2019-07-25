package migrations

import (
	"github.com/jinzhu/gorm"
	entities "gitlab.com/shitposting/datalibrary/entities/autopostingbot"
)

// CreateFingerprints creates the table for the entity `Fingerprint`
func CreateFingerprints(db *gorm.DB) error {
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
	err := db.CreateTable(&entities.Fingerprint{}).Error
	if err != nil {
		return err
	}

	/* EACH FINGERPRINT BELONGS TO A POST */
	return db.Model(&entities.Fingerprint{}).AddForeignKey("post_id", "posts(id)", "CASCADE", "CASCADE").Error
}

// MigrateFingerprints performs the migration for the `fingerprints` table
func MigrateFingerprints(db *gorm.DB) error {
	return db.AutoMigrate(&entities.Fingerprint{}).Error
}

// DropFingerprints drops the `fingerprints` table
func DropFingerprints(db *gorm.DB) error {
	return db.DropTableIfExists(&entities.Fingerprint{}).Error
}
