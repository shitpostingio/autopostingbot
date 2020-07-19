package api

import "github.com/zelenin/go-tdlib/client"

func SendAnimation(chatID, replyToMessageID int64, remoteFileID, caption string, entities []*client.TextEntity) (*client.Message, error) {

	request := client.SendMessageRequest{
		ChatId:           chatID,
		ReplyToMessageId: replyToMessageID,
		InputMessageContent: &client.InputMessageAnimation{
			Animation: &client.InputFileRemote{
				Id: remoteFileID,
			},
			Caption: &client.FormattedText{
				Text:     caption,
				Entities: entities,
			},
		},
	}

	return tdlibClient.SendMessage(&request)

}

func SendPlainAnimation(chatID int64, remoteFileID, caption string) (*client.Message, error) {
	return SendAnimation(chatID, NoReply, remoteFileID, caption, nil)
}

func SendPlainReplyAnimation(chatID, replyToMessageID int64, remoteFileID, caption string) (*client.Message, error) {
	return SendAnimation(chatID, replyToMessageID, remoteFileID, caption, nil)
}

func GetAnimationFileInfoFromMessage(message *client.Message) *client.File {
	return message.Content.(*client.MessageAnimation).Animation.Animation
}
