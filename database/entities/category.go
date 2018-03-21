package database

import (
	"time"
)

// Category Entity
type Category struct {
	ID        uint      `gorm:"AUTO_INCREMENT"`
	Name      string    `gorm:"type:varchar(191);not null;unique"` // Unique name for the category
	Posts     []Post    `gorm:"many2many:category_posts;"`         // Category has and belongs to many posts, use `category_posts` as join table
	CreatedAt time.Time // Timestamp of the creation inside the database
	UpdatedAt time.Time // Everytime it's update orm will touch this column
}
