package entities

import (
	"time"
)

// User Entity
type User struct {
	ID         uint      `gorm:"AUTO_INCREMENT,primary_key"`
	TelegramID int       `gorm:"not null;unique"` // Telegram id from telegram API
	Posts      []Post    // One use can have many posts
	CreatedAt  time.Time // Timestamp of the creation inside the database
}
