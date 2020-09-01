package api

import "github.com/zelenin/go-tdlib/client"

func GetMessage(chatID, messageID int64) (*client.Message, error) {

	message, err := tdlibClient.GetMessage(&client.GetMessageRequest{
		ChatId:    chatID,
		MessageId: messageID,
	})

	return message, err

}

