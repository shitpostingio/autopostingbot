package database

import (
	"github.com/jinzhu/gorm"
	entities "gitlab.com/shitposting/datalibrary/entities/autopostingbot"
)

// GetFingerprintsByAHash returns all the `Fingerprint` entities
// whose `AHash` is equal to the one passed
func GetFingerprintsByAHash(aHash string, db *gorm.DB) (fingerprints []entities.Fingerprint) {
	db.Where("a_hash = ?", aHash).Find(&fingerprints)
	return
}
