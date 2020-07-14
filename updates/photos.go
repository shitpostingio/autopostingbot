package updates

import (
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"gitlab.com/shitposting/autoposting-bot/api"
)

func handlePhoto(message *client.Message) {

	//
	messageContent := message.Content.(*client.MessagePhoto)

	photoIndex := messageContent.Photo.Sizes[len(messageContent.Photo.Sizes) - 1]
	fileID := photoIndex.Photo.Remote.Id
	fileUniqueID := photoIndex.Photo.Remote.UniqueId

	log.Info("FileID: ", fileID, " UniqueID: ", fileUniqueID)
	_, _ = api.SendPlainPhoto(message.ChatId, fileID, "USH")


	//utf16Text := utf16.Encode([]rune(messageContent.Caption.Text))



}

