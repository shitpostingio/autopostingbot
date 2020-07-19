package api

import "github.com/zelenin/go-tdlib/client"

const (
	NoReply = 0
)

func SendText(chatID, replyToMessageID int64, text string, entities []*client.TextEntity) (*client.Message, error) {

	request := client.SendMessageRequest{
		ChatId:           chatID,
		ReplyToMessageId: replyToMessageID,
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
	return SendText(chatID, NoReply, text, nil)
}

func SendPlainReplyText(chatID, replyToMessageID int64, text string) (*client.Message, error) {
	return SendText(chatID, replyToMessageID, text, nil)
}
