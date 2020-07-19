package api

import "github.com/zelenin/go-tdlib/client"

func SendText(chatID int64, text string, entities []*client.TextEntity) (*client.Message, error) {

	request := client.SendMessageRequest{
		ChatId: chatID,
		InputMessageContent: &client.InputMessageText{
			Text: &client.FormattedText{
				Text:     text,
				Entities: entities,
			},
		},
	}

	return tdlibClient.SendMessage(&request)

}

func SendPlainText(chatID int64, text string) (*client.Message, error) {
	return SendText(chatID, text, nil)
}
