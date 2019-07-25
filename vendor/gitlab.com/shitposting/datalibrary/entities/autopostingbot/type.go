package autopostingbot

import (
	"time"
)

//Type Entity
//Each `Type` has many `Post` (Posts)
type Type struct {
	ID        uint       `gorm:"AUTO_INCREMENT"`
	Name      string     `gorm:"type:varchar(191);not null;unique"`
	Posts     []Post     `gorm:"foreignkey:TypeID"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time `gorm:"default:null"`
	DeletedAt *time.Time `gorm:"default:null"`
}
