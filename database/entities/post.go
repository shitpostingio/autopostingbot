package entities

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Post Entity
type Post struct {
	ID         uint       `gorm:"primary_key"`
	UserID     uint       `gorm:"index"`             // Foreign key (belongs to), tag `index` will create index for this column
	Type       string     `gorm:"type:varchar(191)"` // Type from telegram API
	Media      string     `gorm:"type:longtext"`     // The URL of the media from telegram API
	Caption    string     `sql:"type:varchar(192) CHARACTER SET utf8 COLLATE utf8_unicode_ci"`
	Categories []Category `gorm:"many2many:category_posts;"` // Post has and belongs to many categories, use `category_posts` as join table
	CreatedAt  time.Time  // Timestamp of the creation inside the database
}

// IsVideo returns true if p is a video, false otherwise.
func (p Post) IsVideo(db *gorm.DB) bool {
	db.Preload("Categories").Where("media = ?", p.Media).First(&p)
	for _, cat := range p.Categories {
		if cat.Name == "video" {
			return true
		}
	}

	return false
}

// IsImage returns true if p is an image, false otherwise.
func (p Post) IsImage(db *gorm.DB) bool {
	db.Preload("Categories").Where("media = ?", p.Media).First(&p)
	for _, cat := range p.Categories {
		if cat.Name == "image" {
			return true
		}
	}

	return false
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
