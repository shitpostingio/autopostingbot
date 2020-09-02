package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// User represents a user in the document store.
type User struct {

	// ID is MongoDB's object ID.
	ID         primitive.ObjectID `bson:"_id,omitempty"`

	// TelegramID is the user's telegram ID.
	TelegramID int32

	// CreatedAt is the timestamp of the addition of the user.
	CreatedAt  time.Time
}
