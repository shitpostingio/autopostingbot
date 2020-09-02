package dbwrapper

import (
	"gitlab.com/shitposting/autoposting-bot/documentstore"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddUser adds a user to the authorized admins.
func AddUser(userID int32) error {
	return documentstore.AddUser(userID, documentstore.UserCollection)
}

// UserIsAuthorized returns true if the user can interact with the bot.
func UserIsAuthorized(userID int32) bool {
	return documentstore.UserIsAuthorized(userID, documentstore.UserCollection)
}

// GetUsers retrieves all the authorized users from the database.
func GetUsers() (*mongo.Cursor, error) {
	return documentstore.GetUsers(documentstore.UserCollection)
}

// DeleteUser deletes a user from the database.
func DeleteUser(userID int32) error {
	return documentstore.DeleteUser(userID, documentstore.UserCollection)
}
