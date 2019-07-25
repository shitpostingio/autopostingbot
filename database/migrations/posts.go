package migrations

import (
	"github.com/jinzhu/gorm"
	entities "gitlab.com/shitposting/datalibrary/entities/autopostingbot"
)

// CreatePosts creates the table for the entity `Post`
func CreatePosts(db *gorm.DB) error {

	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
	err := db.CreateTable(&entities.Post{}).Error
	if err != nil {
		return err
	}

	/* EACH POST BELONGS TO A USER */
	err = db.Model(&entities.Post{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "CASCADE").Error
	if err != nil {
		return err
	}

	/* EACH POST BELONGS TO A TYPE */
	return db.Model(&entities.Post{}).AddForeignKey("type_id", "types(id)", "RESTRICT", "CASCADE").Error
}

// MigratePosts performs the migration for the `posts` table
func MigratePosts(db *gorm.DB) error {
	return db.AutoMigrate(&entities.Post{}).Error
}

// DropPosts drops the `posts` table
func DropPosts(db *gorm.DB) error {
	return db.DropTableIfExists(&entities.Post{}).Error
}
