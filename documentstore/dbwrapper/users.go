package dbwrapper

import (
	"gitlab.com/shitposting/autoposting-bot/documentstore"
)

func UserIsAuthorized(userID int32) bool {
	return documentstore.UserIsAuthorized(userID, documentstore.UserCollection)
}
