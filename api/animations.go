package api

import "github.com/zelenin/go-tdlib/client"

func SendAnimation(chatID int64, remoteFileID, caption string, entities []*client.TextEntity) (*client.Message, error) {

	request := client.SendMessageRequest{
		ChatId: chatID,
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
	return SendAnimation(chatID, remoteFileID, caption, nil)
}

func GetAnimationFileIDsFromMessage(message *client.Message) (fileID, fileUniqueID string) {
	animation := message.Content.(*client.MessageAnimation).Animation.Animation
	return animation.Remote.Id, animation.Remote.UniqueId
}
