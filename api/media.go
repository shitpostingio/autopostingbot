package api

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
)

var (
	sendFunction = map[string]func(int64, int64, string, string, []*client.TextEntity) (*client.Message, error){
		client.TypeAnimation: SendAnimation,
		client.TypePhoto:     SendPhoto,
		client.TypeVideo:     SendVideo,
	}

	fileIDFunction = map[string]func(*client.Message) *client.File{
		client.TypeMessageAnimation: GetAnimationFileInfoFromMessage,
		client.TypeMessagePhoto:     GetPhotoFileInfoFromMessage,
		client.TypeMessageVideo:     GetVideoFileInfoFromMessage,
	}
)

func SendMedia(mediaType string, chatID, replyToMessageID int64, remoteFileID, caption string, entities []*client.TextEntity) (*client.Message, error) {

	send, found := sendFunction[mediaType]
	if !found {
		err := fmt.Errorf("send function not found for media type %s", mediaType)
		log.Error(err)
		return nil, err
	}

	return send(chatID, replyToMessageID, remoteFileID, caption, entities)

}

func GetMediaFileInfo(message *client.Message) (*client.File, error) {

	mediaType := message.Content.MessageContentType()
	getIDs, found := fileIDFunction[mediaType]
	if !found {
		err := fmt.Errorf("get file id function not found for media type %s", mediaType)
		log.Error(err)
		return nil, err
	}

	return getIDs(message), nil

}

func GetMediaFormattedText(message *client.Message) *client.FormattedText {

	switch message.Content.MessageContentType() {
	case client.TypeMessagePhoto:
		return message.Content.(*client.MessagePhoto).Caption
	case client.TypeMessageAnimation:
		return message.Content.(*client.MessageAnimation).Caption
	case client.TypeMessageVideo:
		return message.Content.(*client.MessageVideo).Caption
	default:
		return nil
	}

}
