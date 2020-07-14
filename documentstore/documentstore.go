package documentstore

import (
	"context"
	configuration "gitlab.com/shitposting/autoposting-bot/config"
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

	PostCollection *mongo.Collection
	UserCollection *mongo.Collection

)

// Connect connects to the document store
func Connect(cfg *configuration.DocumentStoreConfiguration) {

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

	/* SAVE COLLECTIONS */
	database = client.Database(cfg.DatabaseName)
	PostCollection = database.Collection(postCollectionName)
	UserCollection = database.Collection(userCollectionName)

}
