package api

import "github.com/zelenin/go-tdlib/client"

func SendVideo(chatID int64, remoteFileID, caption string, entities []*client.TextEntity) (*client.Message, error) {

	request := client.SendMessageRequest{
		ChatId: chatID,
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

func SendPlainVideo(chatID int64, remoteFileID, caption string) (*client.Message, error) {
	return SendVideo(chatID, remoteFileID, caption, nil)
}

func GetVideoFileIDsFromMessage(message *client.Message) (fileID, fileUniqueID string) {
	video := message.Content.(*client.MessageVideo).Video.Video
	return video.Remote.Id, video.Remote.UniqueId
}
