package api

import "github.com/zelenin/go-tdlib/client"

// GetMessage returns the client.Message with the input messageID
// in the input chatID.
func GetMessage(chatID, messageID int64) (*client.Message, error) {

	message, err := tdlibClient.GetMessage(&client.GetMessageRequest{
		ChatId:    chatID,
		MessageId: messageID,
	})

	return message, err

}
