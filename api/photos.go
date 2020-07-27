package api

import (
	"github.com/zelenin/go-tdlib/client"
)

func SendPhoto(chatID, replyToMessageID int64, remoteFileID, caption string, entities []*client.TextEntity) (*client.Message, error) {

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

//TODO: POSSIBILI PROBLEMI CON SIZE = 0
func GetPhotoFileInfoFromMessage(message *client.Message) *client.File {

	//
	photo := message.Content.(*client.MessagePhoto).Photo

	//There's no reason why photo.Sizes should be 0
	//but that would definitely cause issues
	return photo.Sizes[len(photo.Sizes)-1].Photo

}
