package api

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
)

var (

	sendFunction = map[string]func(int64, string, string, []*client.TextEntity) (*client.Message, error){
		client.TypeAnimation: SendAnimation,
		client.TypePhoto:     SendPhoto,
		client.TypeVideo:     SendVideo,
	}

	fileIDFunction = map[string]func(*client.Message) (string, string) {
		client.TypeAnimation: GetAnimationFileIDsFromMessage,
		client.TypePhoto:     GetPhotoFileIDsFromMessage,
		client.TypeVideo:     GetVideoFileIDsFromMessage,
	}

)

func SendMedia(mediaType string, chatID int64, remoteFileID, caption string, entities []*client.TextEntity) (*client.Message, error) {

	send, found := sendFunction[mediaType]
	if !found {
		err := fmt.Errorf("send function not found for media type %s", mediaType)
		log.Error(err)
		return nil, err
	}

	return send(chatID, remoteFileID, caption, entities)

}

func GetMediaFileIDs(message *client.Message) (fileID, fileUniqueID string) {

	mediaType := message.Content.MessageContentType()
	getIDs, found := fileIDFunction[mediaType]
	if !found {
		err := fmt.Errorf("send function not found for media type %s", mediaType)
		log.Error(err)
		return
	}

	return getIDs(message)

}
