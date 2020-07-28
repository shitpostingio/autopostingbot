package dbwrapper

import (
	"gitlab.com/shitposting/autoposting-bot/documentstore"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserIsAuthorized(userID int32) bool {
	return documentstore.UserIsAuthorized(userID, documentstore.UserCollection)
}

func GetUsers() (*mongo.Cursor, error) {
	return documentstore.GetUsers(documentstore.UserCollection)
}
