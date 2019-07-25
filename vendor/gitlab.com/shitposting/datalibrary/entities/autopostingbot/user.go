package autopostingbot

import (
	"time"
)

//User Entity
//Each `User` has many `Posts`
type User struct {
	ID         uint   `gorm:"AUTO_INCREMENT,primary_key"`
	TelegramID int    `gorm:"not null;unique"`
	Handle     string `gorm:"type:varchar(32)"`
	Posts      []Post
	CreatedAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  *time.Time `gorm:"default:null"`
	DeletedAt  *time.Time `gorm:"default:null"`
}
