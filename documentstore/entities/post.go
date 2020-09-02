package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Post represents a post in the document store.
type Post struct {

	// ID is MongoDB's object ID.
	ID        primitive.ObjectID `bson:"_id,omitempty"`

	// AddedBy is the Telegram user ID of the person that added the post.
	AddedBy   int32

	// Media contains information on the media added.
	Media     Media

	// Caption represents the media caption.
	// In case it has markup, it must be encoded in HTML to
	// be processed correctly.
	Caption   string

	// HasError becomes true if there was an error while trying to
	// post the media.
	HasError  bool `bson:",omitempty"`

	// AddedAt is the timestamp of the addition to the database.
	AddedAt   time.Time

	// UpdatedAT is the timestamp of the last modification.
	UpdatedAt *time.Time `bson:",omitempty"`

	// PostedAt is the timestamp of the post on the channel.
	PostedAt  *time.Time `bson:",omitempty"`

	// MessageID is Tdlib's message ID for the post.
	// It must be noted that they are different from normal messageIDs.
	MessageID int64

	// DeletedAt is the timestamp of the deletion from the channel.
	DeletedAt *time.Time `bson:",omitempty"`
}
