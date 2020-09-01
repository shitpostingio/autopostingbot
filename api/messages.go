package api

import "github.com/zelenin/go-tdlib/client"

func GetMessage(chatID, messageID int64) (*client.Message, error) {

	message, err := tdlibClient.GetMessage(&client.GetMessageRequest{
		ChatId:    chatID,
		MessageId: messageID,
	})

	return message, err

}

func GetMessageLink(chatID, messageID int64) (string, error) {

	link, err := tdlibClient.GetMessageLink(&client.GetMessageLinkRequest{
		ChatId:    chatID,
		MessageId: messageID,
	})

	if err != nil {
		return "", err
	}

	return link.Url, nil

}
