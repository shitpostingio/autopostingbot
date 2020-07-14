package documentstore

import (
	"context"
	"fmt"
	"gitlab.com/shitposting/autoposting-bot/documentstore/entities"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AddPost(addedBy int64,media entities.Media, caption string, collection *mongo.Collection) error {

	post := entities.Post{
		AddedBy:   addedBy,
		Media:     media,
		Caption:   caption,
		AddedAt:   time.Now(),
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	_, err := collection.InsertOne(ctx, post)
	if err != nil {
		err = fmt.Errorf("AddPost: %v", err)
	}

	return err

}

//// UpdatePostCaptionByFileID updates the caption of a post given its fileID
//func UpdatePostCaptionByFileID(fileID, caption string) bool {
//
//}
//
//// FindPostByFileID retrieves a post via its fileID
//func FindPostByFileID(fileID string) (post entities.Post) {
//
//}
//
//// FindPostByID retrieves a post entity via its database id
//func FindPostByID(id uint) (post entities.Post) {
//
//}
//
//// DeletePostByFileID deletes a post entity via its fileID
//func DeletePostByFileID(fileID string) error {
//
//}
//
//// GetNextPost retrieves the oldest media in the queue
//func GetNextPost() (entities.Post, error) {
//
//
//}
//
//// GetQueueLength returns the number of the enqueued posts
//func GetQueueLength() (length int) {
//
//}
//
//// GetQueuePositionByDatabaseID returns the position of the selected post in the queue
//func GetQueuePositionByDatabaseID(id uint) (position int) {
//
//}
//
//// MarkPostAsPosted marks a post as posted
//func MarkPostAsPosted(post entities.Post, messageID int) error {
//
//}
//
//// MarkPostAsFailed marks a post as failed
//func MarkPostAsFailed(post entities.Post) error {
//
//
//}
