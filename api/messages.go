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

// GetMessageFormattedText returns the client.FormattedText structure for
// supported message types, nil otherwise.
func GetMessageFormattedText(mc client.MessageContent) *client.FormattedText {

	switch mc.MessageContentType() {
	case client.TypeMessageText:
		return mc.(*client.MessageText).Text
	case client.TypeMessagePhoto:
		return mc.(*client.MessagePhoto).Caption
	case client.TypeMessageAnimation:
		return mc.(*client.MessageAnimation).Caption
	case client.TypeMessageVideo:
		return mc.(*client.MessageVideo).Caption
	default:
		return nil
	}

}
