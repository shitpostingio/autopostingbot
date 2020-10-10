package documentstore

import (
	"context"
	"errors"
	"fmt"
	"github.com/shitpostingio/autopostingbot/documentstore/entities"
	fpcompare "github.com/shitpostingio/image-fingerprinting/comparer"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/xerrors"
	"math"
	"time"
)

// AddPost adds a post to the database.
func AddPost(addedBy int32, media entities.Media, caption string, collection *mongo.Collection) error {

	//
	post := entities.Post{
		AddedBy: addedBy,
		Media:   media,
		Caption: caption,
		AddedAt: time.Now(),
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	_, err := collection.InsertOne(ctx, post)
	if err != nil {
		err = fmt.Errorf("AddPost: %v", err)
	}

	return err

}

// UpdatePostCaptionByUniqueID updates the caption of a post given its uniqueID.
func UpdatePostCaptionByUniqueID(uniqueID, caption string, collection *mongo.Collection) error {

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.M{"media.fileuniqueid": uniqueID}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "caption", Value: caption},
			},
		},
	}

	//
	_, err := collection.UpdateOne(ctx, filter, update, options.Update())
	return err

}

// FindPostByFeatures finds a post by its features.
func FindPostByFeatures(histogram []float64, pHash string, approximation float64, similarityThreshold int, collection *mongo.Collection) (post entities.Post, err error) {

	//
	if histogram == nil {
		err = xerrors.New("FindPostByFeatures: histogram was nil")
		return
	}

	//
	if pHash == "" {
		err = xerrors.New("FindPostByFeatures: pHash was empty")
		return
	}

	//
	average, sum := entities.GetHistogramAverageAndSum(histogram)
	minAvg := math.Trunc(average - 1)
	maxAvg := math.Ceil(average + 1)
	minSum := math.Trunc(sum - (sum * approximation))
	maxSum := math.Ceil(sum + (sum * approximation))

	//
	filter := bson.D{
		{
			Key: "media.histogramaverage",
			Value: bson.D{
				{Key: "$gte", Value: minAvg},
				{Key: "$lte", Value: maxAvg},
			},
		},
		{
			Key: "media.histogramsum",
			Value: bson.D{
				{Key: "$gte", Value: minSum},
				{Key: "$lte", Value: maxSum},
			},
		},
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		err = xerrors.Errorf("FindPostByFeatures: unable to retrieve post: %s", err)
		return
	}

	//
	post, err = findBestMatch(pHash, similarityThreshold, cursor)
	if err != nil {
		err = xerrors.Errorf("FindMediaByFeatures: %s", err)
		return
	}

	return

}

// FindPostByUniqueID retrieves a post via its uniqueID.
func FindPostByUniqueID(uniqueID string, collection *mongo.Collection) (post entities.Post, err error) {

	//
	if uniqueID == "" {
		return post, errors.New("uniqueID empty")
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.M{"media.fileuniqueid": uniqueID}

	//
	result := collection.FindOne(ctx, filter, options.FindOne())
	if result.Err() != nil {
		return post, result.Err()
	}

	//
	err = result.Decode(&post)
	return post, err

}

// DeletePostByUniqueID deletes a post entity via its uniqueID.
func DeletePostByUniqueID(uniqueID string, collection *mongo.Collection) error {

	//
	if uniqueID == "" {
		return errors.New("uniqueID empty")
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.M{"media.fileuniqueid": uniqueID}

	//
	_, err := collection.DeleteOne(ctx, filter, options.Delete())
	return err

}

// GetNextPost retrieves the oldest post in the queue (not yet posted).
func GetNextPost(collection *mongo.Collection) (post entities.Post, err error) {

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.D{
		{
			Key:   "messageid",
			Value: 0,
		},
		{
			Key:   "postedat",
			Value: nil,
		},
		{
			Key:   "haserror",
			Value: nil,
		},
	}

	//
	sortingOptions := options.FindOne().SetSort(bson.M{"addedat": 1})

	//
	err = collection.FindOne(ctx, filter, sortingOptions).Decode(&post)
	return

}

// GetQueueLength returns the number of the enqueued posts.
func GetQueueLength(collection *mongo.Collection) (length int64) {

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.D{
		{
			Key:   "postedat",
			Value: nil,
		},
		{
			Key:   "haserror",
			Value: nil,
		},
	}

	//
	res, err := collection.CountDocuments(ctx, filter, options.Count())
	if err != nil {
		return -1
	}

	return res

}

// GetQueuePositionByAddTime returns the position of the first post added before the input time.
func GetQueuePositionByAddTime(addedAt time.Time, collection *mongo.Collection) (position int) {

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.D{
		{
			Key:   "addedat",
			Value: bson.D{{Key: "$lte", Value: addedAt}},
		},
		{
			Key:   "postedat",
			Value: nil,
		},
		{
			Key:   "haserror",
			Value: nil,
		},
	}

	//
	res, err := collection.CountDocuments(ctx, filter, options.Count())
	if err != nil {
		return -1
	}

	return int(res)

}

// MarkPostAsPosted marks a post as posted.
func MarkPostAsPosted(post *entities.Post, messageID int, collection *mongo.Collection) error {

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.M{"_id": post.ID}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "messageid", Value: messageID},
				{Key: "postedat", Value: time.Now()},
			},
		},
	}

	//
	_, err := collection.UpdateOne(ctx, filter, update, options.Update())
	return err

}

// MarkPostAsFailed marks a post as failed.
func MarkPostAsFailed(post *entities.Post, collection *mongo.Collection) error {

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.M{"_id": post.ID}
	update := bson.D{
		{
			Key:   "$set",
			Value: bson.D{{Key: "haserror", Value: true}},
		}}

	//
	_, err := collection.UpdateOne(ctx, filter, update, options.Update())
	return err

}

// MarkPostAsDeletedByMessageID marks a post as deleted.
func MarkPostAsDeletedByMessageID(messageID int64, collection *mongo.Collection) error {

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.M{"messageid": messageID}
	update := bson.D{
		{
			Key:   "$set",
			Value: bson.D{{Key: "deletedat", Value: time.Now()}},
		}}

	//
	_, err := collection.UpdateOne(ctx, filter, update, options.Update())
	return err

}

// ============================================================================

// findBestMatch finds the best match given the input referencePHash.
func findBestMatch(referencePHash string, similarityThreshold int, cursor *mongo.Cursor) (post entities.Post, err error) {

	defer func() {
		_ = cursor.Close(dsCtx)
	}()

	i := 0
	for cursor.Next(context.TODO()) {

		i++
		// Support variable. If we deserialize directly in media,
		// since IsWhitelisted is an omitempty field, it won't be
		// deserialized in case of it being missing. This way, if
		// a document with it set to true has already been retrieved,
		// it will always keep being true.
		var res entities.Post
		err = cursor.Decode(&res)
		if err == nil && fpcompare.PhotoSimilarity(referencePHash, res.Media.PHash) < similarityThreshold {
			post = res
			log.Debugln("match in ", i, "iterations. FileID", post.Media.FileUniqueID)
			return
		}

	}

	err = xerrors.New("no match found")
	return

}
