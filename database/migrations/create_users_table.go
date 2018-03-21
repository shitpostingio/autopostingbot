package migrations

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/shitposting/autoposting-bot/database/entities" // We need SQLite to perform migrations
)

// CreateUsers will create table for model `User`
// will append "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci" to the SQL statement when creating table `users`
func CreateUsers(db *gorm.DB) {
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci").CreateTable(entities.User{})
}

// DropUsers will drop the Users table
func DropUsers(db *gorm.DB) {
	db.DropTableIfExists(entities.User{})
}
