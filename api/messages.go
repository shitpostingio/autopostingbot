package api

import "github.com/zelenin/go-tdlib/client"

func GetMessage(chatID, messageID int64) (*client.Message, error) {
	return tdlibClient.GetMessage(&client.GetMessageRequest{
		ChatId:    chatID,
		MessageId: messageID,
	})
}
