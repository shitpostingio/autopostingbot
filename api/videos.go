package api

import (
	"github.com/zelenin/go-tdlib/client"
)

func SendVideo(chatID, replyToMessageID int64, remoteFileID, caption string, entities []*client.TextEntity) (*client.Message, error) {

	request := client.SendMessageRequest{
		ChatId:           chatID,
		ReplyToMessageId: replyToMessageID,
		InputMessageContent: &client.InputMessageVideo{
			Video: &client.InputFileRemote{
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

func GetVideoFileInfoFromMessage(message *client.Message) *client.File {
	return message.Content.(*client.MessageVideo).Video.Video
}
