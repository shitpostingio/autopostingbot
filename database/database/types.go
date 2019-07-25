package database

import (
	"github.com/jinzhu/gorm"
	entities "gitlab.com/shitposting/datalibrary/entities/autopostingbot"
)

// GetTypeByName retrieves a type via its name
func GetTypeByName(typeName string, db *gorm.DB) (typeEntity entities.Type) {
	db.Where("name = ?", typeName).First(&typeEntity)
	return
}
