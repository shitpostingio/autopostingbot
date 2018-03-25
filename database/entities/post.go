package entities

import (
	"time"
)

// Post Entity
type Post struct {
	ID         uint
	UserID     uint       `gorm:"index"`                     // Foreign key (belongs to), tag `index` will create index for this column
	Type       string     `gorm:"type:varchar(191)"`         // Type from telegram API
	Media      string     `gorm:"type:longtext"`             // The URL of the media from telegram API
	Caption    string     `gorm:"type:varchar(192)"`         // The caption of the current post, aka the test under the media
	Categories []Category `gorm:"many2many:category_posts;"` // Post has and belongs to many categories, use `category_posts` as join table
	CreatedAt  time.Time  // Timestamp of the creation inside the database
}

// Posts is a collection of Post, which implements the sort interface
type Posts []Post

func (p Posts) Len() int {
	return len(p)
}

func (p Posts) Less(i, j int) bool {
	return p[i].CreatedAt.Before(p[j].CreatedAt)
}

func (p Posts) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
