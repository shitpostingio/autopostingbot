package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	AddedBy   int64
	Media     Media
	Caption   string
	HasError  bool `bson:",omitempty"`
	AddedAt   time.Time
	UpdatedAt time.Time `bson:",omitempty"`
	PostedAt  time.Time `bson:",omitempty"`
	DeletedAt time.Time `bson:",omitempty"`
}
