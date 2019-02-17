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
	MessageID  int        `gorm:"not null"`                  // Message identification string, that Telegram uses to handle forwards
	CreatedAt  time.Time  // Timestamp of the creation inside the database
	PostedAt   time.Time  `gorm:"default:null"`  // Timestamp of the successful post on the channel
	HasError   bool       `gorm:"default:false"` // If true, this Post had some kind of posting error
	PHash      string     `gorm:"type:longtext"` // PerceptionHash for a photo
	AHash      string     `gorm:"type:longtext"` // AverageHash for a photo
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

// IsGIF returns true if p is a GIF, false otherwise.
func (p Post) IsGIF(db *gorm.DB) bool {
	db.Preload("Categories").Where("media = ?", p.Media).First(&p)
	for _, cat := range p.Categories {
		if cat.Name == "gif" {
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
