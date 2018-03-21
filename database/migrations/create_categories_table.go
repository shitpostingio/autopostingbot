package migrations

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/shitposting/autoposting-bot/database/entities"
)

// CreateCategories will create table for model `Category`
// will append "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci" to the SQL statement when creating table `categories`
func CreateCategories(db *gorm.DB) {
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
	db.CreateTable(&entities.Category{})

	// Add unique index for name
	db.Model(entities.Category{}).AddUniqueIndex("categories_name_unique", "name")
}

// DropCategories will drop the Categories table
func DropCategories(db *gorm.DB) {
	db.DropTableIfExists(&entities.Category{})
}
