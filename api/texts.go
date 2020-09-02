package api

import "github.com/zelenin/go-tdlib/client"

const (
	// NoReply can be used in the replyToMessageID field
	// to indicate that the message should not be in reply
	// to anything.
	NoReply = 0
)

// SendText sends a text message to a certain chat.
// If replyToMessageID is not 0, the text will be in reply to that message id.
// text and entities can be used to attach a message with markdown.
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

// SendPlainText simplifies the sending of a plain text message.
func SendPlainText(chatID int64, text string) (*client.Message, error) {
	return SendText(chatID, NoReply, text, nil)
}

// SendPlainReplyText simplifies the sending of a plain text message
// in reply to a certain replyToMessageID.
func SendPlainReplyText(chatID, replyToMessageID int64, text string) (*client.Message, error) {
	return SendText(chatID, replyToMessageID, text, nil)
}
