package api

import (
	"github.com/zelenin/go-tdlib/client"
)

// SharePhoto shares a photo to a certain chat.
// If replyToMessageID is not 0, the photo will be in reply to that message id.
// caption and entities can be used to attach a message with markdown.
func SharePhoto(chatID, replyToMessageID int64, remoteFileID, caption string, entities []*client.TextEntity) (*client.Message, error) {

	request := client.SendMessageRequest{
		ChatId:           chatID,
		ReplyToMessageId: replyToMessageID,
		InputMessageContent: &client.InputMessagePhoto{
			Photo: &client.InputFileRemote{
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

// UploadPhoto shares a photo to a certain chat.
// If replyToMessageID is not 0, the photo will be in reply to that message id.
// caption and entities can be used to attach a message with markdown.
func UploadPhoto(chatID, replyToMessageID int64, localFilePath, caption string, entities []*client.TextEntity) (*client.Message, error) {

	request := client.SendMessageRequest{
		ChatId:           chatID,
		ReplyToMessageId: replyToMessageID,
		InputMessageContent: &client.InputMessagePhoto{
			Photo: &client.InputFileLocal{
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

// GetPhotoFileInfoFromMessage returns the Photo structure
// of a given client.Message.
func GetPhotoFileInfoFromMessage(message *client.Message) *client.File {

	//
	photo := message.Content.(*client.MessagePhoto).Photo

	//There's no reason why photo.Sizes should be 0
	//but that would definitely cause issues
	return photo.Sizes[len(photo.Sizes)-1].Photo

}
