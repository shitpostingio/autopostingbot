package api

import (
	"fmt"
	"github.com/zelenin/go-tdlib/client"
)

// GetUserByID returns a client.User given a userID.
func GetUserByID(userID int32) (*client.User, error) {
	user, err := tdlibClient.GetUser(&client.GetUserRequest{UserId: userID})
	return user, err
}

// GetSenderUserID returns the sender user id, if the message
// was not sent on behalf of a chat
func GetSenderUserID(message *client.Message) (int32, error) {

	if message.Sender.MessageSenderType() != client.TypeMessageSenderUser{
		return 0, fmt.Errorf("the message was not sent by a user")
	}

	sender := message.Sender.(*client.MessageSenderUser)
	return sender.UserId, nil

}
