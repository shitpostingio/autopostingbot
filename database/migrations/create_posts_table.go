package migrations

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/shitposting/autoposting-bot/database/entities"
)

// CreatePosts will create table for model `Post`
// will append "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci" to the SQL statement when creating table `posts`
func CreatePosts(db *gorm.DB) {
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
	db.AutoMigrate(&entities.Post{})

	// Add Foreign key to reference the id on users table with cascade onupdate
	db.Model(&entities.Post{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "CASCADE")
}

// DropPosts will drop the Posts table
func DropPosts(db *gorm.DB) {
	db.DropTableIfExists(&entities.Post{})
}
