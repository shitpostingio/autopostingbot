package documentstore

import (
	"context"
	"fmt"
	"gitlab.com/shitposting/autoposting-bot/documentstore/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func AddUser(userID int32, collection *mongo.Collection) error {

	//
	if userID <= 0 {
		return fmt.Errorf("AddUser: userID <= 0: %d", userID)
	}

	//
	user := entities.User{
		TelegramID: userID,
		CreatedAt:  time.Now(),
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("AddUser: %v", err)
	}

	return nil

}

func UserIsAuthorized(userID int32, collection *mongo.Collection) bool {

	//
	if userID <= 0 {
		return false
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.M{"telegramid": userID}
	return collection.FindOne(ctx, filter, options.FindOne()).Err() == nil

}

func GetUsers(collection *mongo.Collection) (*mongo.Cursor, error) {

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	return collection.Find(ctx, bson.M{}, options.Find())

}

func DeleteUser(userID int32, collection *mongo.Collection) error {

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.M{"telegramid": userID}

	//
	_, err := collection.DeleteOne(ctx, filter, options.Delete())
	return err

}
