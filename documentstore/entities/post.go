package entities

import (
	"github.com/zelenin/go-tdlib/client"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	AddedBy   int32
	Media     Media
	Caption   *client.FormattedText
	HasError  bool `bson:",omitempty"`
	AddedAt   time.Time
	UpdatedAt time.Time `bson:",omitempty"`
	PostedAt  time.Time `bson:",omitempty"`
	MessageID int64 `bson:",omitempty"`
	DeletedAt time.Time `bson:",omitempty"`
}
