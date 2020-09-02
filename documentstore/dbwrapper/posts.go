package dbwrapper

import (
	"gitlab.com/shitposting/autoposting-bot/documentstore"
	"gitlab.com/shitposting/autoposting-bot/documentstore/entities"
	"time"
)

// AddPost adds a post to the database.
func AddPost(addedBy int32, media entities.Media, caption string) error {
	return documentstore.AddPost(addedBy, media, caption, documentstore.PostCollection)
}

// FindPostByFeatures finds a post by its features.
func FindPostByFeatures(histogram []float64, pHash string) (post entities.Post, err error) {
	return documentstore.FindPostByFeatures(histogram, pHash,
		documentstore.MediaApproximation, documentstore.SimilarityThreshold, documentstore.PostCollection)
}

// FindPostByUniqueID retrieves a post via its uniqueID.
func FindPostByUniqueID(uniqueID string) (post entities.Post, err error) {
	return documentstore.FindPostByUniqueID(uniqueID, documentstore.PostCollection)
}

// DeletePostByUniqueID deletes a post entity via its uniqueID.
func DeletePostByUniqueID(uniqueID string) error {
	return documentstore.DeletePostByUniqueID(uniqueID, documentstore.PostCollection)
}

// GetQueueLength returns the number of the enqueued posts.
func GetQueueLength() (length int64) {
	return documentstore.GetQueueLength(documentstore.PostCollection)
}

// GetNextPost retrieves the oldest post in the queue (not yet posted).
func GetNextPost() (entities.Post, error) {
	return documentstore.GetNextPost(documentstore.PostCollection)
}

// GetQueuePositionByAddTime returns the position of the first post added before the input time.
func GetQueuePositionByAddTime(addedAt time.Time) (position int) {
	return documentstore.GetQueuePositionByAddTime(addedAt, documentstore.PostCollection)
}

// MarkPostAsPosted marks a post as posted.
func MarkPostAsPosted(post *entities.Post, messageID int) error {
	return documentstore.MarkPostAsPosted(post, messageID, documentstore.PostCollection)
}

// MarkPostAsFailed marks a post as failed.
func MarkPostAsFailed(post *entities.Post) error {
	return documentstore.MarkPostAsFailed(post, documentstore.PostCollection)
}

// UpdatePostCaptionByUniqueID updates the caption of a post given its uniqueID.
func UpdatePostCaptionByUniqueID(uniqueID, caption string) error {
	return documentstore.UpdatePostCaptionByUniqueID(uniqueID, caption, documentstore.PostCollection)
}

// MarkPostAsDeletedByMessageID marks a post as deleted.
func MarkPostAsDeletedByMessageID(messageID int64) error {
	return documentstore.MarkPostAsDeletedByMessageID(messageID, documentstore.PostCollection)
}
