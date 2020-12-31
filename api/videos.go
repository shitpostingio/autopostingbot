package api

import (
	"github.com/zelenin/go-tdlib/client"
)

// ShareVideo shares a video to a certain chat.
// If replyToMessageID is not 0, the video will be in reply to that message id.
// caption and entities can be used to attach a message with markdown.
func ShareVideo(chatID, replyToMessageID int64, remoteFileID, caption string, entities []*client.TextEntity) (*client.Message, error) {

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

// UploadVideo shares a video to a certain chat.
// If replyToMessageID is not 0, the video will be in reply to that message id.
// caption and entities can be used to attach a message with markdown.
func UploadVideo(chatID, replyToMessageID int64, localFilePath, caption string, entities []*client.TextEntity) (*client.Message, error) {

	request := client.SendMessageRequest{
		ChatId:           chatID,
		ReplyToMessageId: replyToMessageID,
		InputMessageContent: &client.InputMessageVideo{
			Video: &client.InputFileLocal{
				Path: localFilePath,
			},
			Caption: &client.FormattedText{
				Text:     caption,
				Entities: entities,
			},
		},
	}

	return tdlibClient.SendMessage(&request)

}

// GetVideoFileInfoFromMessage returns the Video structure
// of a given client.Message.
func GetVideoFileInfoFromMessage(message *client.Message) *client.File {
	return message.Content.(*client.MessageVideo).Video.Video
}
