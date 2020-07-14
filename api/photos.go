package api

import "github.com/zelenin/go-tdlib/client"

func SendPhoto(chatID int64, remoteFileID, caption string, entities []*client.TextEntity) (*client.Message, error) {

	request := client.SendMessageRequest{
		ChatId: chatID,
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

func SendPlainPhoto(chatID int64, remoteFileID, caption string) (*client.Message, error) {
	return SendPhoto(chatID, remoteFileID, caption, nil)
}

func GetPhotoFileIDsFromMessage(message *client.Message) (fileID, fileUniqueID string) {

	photo := message.Content.(*client.MessagePhoto).Photo
	if len(photo.Sizes) == 0 {
		return
	}

	return photo.Sizes[len(photo.Sizes) - 1].Photo.Remote.Id, photo.Sizes[len(photo.Sizes) - 1].Photo.Remote.UniqueId

}
