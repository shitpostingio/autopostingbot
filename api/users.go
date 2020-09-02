package api

import (
	"github.com/zelenin/go-tdlib/client"
)

// GetUserByID returns a client.User given a userID.
func GetUserByID(userID int32) (*client.User, error) {
	user, err := tdlibClient.GetUser(&client.GetUserRequest{UserId: userID})
	return user, err
}
