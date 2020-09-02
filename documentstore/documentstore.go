package documentstore

import (
	"context"
	"gitlab.com/shitposting/autoposting-bot/config/structs"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (

	/* TIMEOUT */
	opDeadline = 10 * time.Second

	/* COLLECTION NAMES */
	postCollectionName = "posts"
	userCollectionName = "users"
)

var (

	//
	dsCtx    context.Context
	database *mongo.Database

	// PostCollection is the MongoDB post collection.
	PostCollection *mongo.Collection

	// UserCollection is the MongoDB user collection.
	UserCollection *mongo.Collection

	// MediaApproximation is the approximation factor for similarity search in the database.
	MediaApproximation float64

	// SimilarityThreshold is the threshold for picture-to-picture similarity comparisons.
	SimilarityThreshold int
)

// Connect connects to the document store.
func Connect(cfg *structs.DocumentStoreConfiguration, mediaApproximation float64, similarityThreshold int) {

	//
	MediaApproximation = mediaApproximation
	SimilarityThreshold = similarityThreshold

	//
	client, err := mongo.Connect(context.Background(), cfg.MongoDBConnectionOptions())
	if err != nil {
		log.Fatal("Unable to connect to document store:", err)
	}

	pingCtx, cancelPingCtx := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelPingCtx()
	err = client.Ping(pingCtx, readpref.Primary())
	if err != nil {
		log.Fatal("Unable to ping document store:", err)
	}

	//
	dsCtx = context.TODO()

	//
	database = client.Database(cfg.DatabaseName)
	PostCollection = database.Collection(postCollectionName)
	UserCollection = database.Collection(userCollectionName)

}
