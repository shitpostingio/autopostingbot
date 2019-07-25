package autopostingbot

import (
	"time"
)

//Fingerprint Entity
//Each `Fingerprint` belongs to a `Post`
//AHash represents the AverageHash
//PHash represents the PerceptionHash
type Fingerprint struct {
	ID        uint `gorm:"AUTO_INCREMENT"`
	PostID    uint `gorm:"not null"`
	AHash     string
	PHash     string
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time
}
