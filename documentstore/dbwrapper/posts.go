package dbwrapper

import (
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/documentstore"
	"gitlab.com/shitposting/autoposting-bot/documentstore/entities"
)

func AddPost(addedBy int32, media entities.Media, caption *client.FormattedText) error {
	return documentstore.AddPost(addedBy, media, caption, documentstore.PostCollection)
}

func FindPostByFeatures(histogram []float64, pHash string) (post entities.Post, err error) {
	return documentstore.FindPostByFeatures(histogram, pHash, mediaApproximation, documentstore.PostCollection)
}

// FindPostByFileID retrieves a post via its fileID
func FindPostByUniqueID(uniqueID string) (post entities.Post, err error) {
	return documentstore.FindPostByUniqueID(uniqueID, documentstore.PostCollection)
}
