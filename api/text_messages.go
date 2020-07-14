package api

import "github.com/zelenin/go-tdlib/client"

func SendPlainTextMessage(chatID int64, text string) (*client.Message, error) {

	request := client.SendMessageRequest{
		ChatId: chatID,
		InputMessageContent: &client.InputMessageText{
			Text: &client.FormattedText{
				Text: text,
			},
		},
	}

	return tdlibClient.SendMessage(&request)

}

func SendTextMessage(chatID int64, text string) (*client.Message, error) {

	request := client.SendMessageRequest{
		ChatId: chatID,
		InputMessageContent: &client.InputMessageText{
			Text: &client.FormattedText{
				Text: text,
				Entities: []*client.TextEntity{	//TODO: ASSOLUTAMENTE SISTEMARE
					{
						Offset: 0,
						Length: 0,
						Type:   nil,
					},
				},
			},
		},
	}

	return tdlibClient.SendMessage(&request)

}
