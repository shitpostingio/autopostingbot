package autopostingbot

import (
	"time"
)

//Post Entity
//Each `Post` belongs to a `User` (User, UserID)
//Each `Post` belongs to a `Type` (Type, TypeID)
//Each `Post` can have one `Fingerprint` (Fingerprint)
type Post struct {
	ID          uint   `gorm:"primary_key"`
	User        User   `gorm:"foreignkey:UserID"`
	UserID      uint   `gorm:"not null"`
	MessageID   int    `gorm:"not null"`
	Type        Type   `gorm:"foreignkey:TypeID"`
	TypeID      uint   `gorm:"not null"`
	FileID      string `gorm:"type:varchar(190);unique;not null"`
	Caption     string
	Fingerprint *Fingerprint
	PostedAt    *time.Time `gorm:"default:null"`
	HasError    bool       `gorm:"default:false"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   *time.Time `gorm:"default:null"`
	DeletedAt   *time.Time `gorm:"default:null"`
}
